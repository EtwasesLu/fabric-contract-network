# Fabric 网络部署

## 环境
* Ubuntu 16.04.6-server
* Fabric 1.4.1

## 初始化Fabric 

### 将项目复制到Fabric的根目录下

### 设置工作路径
``` bash
export FABRIC_CFG_PATH=$GOPATH/src/github.com/hyperledger/fabric/portal-fabric-network/deploy
```
### 环境清理
``` bash
rm -fr config/*
rm -fr crypto-config/*
```

### 生成证书文件
``` bash
../../build/bin/cryptogen generate --config=./crypto-config.yaml
```

### 生成创世区块
``` bash
../../build/bin/configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./config/genesis.block
```

### 生成通道的创世交易
``` bash
../../build/bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./config/contractchannel.tx -channelID contractchannel
```

### 生成组织关于通道的锚节点（主节点）交易
``` bash
../../build/bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/Org1MSPanchors.tx -channelID contractchannel -asOrg Org1MSP
../../build/bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/Org2MSPanchors.tx -channelID contractchannel -asOrg Org2MSP
```

### 启动Fabric
``` bash
cd ~/go/src/github.com/hyperledger/fabric/portal-fabric-network/deploy/
./build.sh -- 启动
./stip.sh -- 停止
./teardown.sh -- 停止并清除docker镜像
``` 

### 进入cli
``` bash
docker exec -it cli bash
```

### 创建通道
``` bash
peer channel create -o orderer.liqi.com:7050 -c contractchannel -f /etc/hyperledger/configtx/contractchannel.tx
```

### 加入通道
``` bash
peer channel join -b contractchannel.block
```

### 设置主节点
``` bash
peer channel update -o orderer.liqi.com:7050 -c contractchannel -f /etc/hyperledger/configtx/Org1MSPanchors.tx
```

### 链码安装
``` bash
peer chaincode install -n contract -v 1.0.0 -l golang -p github.com/chaincode/contract
```

### 链码实例化
``` bash
peer chaincode instantiate -o orderer.liqi.com:7050 -C contractchannel -n contract -l golang -v 1.0.0 -c '{"Args":["init"]}'
```
### 链码交互
``` bash
peer chaincode invoke -C contractchannel -n contract -c '{"Args":["insertContract", "contract1", "goodName1", "goodCode1", "accountCode1", "100"]}'
peer chaincode invoke -C contractchannel -n contract -c '{"Args":["queryContract", "contract1"]}'
```
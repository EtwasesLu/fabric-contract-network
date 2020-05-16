/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// 合同
type Contract struct {
	Code string `json:"code"` // messagepack || protobuf
	GoodsName string `json:"goodsName"` //商品名称
	GoodsCode string `json:"goodsCode"` //商品编码
	AccountCode string `json:"accountCode"` //订单所属账户（买家）
	TotalPrice string `json:"totalPrice"` //总金额
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	println("chaincode contract init.")
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	switch function {
		case "insertContract":
			return s.insertContract(APIstub, args)
		case "queryContract":
			return s.queryContract(APIstub, args)
		default:
			return shim.Error(fmt.Sprintf("unsupported function: %s", function))
	}
}

// 新增合同
func (s *SmartContract) insertContract(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	// 套路1：检查参数的个数
	if len(args) != 5 {
		return shim.Error("not enough args")
	}

	// 套路2：验证参数的正确性
	code := args[0]
	goodsName := args[1]
	goodsCode := args[2]
	accountCode := args[3]
	totalPrice := args[4]
	if code == "" || goodsName == "" || goodsCode == "" || accountCode == "" || totalPrice == "" {
		return shim.Error("invalid args")
	}

	contractKey, err := stub.CreateCompositeKey("contract", []string{
		code,
	})
	if err != nil {
		return shim.Error(fmt.Sprintf("create key error: %s", err))
	}

	// 套路3：验证数据是否存在 应该存在 or 不应该存在
	if contractBytes, err := stub.GetState(contractKey); err == nil && len(contractBytes) != 0 {
		return shim.Error("contract already exist")
	}

	// 套路4：写入状态
	contract := &Contract{
		Code:   code,
		GoodsName:     goodsName,
		GoodsCode: goodsCode,
		AccountCode: accountCode,
		TotalPrice: totalPrice,
	}

	// 序列化对象
	contractBytes, err := json.Marshal(contract)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error %s", err))
	}

	if err := stub.PutState(contractKey, contractBytes); err != nil {
		return shim.Error(fmt.Sprintf("put contract error %s", err))
	}

	// 成功返回
	return shim.Success(nil)
}

// 合同查询
func (s *SmartContract) queryContract(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	// 套路1：检查参数的个数
	if len(args) != 1 {
		return shim.Error("not enough args")
	}

	// 套路2：验证参数的正确性
	contractCode := args[0]
	if contractCode == "" {
		return shim.Error("invalid args")
	}

	contractKey, err := stub.CreateCompositeKey("contract", []string{
		contractCode,
	})
	if err != nil {
		return shim.Error(fmt.Sprintf("create key error: %s", err))
	}

	// 套路3：验证数据是否存在 应该存在 or 不应该存在
	contractBytes, err := stub.GetState(contractKey)
	if err != nil || len(contractBytes) == 0 {
		return shim.Error("contract not found")
	}

	return shim.Success(contractBytes)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

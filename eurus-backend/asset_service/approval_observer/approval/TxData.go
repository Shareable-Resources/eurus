package approval

// func ExtractTransferState() func(string) int {
// 	// innerMap is captured in the closure returned below
// 	innerMap := map[string]int{
// 		"Failed": 0,
// 		"ABIError": 1,
// 		"Success": 2,
// 	}

// 	return func(key string) int {
// 		return innerMap[key]
// 	}
// }

// type TransferModel struct{
// 	Recipient *common.Address
// 	Amount *big.Int
// }

// func ExtractTransferData(txData []byte)(map[string]interface{},error, int){
// 	erc20ABI,err := abi.JSON(strings.NewReader(contract.EurusERC20ABI))
// 	// decode txInput method signature
// 	method, err := erc20ABI.MethodById(txData[:8])
// 	if err != nil {
// 		return nil,errors.New("Input Params Error"),ExtractTransferState()("Failed")
// 	}
// 	if(method.RawName!="transfer"){
// 		return nil,errors.New("Method not found"),ExtractTransferState()("Failed")
// 	}
// 	meth, ok := erc20ABI.Methods["transfer"]
// 	if !ok {
// 		return nil,errors.New("Method not found"),ExtractTransferState()("Failed")
// 	}
// 	dataMap :=make(map[string]interface{})
// 	err=meth.Inputs.UnpackIntoMap(dataMap,txData[4:])
// 	if(err!=nil){
// 		return nil, err, ExtractTransferState()("ABIError")
// 	}
// 	return dataMap, nil,ExtractTransferState()("Success")

// }

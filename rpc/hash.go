package rpc

import (
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/rpc/types"
)

func fmtExecuteCalldataStrings(nonce *big.Int, calls []types.FunctionCall) (calldataStrings []string) {
	callArray := fmtExecuteCalldata(nonce, calls)
	for _, data := range callArray {
		calldataStrings = append(calldataStrings, fmt.Sprintf("0x%s", data.Text(16)))
	}
	return calldataStrings
}

/*
Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func fmtExecuteCalldata(nonce *big.Int, calls []types.FunctionCall) (calldataArray []*big.Int) {
	callArray := []*big.Int{big.NewInt(int64(len(calls)))}

	for _, tx := range calls {
		address, _ := big.NewInt(0).SetString(tx.ContractAddress.Hex(), 0)
		callArray = append(callArray, address, caigo.GetSelectorFromName(tx.EntryPointSelector))

		if len(tx.CallData) == 0 {
			callArray = append(callArray, big.NewInt(0), big.NewInt(0))

			continue
		}

		callArray = append(callArray, big.NewInt(int64(len(calldataArray))), big.NewInt(int64(len(tx.CallData))))
		for _, cd := range tx.CallData {
			calldataArray = append(calldataArray, caigo.SNValToBN(cd))
		}
	}

	callArray = append(callArray, big.NewInt(int64(len(calldataArray))))
	callArray = append(callArray, calldataArray...)
	callArray = append(callArray, nonce)
	return callArray
}
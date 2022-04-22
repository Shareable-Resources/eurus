package sc_validator

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
)

func isSubSet(fullSet []common.Address, expectedList []string) (bool, []int) {
	var notFoundIndex []int

	for i, comparer := range expectedList {
		var isFound bool = false
		for _, addr := range fullSet {
			if bytes.Equal(common.HexToAddress(comparer).Bytes(), addr.Bytes()) {
				isFound = true
				break
			}
		}
		if !isFound {
			notFoundIndex = append(notFoundIndex, i)
		}
	}
	return true, notFoundIndex
}

func isSubSetAddressList(fullSet []common.Address, expectedList []common.Address) (bool, []int) {
	var notFoundIndex []int

	for i, comparer := range expectedList {
		var isFound bool = false
		for _, addr := range fullSet {
			if bytes.Equal(comparer.Bytes(), addr.Bytes()) {
				isFound = true
				break
			}
		}
		if !isFound {
			notFoundIndex = append(notFoundIndex, i)
		}
	}
	return true, notFoundIndex
}

func isAddressInStringList(fullSet []string, find common.Address) (bool, int) {
	for i, addrStr := range fullSet {
		if bytes.Equal(common.HexToAddress(addrStr).Bytes(), find.Bytes()) {
			return true, i
		}
	}
	return false, -1
}

func isAddressInList(fullSet []common.Address, find common.Address) (bool, int) {
	for i, addrStr := range fullSet {
		if bytes.Equal(addrStr.Bytes(), find.Bytes()) {
			return true, i
		}
	}
	return false, -1
}

func isEqualAddress(addr common.Address, addrStr string) bool {

	if bytes.Equal(common.HexToAddress(addrStr).Bytes(), addr.Bytes()) {
		return true
	}

	return false
}

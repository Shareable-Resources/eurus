package sc_validator

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func printValidationHeader(title string) {
	fmt.Println("========" + title + "========")
}

func printValidationResult(title string, isSuccess bool) {
	var caption string
	var color string

	if isSuccess {
		caption = "PASS"
		color = colorGreen
	} else {

		caption = "FAILED"
		color = colorRed
	}
	fmt.Printf("%s========%s %s========%s\r\n", color, title, caption, colorReset)
}

func printAddressList(title string, addrList []common.Address, config *ValidatorConfig) {
	fmt.Printf("%s. Count: %d\r\n", title, len(addrList))

	for _, addr := range addrList {
		agent := config.GetNameByAddress(addr)
		fmt.Printf("%s\t- %s\r\n", addr.Hex(), agent)
	}
	fmt.Println()
}

func printStringList(title string, addrList []string, config *ValidatorConfig) {
	fmt.Printf("%s. Count: %d\r\n", title, len(addrList))

	for _, addr := range addrList {
		agent := config.GetNameByAddress(common.HexToAddress(addr))
		fmt.Printf("%s\t- %s\r\n", addr, agent)
	}
	fmt.Println()
}

func itemToString(itemIndexList []int, strList []string) string {
	var message string = ""
	for _, i := range itemIndexList {
		message += fmt.Sprintf("Index:%d=>%s\r\n", i, strList[i])
	}
	return message
}

func itemAddressToString(itemIndexList []int, addrList []common.Address) string {
	var message string = ""
	for _, i := range itemIndexList {
		message += fmt.Sprintf("Index:%d=>%s\r\n", i, addrList[i].Hex())
	}
	return message
}

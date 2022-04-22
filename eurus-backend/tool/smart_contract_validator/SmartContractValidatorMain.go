package main

import (
	"eurus-backend/tool/smart_contract_validator/sc_validator"
	"fmt"
)

func stdinReader() int {
	var selectItem int
	fmt.Print("> ")
	fmt.Scanf("%d", &selectItem)
	return selectItem
}

func selectMenu(me *sc_validator.UserSelect) {
	fmt.Println("Please input the number of the network.\n 1. Mainnet\n 2. Side Chain\n 3. All Network\n 4. Exit")
	for {
		selectOption := stdinReader()
		if selectOption > int(sc_validator.Exit) || selectOption == int(sc_validator.StatusError) {
			fmt.Println("The number must be between 1 to ", sc_validator.Exit)
			continue
		}
		me.Network = sc_validator.NetworkSelect(selectOption)
		break
	}
}

func main() {
	var validator *sc_validator.SmartContractValidator = sc_validator.NewSmartContractValidator()
	err := validator.LoadConfig("ValidatorConfig.json", "SmartContractDeploy.json")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	var programStarted bool = false

	for {

		if programStarted == false {
			selectMenu(validator.UserSelect)
			if validator.UserSelect.Network == sc_validator.Exit {
				break
			}
			programStarted = true
		}

		err = validator.Connect()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		if validator.UserSelect.Network == sc_validator.SideChain || validator.UserSelect.Network == sc_validator.AllNetwork {
			err = validator.ValidatePlatformWallet()
			if err != nil {
				fmt.Println(err)
			}

			err = validator.ValidateWithdrawSmartContract()
			if err != nil {
				fmt.Println(err)
			}

			err = validator.ValidateApprovalWalletForUser()
			if err != nil {
				fmt.Println(err)
			}

			err = validator.ValidateERC20()
			if err != nil {
				fmt.Println(err)
			}

			err = validator.ValidateInternalSmartContractConfig()
			if err != nil {
				fmt.Println(err)
			}

			err = validator.ValidateExternalSmartContractConfig()
			if err != nil {
				fmt.Println(err)
			}
		}
		if validator.UserSelect.Network == sc_validator.MainNet || validator.UserSelect.Network == sc_validator.AllNetwork {
			err = validator.ValidateMainnetEurusInternalConfig()
			if err != nil {
				fmt.Println(err)
			}

			err = validator.ValidateMainnetEurusPlatformWalletConfig()
			if err != nil {
				fmt.Println(err)
			}
		}
		programStarted = false
	}
}

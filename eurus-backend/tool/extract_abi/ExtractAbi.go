package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	var filePath string
	flag.StringVar(&filePath, "file", "", "Input file path")
	flag.Parse()

	if filePath == "" {
		fmt.Println("No file input")
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	abiMap, ok := jsonMap["abi"]
	if !ok {
		fmt.Println("abi field not found")
		return
	}

	abiData, err := json.Marshal(abiMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(abiData))
}

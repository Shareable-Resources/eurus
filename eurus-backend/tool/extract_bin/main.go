package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Data struct {
	Abi              interface{} `json:"abi"`
	DeployedByteCode string      `json:"bytecode"`
}

func main() {

	inputFolder := flag.String("f", "input folder", "a string")
	abiFolder := flag.String("abi", "abi output folder", "a string")
	binFolder := flag.String("bin", "bin output folder", "a string")
	flag.Parse()

	fmt.Println(*inputFolder)
	files, err := ioutil.ReadDir(*inputFolder)
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(*abiFolder, 0777)
	os.MkdirAll(*binFolder, 0777)

	for _, f := range files {

		fileName := strings.Split(f.Name(), ".")[0]

		content, err := ioutil.ReadFile(*inputFolder + f.Name())
		if err != nil {
			log.Fatal(err)
		}

		data := new(Data)
		json.Unmarshal(content, data)
		jsonString, err := json.Marshal(data.Abi)

		if err != nil {
			log.Fatal(err)
		}
		f1, err := os.Create(*binFolder + fileName + ".bin")

		if err != nil {
			log.Fatal(err)
		}

		defer f1.Close()

		data.DeployedByteCode=strings.TrimPrefix(data.DeployedByteCode, "0x")
		_, err2 := f1.WriteString(data.DeployedByteCode)

		if err2 != nil {
			log.Fatal(err2)
		}

		f2, err := os.Create(*abiFolder + fileName + ".abi")

		if err != nil {
			log.Fatal(err)
		}

		defer f2.Close()

		_, err2 = f2.WriteString(string(jsonString))

		if err2 != nil {
			log.Fatal(err2)
		}

	}
}

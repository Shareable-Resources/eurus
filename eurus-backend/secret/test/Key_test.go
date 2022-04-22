package test

import (
	"eurus-backend/secret"
	"fmt"
	"strconv"
	"testing"
)

const privateKeyStr string = "rLoBAEHC8n5WjB+2FkQsqQh8jyHOpyZrg1HpMaeqYq8I/NvKswReJWutnBekVGQxM1qnmHgxY7lNCH4OSHMXow8ccpgtWTIpaZhbhRzsZxML1l0d0M5uTFwOozZAVjE6lHW5aglSNYsVuuNMcavDmo5ZChbztraJsby1nuJiqODcNRA2ybipFglbB6It8HMp13yhJZaerRMLptMUVsWxn4MN/XN1WINwXDmXTFIyHOPEl5dNQe28rhixWgIDHvUh6vRAUqPTKSP/S/UVzwwTwAxQW1zlixyKejmkGnEEDnWaxHXCDo5o22w6Pg2Cgg20TOFEbkZC6T6D1BnMixnw94ubV/EdQsslcK76r4S0Iiq6rOrwxfqVlH4YjHztG2A4pjYcHnmORL5zX9nyQdLg2Qtn0gI9GAqhghLGYMqkk54h7ftcm9WYvL3UrDmhiXmlSqWR2LA9PHSDoBmCQv/T6gvbVl370wBXDpEFlwo132AYLqPHK37f4/O0ohXB+y49QbquBiEBasz7W+jXkR5X6zROIOlEabYjA5K9DyeS5gP+YG0LydOEPdr2R7Igsgerz5wsRsKnnIWPghjEJmdV5TX9ioNLv01B/ilXVnFhJnV/YL8qCFXzTpZ8jC6dYzvq2Ee0FtxxMPxJ2cmlVMzr3G/mKhD02zFFSsOpvYA7bd4gjW7ZMW8mPCJxSWF4X6IO3HC05g+KVs7+6hyANXBNuHPgDwcYaiDvLe3uKxTM9Y/0P3js11vcthgQzxGgtkRb+wUdAvT7JcS46wYPFolO1CvNb5EwHqJqNUg/TzbvtMMTM2OTpnrM9TqIDe8DuraNjslKzjmiH70A/PppBWY+11lS2HxiBhA/takmf34ZVaKdULNLJ4fnA1XDiUaOW5++JDl7j4wTZLG9N2JHtvqGrkm+0YmpU5kgVlA3p1/+5z2Zy5sISypFBCL9zZtONR2rurwfn4kBqOVKGNv2MIUlAw+CAVLArimhM+eaJiNYlF/4qqhq4zP9zHf8R8HbYNtVEOtzj0By64a7OmKPsm9Kd7qoO5/zFzz1BvGxytt3bSdAZZ2No5uZutjb3/X4XXE/XNXELcx4Zpsjceg2a2E38sj3tLb24NtZIGZiEEUBKTYsW0oNB7zixyJ6mQhgBRCHin9pMLj31Iy4cY4iDL06NIKVOMgeg0ndQJUpV8k4dLfi6M5OSB+OHXLHbTdyrRS62l0M7DECoJ20hdBO07z3yW62Jq9prwKUB6OyRjRL9Fn+xHx4gXWzu7a8fn6N8tvK/s5C6g=="

func TestGenerateMintBurnKey(t *testing.T) {
	wallet, account, hdWalletAddr, _ := secret.GenerateMintBurnKey("please change it", privateKeyStr, strconv.FormatInt(55, 10))

	privateKey, err := wallet.PrivateKeyHex(*account)
	if err != nil {
		t.Errorf("private key error: %s", err.Error())
	}
	fmt.Println("Wallet address: ", hdWalletAddr)
	fmt.Println("Wallet Private Key: ", privateKey)
}

func TestGenerateCentralizeHdWallet(t *testing.T) {

	wallet, acc, err := secret.GenerateCentralizedUserMainnetPrivateKey("please change it", 301)
	if err != nil {
		t.Fatal(err)
	}
	privateKey, err := wallet.PrivateKeyHex(*acc)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Address: ", acc.Address.Hex())
	fmt.Println("Private key: ", privateKey)

}

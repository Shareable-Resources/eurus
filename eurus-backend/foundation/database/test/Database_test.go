package test

import (
	"eurus-backend/foundation/database"
	"fmt"
	"testing"
)

//var key string = "C639A572E14D5075C526FDDD43E4ECF6B095EA17783D32EF3D2710AF9F359DD4"
/*
func TestInitDatabase(t *testing.T) {
	keyBytes, err := hex.DecodeString(key)
	password := "admin123456"
	encryptedPassword, err := crypto.Encrypt([]byte(password), keyBytes)
	if err != nil {
		t.Errorf("%s", err)
	} else {
		t.Log(encryptedPassword)
	}
	database := database.Database{IP: "18.166.178.116", Port: 5432, UserName: "admin", Password: encryptedPassword, DBName: "postgres", SchemaName: "helloWorld", IdleConns: 10, MaxOpenConns: 100}
	err = database.SetDB(key)
	if err != nil {
		t.Errorf("%s", err)
	} else {
		t.Log("Database Connection is Success!")
	}

}
*/
func TestDatabasePW(t *testing.T){
	database := database.Database{IP: "18.166.178.116", Port: 9999, UserName: "admin", Password: "kHSUOn/LW3vdKHIYX0tH8ovMuEHMhpj+VcoDHjd6xBM=", DBName: "postgres", SchemaName: "helloWorld", IdleConns: 10, MaxOpenConns: 100}
	database.SetAESKey("XUFAKrxLKna5cZ2REBfFkii0btPBEehRApCbHPtQ6g8=")
	_,err:=database.GetConn()
	fmt.Println(err)
}
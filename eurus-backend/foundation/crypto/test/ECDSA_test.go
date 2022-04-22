package crypto

import (
	"crypto/ecdsa"
	"encoding/pem"
	"eurus-backend/foundation/crypto"
	crypto_ "github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestParseECDSA(t *testing.T){
	pemBytes,_:=pem.Decode([]byte("-----BEGIN PRIVATE KEY-----\nMIGEAgEAMBAGByqGSM49AgEGBSuBBAAKBG0wawIBAQQg3vVWGMC58PbJNfzy4XqEOFVkEq07KZBaza7jdorBkv6hRANCAARbu6c8hLdRkDcZg2Yts5sbBNT/9T4SANZSSQNfDRwDsD1ZQvnpYl1L3VgZG5j8tTT83JHD+scDoSO0mlyqW4PI\n-----END PRIVATE KEY-----"))
	key,err:=crypto.ParsePKCS8ECPrivateKey(pemBytes.Bytes)
	if(err!=nil){
		t.Errorf("%s", err)
		panic(err)
	}
	privateKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		t.Errorf("expected key to be of type *ecdsa.PrivateKey, but actual was %T", key)
	}
	addr:=crypto_.PubkeyToAddress(privateKey.PublicKey)
	t.Log(addr.Hex())
}
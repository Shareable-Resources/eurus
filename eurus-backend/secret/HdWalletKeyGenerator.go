package secret

import (
	"crypto/ecdsa"
	"eurus-backend/foundation/crypto"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

func GenerateServerHDWallet(mnemonicPhase string, password string, index string) (*hdwallet.Wallet, *accounts.Account, string, error) {
	return crypto.GenerateHdWalletKey(mnemonicPhase, mintBurnPassword+password, HdWalletServerIndex, index)
}

func GenerateMintBurnKey(mnemonicPhase string, password string, index string) (*hdwallet.Wallet, *accounts.Account, string, error) {
	return crypto.GenerateHdWalletKey(mnemonicPhase, mintBurnPassword+password, HdWalletServerIndex, index)
}

func GenerateCentralizedUserMainnetPrivateKey(userMainnetWalletMnemonic string, userID uint64) (*ecdsa.PrivateKey, *common.Address, error) {
	seed := bip39.NewSeed(userMainnetWalletMnemonic, userMainnetWalletPassPhrase)
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return nil, nil, err
	}

	path := hdwallet.MustParseDerivationPath(userMainnetWalletDerivedPath + strconv.FormatUint(userID, 10))

	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, nil, err
	}

	priKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, nil, err
	}

	return priKey, &account.Address, nil
}

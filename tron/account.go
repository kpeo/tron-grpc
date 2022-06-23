package tron

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	addr "github.com/fbsobreira/gotron-sdk/pkg/address"
)

// CreateAddress creates a new address and private key
func CreateAddress() (key string, address string) {
	priv, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return "", ""
	}
	for len(priv.D.Bytes()) != 32 {
		priv, err = btcec.NewPrivateKey(btcec.S256())
		if err != nil {
			continue
		}
	}
	key = hex.EncodeToString(priv.D.Bytes())
	address = addr.PubkeyToAddress(priv.ToECDSA().PublicKey).String()
	return
}

// CreateAddressBySeed creates a new address with seed bytes
func CreateAddressBySeed(seed []byte) (string, error) {
	if len(seed) != 32 {
		return "", fmt.Errorf("seed len=%d is not equal 32", len(seed))
	}
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), seed)
	if priv == nil {
		return "", errors.New("private key is nil")
	}
	a := addr.PubkeyToAddress(priv.ToECDSA().PublicKey)
	return a.String(), nil
}

// AddressB58ToHex converts base58 to hex string
func AddressB58ToHex(b58 string) (string, error) {
	a, err := addr.Base58ToAddress(b58)
	if err != nil {
		return "", err
	}
	return a.Hex(), nil
}

// AddressHexToB58 converts hex address to base58 string
func AddressHexToB58(hexAddress string) string {
	a := addr.HexToAddress(hexAddress)
	return a.String()
}

// clearKey zeroes a private key in memory
func ClearKey(k *ecdsa.PrivateKey) {
	b := k.D.Bits()
	for i := range b {
		b[i] = 0
	}
}
package tron

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
)

const test_Base58Address = "TNUC9Qb1rRpS5CbWLmNMxXBjyFoydXjWFR"
const test_HexAddress = "0x41891cdb91d149f23b1a45d9c5ca78a88d0cb44c18"

func Test_CreateAddress(t *testing.T) {
	key, addr := CreateAddress()
	if len(key) != 64 || len(addr) != 34 {
		t.Fatalf("Can't generate a new address and private key")
	}
}

func Test_CreateAddressBySeed(t *testing.T) {
	key, addr := CreateAddress()
	keyBytes, err := common.HexStringToBytes(key)
	if err != nil {
		t.Fatalf("HexStringToBytes error: %v", err)
	}

	addr2, err := CreateAddressBySeed(keyBytes)
	if err != nil {
		t.Fatalf("CreateAddressBySeed error: %v", err)
	}

	if addr2 != addr {
		t.Fatalf("Address created by seed (%s) is not match the required address (%s)", addr2, addr)
	}
}

func Test_AddressB58ToHex(t *testing.T) {
	hex, err := AddressB58ToHex(test_Base58Address)
	if err != nil {
		t.Fatalf("AddressB58ToHex error: %v", err)
	}

	if hex != test_HexAddress {
		t.Fatalf("hex addr=%s is not equal to base58 addr=%s", hex, test_HexAddress)
	}
}

func Test_AddressHexToB58(t *testing.T) {
	addr := AddressHexToB58(test_HexAddress)
	if addr != test_Base58Address {
		t.Fatalf("hex addr=%s is not equal to base58 addr=%s", test_HexAddress, test_Base58Address)
	}
}

func Test_ClearKey(t *testing.T) {
	key, _ := CreateAddress()
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		t.Fatalf("can't decode hex private key: %v", err)
	}

	priv := crypto.ToECDSAUnsafe(keyBytes)
	ClearKey(priv)
	for _, v := range priv.D.Bytes() {
		if v != 0 {
			t.Fatalf("Private key is not cleared")
			break
		}
	}
}
package tron

import (
	"testing"
)

func Test_CreateAddress(t *testing.T) {
	key, address := CreateAddress()
	if len(key) == 0 || len(address) == 0 {
		t.Error("Can't generate a new address and private key")
	}
}

func Test_CreateAddressBySeed(t *testing.T) {
	t.Error("Not implemented")
}

func Test_AddressB58ToHex(t *testing.T) {
	t.Error("Not implemented")
}

func Test_AddressHexToB58(t *testing.T) {
	t.Error("Not implemented")
}

func Test_ClearKey(t *testing.T) {
	t.Error("Not implemented")
}
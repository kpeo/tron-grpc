package tron

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/fatih/structs"
	addr "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

// Please fill the values for the tests
// (you can get them via https://nileex.io/join/getJoinPage)

const test_TxAddressFrom string = "TW8b5Z4mXbeRNQmAiSQiQrZSEZjgiQbhAW"
const test_TxAddressFromKey string = "4b9bfeb271fd55de4721d5f22f2b4068"
const test_TxAddressFromKey2 string = "1c6879a482e7250f01ecfa1a9ddd53d8"

const test_TxTransferAmount int64 = 100
const test_TxTransferFeeLimit int64 = 100000

const test_TxTrc10AssetId string = "1000016" // TRZ (TRONZ)
const test_TxTrc10Balance int64 = 240

const test_TxTrc20Contract string = "TXLAQ63Xg1NAzckPwKHvzw7CSEmLMEqcdj" //USDT
const test_TxTrc20Balance int64 = 11000

const test_TxTransactionId string = "ad8b0d5847d5c83a91a78d5ac288c5aba60fac4654a3c544be1493cf67b0e798"
const test_TxTransactionFee int64 = 100000

const test_TxAddressWithBalance string = "TXGE3trrKNejRrkF8BT7AQFd6ZAExP8qCd"
const test_TxAddressBalance int64 = 100

const test_TxTransferId string = "b30d2d93ade310b8d5545155683f7a62f6572f72a889ac7f86c2a10f63c57fe8"

const trc20TransferMethodSignature string = "a9059cbb"

func Test_TransferTrx(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNode, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("NewTronClient error: %v", err)
	}
	_, addr := CreateAddress()
	tx, err := c.Transfer(test_TxAddressFrom, addr, 100)
	if err != nil {
		t.Fatalf("Transfer error: %v", err)
	}
	signTx, err := SignTransaction(tx.Transaction, test_TxAddressFromKey+test_TxAddressFromKey2)
	if err != nil {
		t.Fatalf("SignTransaction error: %v", err)
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		t.Fatalf("BroadcastTransaction error: %v", err)
	}
	fmt.Printf("Transfer passed, txid=%s\n", common.BytesToHexString(tx.GetTxid()))
}

func Test_GetBalance(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNode, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("NewTronClient error: %v", err)
	}
	accb, err := c.GetTrxBalance(test_TxAddressWithBalance)
	if err != nil {
		t.Fatalf("GetTrxBalance error: %v", err)
	}
	if accb != test_TxAddressBalance {
		t.Fatalf("Current balance %d for address %s doesn't match the required value %d", accb, test_TxAddressWithBalance, test_TxAddressBalance)
	}
}

func Test_TransferTrc20(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNode, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("NewTronClient error: %v", err)
	}
	
	amount := big.NewInt(110)
	amount = amount.Mul(amount, big.NewInt(test_TxTransferAmount))

	_, addr := CreateAddress()

	tx, err := c.TransferTrc20(test_TxAddressFrom, addr,
			test_TxTrc20Contract, amount, test_TxTransferFeeLimit)
	if err != nil {
		t.Fatalf("TransferTrc20 error: %v", err)
	}
	signTx, err := SignTransaction(tx.Transaction, test_TxAddressFromKey+test_TxAddressFromKey2)
	if err != nil {
		t.Fatalf("SignTransaction error: %v", err)
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		t.Fatalf("BroadcastTransaction error: %v", err)
	}
	fmt.Printf("TransferTrc20 passed, txid=%s\n", common.BytesToHexString(tx.GetTxid()))
}

func Test_GetTrc20Balance(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNode, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("NewTronClient error: %v", err)
	}
	trc20, err := c.GetTrc20Balance(test_TxAddressWithBalance, test_TxTrc20Contract)
	if err != nil {
		t.Fatalf("GetTrc20Balance error: %v", err)
	}
	if trc20.Cmp(big.NewInt(test_TxTrc20Balance)) != 0 {
		t.Fatalf("Current TRC20 balance %d for address %s doesn't match the required value %d", trc20.Int64(), test_TxAddressWithBalance, test_TxTrc20Balance)
	} 
}

func Test_TransferTrc10(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNode, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("NewTronClient error: %v", err)
	}
	from, err := addr.Base58ToAddress(test_TxAddressFrom)
	if err != nil {
		t.Fatalf("Base58ToAddress(from) error: %v", err)
	}
	_, a := CreateAddress()
	to, err := addr.Base58ToAddress(a)
	if err != nil {
		t.Fatalf("Base58ToAddress(to) error: %v", err)
	}
	tx, err := c.grpc.TransferAsset(from.String(), to.String(), test_TxTrc10AssetId, 120)
	if err != nil {
		t.Fatalf("TransferAsset error: %v", err)
	}
	signTx, err := SignTransaction(tx.Transaction, test_TxAddressFromKey+test_TxAddressFromKey2)
	if err != nil {
		t.Fatalf("SignTransaction error: %v", err)
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		t.Fatalf("BroadcastTransaction error: %v", err)
	}
	fmt.Printf("TransferAsset passed, txid=%s\n", common.BytesToHexString(tx.GetTxid()))
}

func Test_GetTrc10Balance(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNode, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("NewTronClient error: %v", err)
	}
	amount, err := c.GetTrc10Balance(test_TxAddressWithBalance, test_TxTrc10AssetId)
	if err != nil {
		t.Fatalf("GetTrc10Balance error: %v", err)
	}
	if amount != test_TxTrc10Balance {
		t.Errorf("TRC10 balance (%d) doesn't match required value %d", amount, test_TxTrc10Balance)
	}
}

func Test_GetBlock(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNode, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("NewTronClient error: %v", err)
	}
	lb, err := c.grpc.GetNowBlock()
	if err != nil {
		t.Fatalf("GetNowBlock error: %v", err)
	}
	lbs := hex.EncodeToString(lb.Blockid)
	if len(lbs) != 64 {
		t.Errorf("Malformed Block ID: %s", lbs)

	}
	if lb.BlockHeader.RawData.Number < 1 {
		t.Errorf("Block number (%d) should be greater than 0", lb.BlockHeader.RawData.Number)
	}
}

func Test_GetTxInfoById(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNode, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("NewTronClient error: %v", err)
	}
	tx, err := c.grpc.GetTransactionInfoByID(test_TxTransactionId)
	if err != nil {
		t.Fatalf("GetTransactionInfoByID error: %v", err)
	}

	fee := tx.Receipt.GetEnergyFee() + tx.Receipt.GetNetFee()
	if fee != test_TxTransactionFee {
		t.Fatalf("Transaction fee (%d) doesn't match required value %d", fee, test_TxTransactionFee)
	}
}

func Test_GetTransaction(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNode, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("NewTronClient error: %v", err)
	}

	txById, err := c.grpc.GetTransactionByID(test_TxTransferId)
	if err != nil {
		t.Fatalf("GetTransactionByID error: %v", err)
	}
	var tc core.TriggerSmartContract
	if err = ptypes.UnmarshalAny(txById.GetRawData().GetContract()[0].GetParameter(), &tc); err != nil {
		t.Fatalf("UnmarshalAny error: %v", err)
	}
	tm := structs.Map(tc)
	tmData := tm["Data"]
	tmd := tmData.([]uint8)
	data := hex.EncodeToString(tmd)
	if !strings.HasPrefix(data, trc20TransferMethodSignature) {
		t.Errorf("TRC20 prefixed with incorrect signature")
	}
}

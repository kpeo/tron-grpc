package tron

import (
	"encoding/hex"
	"encoding/json"
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

const test_TronAddressFrom string = ""
const test_TronAddressFromPrivateKey string = ""
const test_TronAddressTo string = ""

const test_TronTransferAmount int64 = 100
const test_TronTransferFeeLimit int64 = 100000

const test_TronAssetId string = "1002000" // BTTOLD

const test_TronTransferContract string = ""
const test_TronTransactionId string = ""

const trc20TransferMethodSignature string = "a9059cbb"

func Test_TransferTrx(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNodeCorrect, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	tx, err := c.Transfer(test_TronAddressFrom, test_TronAddressTo, test_TronTransferAmount)
	if err != nil {
		fmt.Println(111)
		t.Fatal(err)
	}
	signTx, err := SignTransaction(tx.Transaction, test_TronAddressFromPrivateKey)
	if err != nil {
		fmt.Println(222)
		t.Fatal(err)
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		fmt.Println(333)
		t.Fatal(err)
	}
	fmt.Println(common.BytesToHexString(tx.GetTxid()))
}

func Test_GetBalance(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNodeCorrect, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	acc, err := c.GetTrxBalance(test_TronAddressFrom)
	if err != nil {
		t.Fatal(err)
	}
	d, _ := json.Marshal(acc)
	fmt.Println(string(d))
	fmt.Println(acc.GetBalance())

}

func Test_GetTrc20Balance(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNodeCorrect, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	amount, err := c.GetTrc20Balance(test_TronAddressFrom, test_TronAddressTo)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(amount.String())

}

func Test_TransferTrc20(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNodeCorrect, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	amount := big.NewInt(1)
	amount = amount.Mul(amount, big.NewInt(test_TronTransferAmount))
	tx, err := c.TransferTrc20(test_TronAddressFrom, test_TronAddressTo,
	test_TronTransferContract, amount, test_TronTransferFeeLimit)
	if err != nil {
		t.Fatal(err)
	}
	signTx, err := SignTransaction(tx.Transaction, test_TronAddressFromPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		t.Fatal(err)

	}
	fmt.Println(common.BytesToHexString(tx.GetTxid()))
}

func Test_TransferTrc10(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNodeCorrect, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	from, _ := addr.Base58ToAddress(test_TronAddressFrom)
	to, _ := addr.Base58ToAddress(test_TronAddressTo)
	tx, err := c.grpc.TransferAsset(from.String(), to.String(), test_TronAssetId, test_TronTransferAmount)
	if err != nil {
		t.Fatal(err)
	}
	signTx, err := SignTransaction(tx.Transaction, test_TronAddressFromPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	err = c.BroadcastTransaction(signTx)
	if err != nil {
		t.Fatal(err)

	}
	fmt.Println(common.BytesToHexString(tx.GetTxid()))
}

func Test_GetTrc10Balance(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNodeCorrect, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	amount, err := c.GetTrc10Balance(test_TronAddressFrom, test_TronAssetId)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(amount)
}

func Test_GetBlock(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNodeCorrect, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	lb, err := c.grpc.GetNowBlock()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(lb.BlockHeader.RawData.Number)
	fmt.Println(hex.EncodeToString(lb.Blockid))
}

func Test_GetTxByTxid(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNodeCorrect, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	ti, err := c.grpc.GetTransactionInfoByID(test_TronTransactionId)
	if err != nil {
		t.Fatal(err)
	}

	fee := ti.Receipt.GetEnergyFee() + ti.Receipt.GetNetFee()
	fmt.Println(fee)
}

func Test_GetTransaction(t *testing.T) {
	c, err := NewTronClient(test_TronGrpcNodeCorrect, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	txInfo, err := c.grpc.GetTransactionByID(test_TronTransactionId)
	if err != nil {
		t.Fatal(err)
	}
	d, _ := json.Marshal(txInfo)
	fmt.Println(string(d))
	r, err := c.grpc.GetTransactionInfoByID(test_TronTransactionId)
	if err != nil {
		t.Fatal(err)
	}
	dd, _ := json.Marshal(r)
	fmt.Println(string(dd))
	var cc core.TriggerSmartContract
	if err = ptypes.UnmarshalAny(txInfo.GetRawData().GetContract()[0].GetParameter(), &cc); err != nil {
		t.Fatal(err)
	}
	tv := structs.Map(cc)
	i := tv["Data"]
	da := i.([]uint8)
	data := hex.EncodeToString(da)
	if !strings.HasPrefix(data, trc20TransferMethodSignature) {
		t.Errorf("TRC20 prefixed with incorrect signature")
	}
}

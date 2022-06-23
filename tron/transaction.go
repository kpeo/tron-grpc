package tron

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"google.golang.org/protobuf/proto"
)

// GetTrc10Balance get TRC10 balance of assetId for address
func (c *TronClient) GetTrc10Balance(addr, assetId string) (int64, error) {
	err := c.keepConnect()
	if err != nil {
		return 0, err
	}
	acc, err := c.grpc.GetAccount(addr)
	if err != nil || acc == nil {
		return 0, fmt.Errorf("can't get account %s: %v", addr, err)
	}
	for key, value := range acc.AssetV2 {
		if key == assetId {
			return value, nil
		}
	}
	return 0, fmt.Errorf("can't find asset_id=%s for account %s", assetId, addr)
}

// GetTrxBalance get TRX balance for address
func (c *TronClient) GetTrxBalance(addr string) (*core.Account, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.grpc.GetAccount(addr)
}

// GetTrc20Balance get TRC20 balance of contract address for address 
func (c *TronClient) GetTrc20Balance(addr, contractAddress string) (*big.Int, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.grpc.TRC20ContractBalance(addr, contractAddress)
}

// Transfer transfer TRC10
func (c *TronClient) Transfer(from, to string, amount int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.grpc.Transfer(from, to, amount)
}

// TransferTrc10 transfer TRC10 by assetId
func (c *TronClient) TransferTrc10(from, to, assetId string, amount int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	fromAddr, err := address.Base58ToAddress(from)
	if err != nil {
		return nil, fmt.Errorf("incorrect from-address")
	}
	toAddr, err := address.Base58ToAddress(to)
	if err != nil {
		return nil, fmt.Errorf("incorrect to-address")
	}
	return c.grpc.TransferAsset(fromAddr.String(), toAddr.String(), assetId, amount)
}

// TransferTrc10 transfer TRC20 by contract address
func (c *TronClient) TransferTrc20(from, to, contract string, amount *big.Int, feeLimit int64) (*api.TransactionExtention, error) {
	err := c.keepConnect()
	if err != nil {
		return nil, err
	}
	return c.grpc.TRC20Send(from, to, contract, amount, feeLimit)
}

// BroadcastTransaction broadcast transaction to the chain
func (c *TronClient) BroadcastTransaction(transaction *core.Transaction) error {
	err := c.keepConnect()
	if err != nil {
		return err
	}
	result, err := c.grpc.Broadcast(transaction)
	if err != nil {
		return fmt.Errorf("broadcast transaction error: %v", err)
	}
	if result.Code != 0 {
		return fmt.Errorf("bad transaction: %v", string(result.GetMessage()))
	}
	if result.Result {
		return nil
	}
	d, _ := json.Marshal(result)
	return fmt.Errorf("tx send fail: %s", string(d))
}

func SignTransaction(transaction *core.Transaction, privateKey string) (*core.Transaction, error) {
	privateBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("can't decode hex private key: %v", err)
	}
	priv := crypto.ToECDSAUnsafe(privateBytes)
	defer ClearKey(priv)
	rawData, err := proto.Marshal(transaction.GetRawData())
	if err != nil {
		return nil, fmt.Errorf("can't unmarshall tx raw data: %v", err)
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	signature, err := crypto.Sign(hash, priv)
	if err != nil {
		return nil, fmt.Errorf("can't sign tx: %v", err)
	}
	transaction.Signature = append(transaction.Signature, signature)
	return transaction, nil
}

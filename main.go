package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"google.golang.org/grpc"

	"tron-grpc/tron"
)

// getEnv Load .env variables for map initialized with key's names
func getEnvMap(config map[string]string) error {
	// assume map is already initialized
	if len(config) == 0 {
		return errors.New("map is empty")
	}
	if err := godotenv.Load(); err != nil {
		return errors.New("error loading .env file")
	}
	for k := range config {
		if tmp, ok := os.LookupEnv(k); ok {
			config[k] = tmp
		} else {
			return errors.New("key is not present in .env file")
		}
	}
	return nil
}

// Temporary test code for interviewing
// Please fill the required variables in .env file to proceed 
func main() {
	config := map[string]string{
		"TRON_GRPC_URL": "",
		"TRON_FROM_ADDRESS": "",
		"TRON_FROM_KEY": "",
		"TRON_TRANSFER_AMOUNT": "",
	}

    if err := getEnvMap(config); err != nil {
		fmt.Printf("Can't load .env: %v\n", err)
	}

	tx_amount, err := strconv.ParseInt(config["TRON_TRANSFER_AMOUNT"], 10, 64)
	if err != nil {
		fmt.Printf("Can't convert %s: %v\n", config["TRON_TRANSFER_AMOUNT"], err)
	}
	key, to_address := tron.CreateAddress()
	fmt.Printf("New address: %s, Private key: %s\n", to_address, key)

	client, err := tron.NewTronClient(config["TRON_GRPC_URL"], grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Can't connect to node %s: %v\n", config["TRON_GRPC_URL"], err)
		return
	}

	tx, err := client.Transfer(config["TRON_FROM_ADDRESS"], to_address, tx_amount)
	if err != nil {
		fmt.Printf("Can't transfer %d TRX from %s to %s: %v\n", tx_amount, config["TRON_FROM_ADDRESS"], to_address, err)
		return
	}

	signTx, err := tron.SignTransaction(tx.Transaction, config["TRON_FROM_KEY"])
	if err != nil {
		fmt.Printf("Can't sign transaction: %v\n", err)
		return
	}
	
	err = client.BroadcastTransaction(signTx)
	if err != nil {
		fmt.Printf("Can't broadcast transaction: %v\n", err)
		return
	}

	fmt.Printf("Transaction hash: %s\n", common.BytesToHexString(tx.GetTxid()))
}

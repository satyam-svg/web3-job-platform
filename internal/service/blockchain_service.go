package service

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func VerifyPayment(txHash string) bool {
	client, err := ethclient.Dial("https://rpc-mumbai.maticvigil.com")
	if err != nil {
		log.Println("RPC connection error:", err)
		return false
	}
	defer client.Close()

	tx, isPending, err := client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil || isPending {
		log.Println("Transaction not found or still pending")
		return false
	}

	adminAddress := common.HexToAddress("5RAdGvEGs6SvNYif1yYqRSDUZhAbH6eMiwCzVhfmxYQ") // Replace this!
	if tx.To() != nil && *tx.To() == adminAddress {
		value := tx.Value()
		expected := big.NewInt(10000000000000000) // 0.01 MATIC in wei
		return value.Cmp(expected) >= 0
	}

	return false
}

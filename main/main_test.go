package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
)

// receipt.EffectiveGasPrice == tx.gasPrice
func TestIssue1(t *testing.T) {
	client, err := ethclient.Dial("https://bsc.nodereal.io")
	if err != nil {
		panic(err)
	}
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	txs := block.Transactions()
	if len(txs) != 0 {
		fmt.Println("Latest block:", block.Number())
		fmt.Printf("first tx hash:%s, gasPrice:%d\n", txs[0].Hash(), txs[0].GasPrice())
		receipt, err := client.TransactionReceipt(context.Background(), txs[0].Hash())
		if err != nil {
			panic(err)
		}
		fmt.Println("Latest transaction receipt:", receipt.EffectiveGasPrice)
	} else {
		fmt.Println("empty block")
	}
}

// sum(receipt.EffectiveGasPrice * receipt.gasUsed) == blockRewards
func TestIssue2(t *testing.T) {
	client, err := ethclient.Dial("https://bsc.nodereal.io")
	if err != nil {
		panic(err)
	}
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Latest block:", block.Number())
	txs := block.Transactions()
	totalRewards := big.NewInt(0)
	for _, tx := range txs {
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			panic(err)
		}
		val := receipt.EffectiveGasPrice.Mul(receipt.EffectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
		totalRewards.Add(totalRewards, val)
	}
	fmt.Println("total rewards:", totalRewards)
}

package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
	"sort"
	"sync"
)

const gwei = 1000000000

func main() {
	client, err := ethclient.Dial("https://bsc.nodereal.io")
	if err != nil {
		panic(err)
	}
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	blockNumber := block.Number()
	count := 50
	blockChan := make(chan *types.Block)
	wg := sync.WaitGroup{}
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			block, err := client.BlockByNumber(context.Background(), blockNumber.Sub(blockNumber, big.NewInt(1)))
			if err != nil {
				panic(err)
			}
			blockChan <- block
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(blockChan)
	}()
	//wg.Wait()
	//close(blockChan)
	file, err := os.OpenFile("count50blocks0925.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	nonZeroTxsCnt := big.NewInt(0)
	nonZeroTxsSum := big.NewInt(0)
	var allTxsGasPrice []*big.Int
	for block := range blockChan {
		fmt.Println("block number:", block.Number())
		text1 := "block number:" + block.Number().Text(10)
		_, err := fmt.Fprintln(file, text1)
		if err != nil {
			panic(err.Error())
		}
		txs := block.Transactions()
		for _, tx := range txs {
			if tx.GasPrice().Cmp(common.Big0) > 0 {
				nonZeroTxsCnt.Add(nonZeroTxsCnt, big.NewInt(1))
				nonZeroTxsSum.Add(nonZeroTxsSum, tx.GasPrice())
				allTxsGasPrice = append(allTxsGasPrice, tx.GasPrice())
			}
			fmt.Println("tx:", tx.Hash(), "gas price:", tx.GasPrice().Div(tx.GasPrice(), big.NewInt(gwei)))
			text2 := "tx: " + tx.Hash().String() + " gas price: " + tx.GasPrice().Div(tx.GasPrice(), big.NewInt(gwei)).Text(10)
			_, err := fmt.Fprintln(file, text2)
			if err != nil {
				panic(err.Error())
			}
		}
	}
	sort.Sort(bigIntArray(allTxsGasPrice))
	avgTxs := nonZeroTxsSum.Div(nonZeroTxsSum, nonZeroTxsCnt)
	//avgTxs.Div(avgTxs, big.NewInt(gwei))
	medianTxs := allTxsGasPrice[(len(allTxsGasPrice)-1)*50/100]
	//medianTxs.Div(medianTxs, big.NewInt(gwei))
	percentileTxs := allTxsGasPrice[(len(allTxsGasPrice)-1)*60/100]
	//percentileTxs.Div(percentileTxs, big.NewInt(gwei))
	fmt.Println("avgTxsGasPrice:", avgTxs.String(), "medianTxsGasPrice:", medianTxs.String(), "60%percentileTxsGasPrice:", percentileTxs.String())
	text := "avgTxsGasPrice: " + avgTxs.String() + " medianTxsGasPrice: " + medianTxs.String() + " 60%percentileTxsGasPrice: " + percentileTxs.String()
	_, err = fmt.Fprintln(file, text)
	if err != nil {
		panic(err.Error())
	}
}

type bigIntArray []*big.Int

func (s bigIntArray) Len() int           { return len(s) }
func (s bigIntArray) Less(i, j int) bool { return s[i].Cmp(s[j]) < 0 }
func (s bigIntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

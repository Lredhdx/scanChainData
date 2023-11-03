package main

import (
	"container/list"
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"sort"
)

const gwei = 1000000000

func main() {
	client, err := ethclient.Dial("https://bsc.nodereal.io")
	if err != nil {
		panic(err)
	}
	blockNumber := big.NewInt(33050608)
	avgQueue := list.New()
	medianQueue := list.New()

	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1", "blockNumber")
	f.SetCellValue("Sheet1", "B1", "avgGasPrice")
	f.SetCellValue("Sheet1", "C1", "medianGasPrice")
	f.SetCellValue("Sheet1", "D1", "latest20BlocksAvgGasPrice")
	f.SetCellValue("Sheet1", "E1", "latest20BlocksMedianGasPrice")
	count := 1
	for blockNumber.Cmp(big.NewInt(33079408)) < 0 {
		count++
		block, err := client.BlockByNumber(context.Background(), blockNumber.Add(blockNumber, big.NewInt(1)))
		if err != nil {
			block, _ = client.BlockByNumber(context.Background(), blockNumber)
		}
		nonZeroTxsCnt := big.NewInt(0)
		nonZeroTxsSum := big.NewInt(0)
		var allTxsGasPrice []*big.Int
		txs := block.Transactions()
		for _, tx := range txs {
			if tx.GasPrice().Cmp(common.Big0) > 0 {
				nonZeroTxsCnt.Add(nonZeroTxsCnt, big.NewInt(1))
				nonZeroTxsSum.Add(nonZeroTxsSum, tx.GasPrice())
				allTxsGasPrice = append(allTxsGasPrice, tx.GasPrice())
			}
		}
		avgTxs := big.NewInt(3)
		if nonZeroTxsCnt.Cmp(big.NewInt(0)) != 0 {
			avgTxs = nonZeroTxsSum.Div(nonZeroTxsSum, nonZeroTxsCnt)
		}
		sort.Sort(bigIntArray(allTxsGasPrice))
		medianGasPrice := big.NewInt(3)
		if len(allTxsGasPrice) != 0 {
			medianGasPrice = allTxsGasPrice[(len(allTxsGasPrice)-1)*50/100]
		}
		cell1 := fmt.Sprintf("A%d", count)
		f.SetCellValue("Sheet1", cell1, blockNumber)
		cell2 := fmt.Sprintf("B%d", count)
		f.SetCellValue("Sheet1", cell2, avgTxs.Div(avgTxs, big.NewInt(gwei)))
		cell3 := fmt.Sprintf("C%d", count)
		f.SetCellValue("Sheet1", cell3, medianGasPrice.Div(medianGasPrice, big.NewInt(gwei)))
		if avgQueue.Len() < 20 {
			avgQueue.PushBack(avgTxs)
		} else {
			sum := big.NewInt(0)
			for item := avgQueue.Front(); item != nil; item = item.Next() {
				sum = new(big.Int).Add(sum, item.Value.(*big.Int))
			}
			avg := sum.Div(sum, big.NewInt(20))
			cell := fmt.Sprintf("D%d", count)
			f.SetCellValue("Sheet1", cell, avg)
			avgQueue.Remove(avgQueue.Front())
			avgQueue.PushBack(avgTxs)
		}
		if medianQueue.Len() < 20 {
			medianQueue.PushBack(medianGasPrice)
		} else {
			var allMedianGasPrice []*big.Int
			for item := avgQueue.Front(); item != nil; item = item.Next() {
				allMedianGasPrice = append(allMedianGasPrice, item.Value.(*big.Int))
			}
			sort.Sort(bigIntArray(allMedianGasPrice))
			medianGasPriceAgain := allMedianGasPrice[(len(allMedianGasPrice)-1)*50/100]
			cell := fmt.Sprintf("E%d", count)
			f.SetCellValue("Sheet1", cell, medianGasPriceAgain)
			medianQueue.Remove(medianQueue.Front())
			medianQueue.PushBack(medianGasPrice)
		}
		fmt.Println("complete block number:", blockNumber)
		if count%1000 == 0 {
			err = f.SaveAs("20231030.xlsx")
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

type bigIntArray []*big.Int

func (s bigIntArray) Len() int           { return len(s) }
func (s bigIntArray) Less(i, j int) bool { return s[i].Cmp(s[j]) < 0 }
func (s bigIntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

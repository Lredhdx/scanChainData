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

	block, err := client.BlockByNumber(context.Background(), big.NewInt(400000000))
	if err != nil {

	}
	fmt.Println("blockNumber", block)

	//startBlockNumber1 := big.NewInt(7390000)  // 2021.5.14
	//startBlockNumber2 := big.NewInt(9570000)  // 2021.7.29
	//startBlockNumber3 := big.NewInt(12939600) // 2021.11.25
	//startBlockNumber4 := big.NewInt(17749200) // 2022.5.12
	//startBlockNumber5 := big.NewInt(33332309) // 2023.11.09
	//go func() {
	//	_, err := scanData(startBlockNumber1, "20210514.xlsx")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//go func() {
	//	_, err := scanData(startBlockNumber2, "20210729.xlsx")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//go func() {
	//	_, err := scanData(startBlockNumber3, "20211125.xlsx")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//go func() {
	//	_, err := scanData(startBlockNumber4, "20220512.xlsx")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//go func() {
	//	_, err := scanData(startBlockNumber5, "20231109.xlsx")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//time.Sleep(10 * time.Hour)
}

func scanData(blockNumber *big.Int, fileName string) (bool, error) {
	startBlock := big.NewInt(blockNumber.Int64())
	endBlock := blockNumber.Add(blockNumber, big.NewInt(86400))
	client, err := ethclient.Dial("https://bsc.nodereal.io")
	if err != nil {
		panic(err)
	}
	avgQueue := list.New()
	medianQueue := list.New()

	avgQueue40 := list.New()
	medianQueue40 := list.New()

	avgQueue60 := list.New()
	medianQueue60 := list.New()

	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1", "blockNumber")
	f.SetCellValue("Sheet1", "B1", "avgGasPrice")
	f.SetCellValue("Sheet1", "C1", "medianGasPrice")
	f.SetCellValue("Sheet1", "D1", "latest20BlocksAvgGasPrice")
	f.SetCellValue("Sheet1", "E1", "latest20BlocksMedianGasPrice")
	f.SetCellValue("Sheet1", "F1", "latest40BlocksAvgGasPrice")
	f.SetCellValue("Sheet1", "G1", "latest40BlocksMedianGasPrice")
	f.SetCellValue("Sheet1", "H1", "latest60BlocksAvgGasPrice")
	f.SetCellValue("Sheet1", "I1", "latest60BlocksMedianGasPrice")
	count := 1
	for startBlock.Cmp(endBlock) < 0 {
		count++
		block, err := client.BlockByNumber(context.Background(), startBlock.Add(startBlock, big.NewInt(1)))
		if err != nil {
			block, _ = client.BlockByNumber(context.Background(), startBlock)
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
		if avgQueue40.Len() < 40 {
			avgQueue40.PushBack(avgTxs)
		} else {
			sum := big.NewInt(0)
			for item := avgQueue40.Front(); item != nil; item = item.Next() {
				sum = new(big.Int).Add(sum, item.Value.(*big.Int))
			}
			avg := sum.Div(sum, big.NewInt(40))
			cell := fmt.Sprintf("F%d", count)
			f.SetCellValue("Sheet1", cell, avg)
			avgQueue40.Remove(avgQueue40.Front())
			avgQueue40.PushBack(avgTxs)
		}
		if medianQueue40.Len() < 40 {
			medianQueue40.PushBack(medianGasPrice)
		} else {
			var allMedianGasPrice []*big.Int
			for item := medianQueue40.Front(); item != nil; item = item.Next() {
				allMedianGasPrice = append(allMedianGasPrice, item.Value.(*big.Int))
			}
			sort.Sort(bigIntArray(allMedianGasPrice))
			medianGasPriceAgain := allMedianGasPrice[(len(allMedianGasPrice)-1)*50/100]
			cell := fmt.Sprintf("G%d", count)
			f.SetCellValue("Sheet1", cell, medianGasPriceAgain)
			medianQueue40.Remove(medianQueue40.Front())
			medianQueue40.PushBack(medianGasPrice)
		}
		if avgQueue60.Len() < 60 {
			avgQueue60.PushBack(avgTxs)
		} else {
			sum := big.NewInt(0)
			for item := avgQueue60.Front(); item != nil; item = item.Next() {
				sum = new(big.Int).Add(sum, item.Value.(*big.Int))
			}
			avg := sum.Div(sum, big.NewInt(60))
			cell := fmt.Sprintf("H%d", count)
			f.SetCellValue("Sheet1", cell, avg)
			avgQueue60.Remove(avgQueue60.Front())
			avgQueue60.PushBack(avgTxs)
		}
		if medianQueue60.Len() < 60 {
			medianQueue60.PushBack(medianGasPrice)
		} else {
			var allMedianGasPrice []*big.Int
			for item := medianQueue60.Front(); item != nil; item = item.Next() {
				allMedianGasPrice = append(allMedianGasPrice, item.Value.(*big.Int))
			}
			sort.Sort(bigIntArray(allMedianGasPrice))
			medianGasPriceAgain := allMedianGasPrice[(len(allMedianGasPrice)-1)*50/100]
			cell := fmt.Sprintf("I%d", count)
			f.SetCellValue("Sheet1", cell, medianGasPriceAgain)
			medianQueue60.Remove(medianQueue60.Front())
			medianQueue60.PushBack(medianGasPrice)
		}
		fmt.Println("complete block number:", startBlock)
		if count%1000 == 0 {
			err = f.SaveAs(fileName)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return true, nil
}

type bigIntArray []*big.Int

func (s bigIntArray) Len() int           { return len(s) }
func (s bigIntArray) Less(i, j int) bool { return s[i].Cmp(s[j]) < 0 }
func (s bigIntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

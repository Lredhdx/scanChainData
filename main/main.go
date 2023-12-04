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
	"time"
)

const gwei = 1000000000

func main() {
	startBlockNumber1 := big.NewInt(33562709) // 2023.11.17
	startBlockNumber2 := big.NewInt(33591509) // 2023.11.18
	startBlockNumber3 := big.NewInt(33620309) // 2023.11.19
	startBlockNumber4 := big.NewInt(33649109) // 2023.11.20
	startBlockNumber5 := big.NewInt(33677909) // 2023.11.21
	startBlockNumber6 := big.NewInt(33706709) // 2023.11.22
	startBlockNumber7 := big.NewInt(33735509) // 2023.11.23

	startBlockNumber8 := big.NewInt(33764309)  // 2023.11.24
	startBlockNumber9 := big.NewInt(33793109)  // 2023.11.25
	startBlockNumber10 := big.NewInt(33821909) // 2023.11.26
	startBlockNumber11 := big.NewInt(33850709) // 2023.11.27
	startBlockNumber12 := big.NewInt(33879509) // 2023.11.28
	startBlockNumber13 := big.NewInt(33908309) // 2023.11.29
	startBlockNumber14 := big.NewInt(33937109) // 2023.11.30

	go func() {
		_, err := scanData(startBlockNumber1, "20231117.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber2, "20231118.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber3, "20231119.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber4, "20231120.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber5, "20231121.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber6, "20231122.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber7, "20231123.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber8, "20231124.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber9, "20231125.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber10, "20231126.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber11, "20231127.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber12, "20231128.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber13, "20231129.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		_, err := scanData(startBlockNumber14, "20231130.xlsx")
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(10 * time.Hour)
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
			time.Sleep(3 * time.Minute)
			block, _ = client.BlockByNumber(context.Background(), startBlock)
		}
		nonZeroTxsCnt := big.NewInt(0)
		nonZeroTxsSum := big.NewInt(0)
		var allTxsGasPrice []*big.Int
		avgTxs := big.NewInt(3)
		if block != nil {
			txs := block.Transactions()
			for _, tx := range txs {
				if tx.GasPrice().Cmp(common.Big0) > 0 {
					nonZeroTxsCnt.Add(nonZeroTxsCnt, big.NewInt(1))
					nonZeroTxsSum.Add(nonZeroTxsSum, tx.GasPrice())
					allTxsGasPrice = append(allTxsGasPrice, tx.GasPrice())
				}
			}
			if nonZeroTxsCnt.Cmp(big.NewInt(0)) != 0 {
				avgTxs = nonZeroTxsSum.Div(nonZeroTxsSum, nonZeroTxsCnt)
			}
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

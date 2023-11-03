package main

import (
	"container/list"
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
	"sort"
	"sync"
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

// test deposit txs
func TestIssue3(t *testing.T) {
	var (
		DepositEventABI     = "TransactionDeposited(address,address,uint256,bytes)"
		DepositEventABIHash = crypto.Keccak256Hash([]byte(DepositEventABI))
	)
	//https://opbnb-mainnet.bnbchain.org
	// client, err := ethclient.Dial("https://bsc-mainnet.nodereal.io/v1/5bb2d1c141b648beb19ff604c6714c1d")
	client, err := ethclient.Dial("https://opbnb-mainnet.bnbchain.org")
	if err != nil {
		panic(err)
	}
	//_, err = client.BlockByNumber(context.Background(), new(big.Int).SetUint64(42238490))
	//if err != nil {
	//	fmt.Println("error")
	//	fmt.Println(err)
	//}
	// 32238558
	address := common.HexToAddress("0x1876EA7702C0ad0C6A2ae6036DE7733edfBca519")
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(42238490),
		ToBlock:   new(big.Int).SetUint64(42238590),
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{DepositEventABIHash}},
	}
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(logs))
}

// scan data
func TestIssue4(t *testing.T) {
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

func TestIssue5(t *testing.T) {
	client, err := ethclient.Dial("https://bsc.nodereal.io")
	if err != nil {
		panic(err)
	}
	blockNumber := big.NewInt(10000000)
	avgQueue := list.New()
	medianQueue := list.New()
	//avgGasPriceFile, err := os.OpenFile("latest20blocksAvgGasPrice.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	panic(err.Error())
	//}
	//defer avgGasPriceFile.Close()
	//
	//medianGasPriceFile, err := os.OpenFile("latest20blocksMedianGasPrice.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	panic(err.Error())
	//}
	//defer medianGasPriceFile.Close()
	//
	//singleGasPriceFile, err := os.OpenFile("singleBlockGasPriceFile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	panic(err.Error())
	//}
	//defer singleGasPriceFile.Close()

	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1", "blockNumber")
	f.SetCellValue("Sheet1", "B1", "avgGasPrice")
	f.SetCellValue("Sheet1", "C1", "medianGasPrice")
	f.SetCellValue("Sheet1", "D1", "latest20BlocksAvgGasPrice")
	f.SetCellValue("Sheet1", "E1", "latest20BlocksMedianGasPrice")
	count := 1
	for blockNumber.Cmp(big.NewInt(10000100)) < 0 {
		count++
		block, _ := client.BlockByNumber(context.Background(), blockNumber.Add(blockNumber, big.NewInt(1)))
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
		avgTxs := nonZeroTxsSum.Div(nonZeroTxsSum, nonZeroTxsCnt)
		sort.Sort(bigIntArray(allTxsGasPrice))
		medianGasPrice := allTxsGasPrice[(len(allTxsGasPrice)-1)*50/100]
		fmt.Println("complete block:", blockNumber)
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
			//fmt.Printf("block number: %d, latest 20 block avg gas price: %d\n", blockNumber, avg)
			//text := "block number:" + block.Number().Text(10) + "		latest 20 block avg gas price:" + avg.Text(10)
			//_, err := fmt.Fprintln(avgGasPriceFile, text)
			//if err != nil {
			//	panic(err.Error())
			//}
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
			medianGasPriceAgain := allTxsGasPrice[(len(allMedianGasPrice)-1)*50/100]
			//fmt.Printf("block number: %d, latest 20 block median gas price: %d\n",
			//	blockNumber, medianGasPriceAgain.Div(medianGasPriceAgain, big.NewInt(gwei)))
			//text := "block number:" + block.Number().Text(10) + "		latest 20 block median gas price:" +
			//	medianGasPriceAgain.Text(10)
			//_, err := fmt.Fprintln(medianGasPriceFile, text)
			//if err != nil {
			//	panic(err.Error())
			//}
			cell := fmt.Sprintf("E%d", count)
			f.SetCellValue("Sheet1", cell, medianGasPriceAgain.Div(medianGasPriceAgain, big.NewInt(gwei)))
			medianQueue.Remove(medianQueue.Front())
			medianQueue.PushBack(medianGasPrice)

		}
		//fmt.Println("block:", blockNumber, "avg gas price", avgTxs.Div(avgTxs, big.NewInt(gwei)),
		//	"median gas price", medianGasPrice.Div(medianGasPrice, big.NewInt(gwei)))
		//text := "block number:" + block.Number().Text(10) + "		avg gas price:" + avgTxs.Text(10) + "		median gas price:" + medianGasPrice.Text(10)
		//_, err := fmt.Fprintln(singleGasPriceFile, text)
		//if err != nil {
		//	panic(err.Error())
		//}
	}
	err = f.SaveAs("test.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

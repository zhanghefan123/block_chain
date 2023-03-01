// Package cmd /*
package cmd

import (
	BlockChain "block_chain/demo/block_chain"
	"block_chain/demo/tools"
	"block_chain/global"
	"bufio"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"os"
)

type TransactionInformation struct {
	src  string
	dest []string
	cost []float64
}

var TransactionsInformationMap = make(map[string]*TransactionInformation)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: `create a normal transaction to send money`,
	Long:  `create a normal transaction to send money`,
	Run: func(cmd *cobra.Command, args []string) {
		CreateBlockAndInsertTransactions()
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringVarP(&global.DestinationAddressOfMiner, global.DestinationAddressOfMinerFlag, "a", "", "address of the miner")
	sendCmd.Flags().StringVarP(&global.NameOfMiner, global.NameOfMinerFlag, "n", "", "name of the miner")
	sendCmd.Flags().StringVarP(&global.PathOfTransactionsFile, global.PathOfTransactionsFileFlag, "p", "", "path of transaction file")
	tools.MarkRequiredFlag(sendCmd, global.DestinationAddressOfMinerFlag)
	tools.MarkRequiredFlag(sendCmd, global.NameOfMinerFlag)
	tools.MarkRequiredFlag(sendCmd, global.PathOfTransactionsFileFlag)

}

// GetTransactionInformation 从文件中读取交易信息
func GetTransactionInformation() map[string]*TransactionInformation {
	path, err := os.Executable()
	if err != nil {
		return nil
	} else {
		fmt.Println(path)
	}
	file, err := os.Open(global.PathOfTransactionsFile)
	if err != nil {
		fmt.Printf("cannot open file: %v", err)
		return nil
	}
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if len(line) == 0 {
			// 在这里退出了程序，为什么ReadLine会出现错误呢？
			// 说明已经读取到了文件的末尾
			break
		} else if err != nil {
			fmt.Printf("read file failed, err: %v", err)
			return nil
		} else {
			// 说明还有数据
			// 读取数据
			var src, dest string
			var cost float64
			_, err := fmt.Sscanf(string(line), "%s %s %f", &src, &dest, &cost)
			if err != nil {
				fmt.Printf("read data failed, err: %v", err)
				return nil
			}
			// 这里需要判断是否有这个键
			if _, ok := TransactionsInformationMap[src]; ok {
				TransactionsInformationMap[src].dest = append(TransactionsInformationMap[src].dest, dest)
				TransactionsInformationMap[src].cost = append(TransactionsInformationMap[src].cost, cost)
			} else {
				TransactionsInformationMap[src] = &TransactionInformation{
					src:  src,
					dest: []string{dest},
					cost: []float64{cost},
				}
			}

		}
	}
	return TransactionsInformationMap
}

// CreateBlockAndInsertTransactions 创建区块并且进行交易的插入
func CreateBlockAndInsertTransactions() {
	// 从文件之中进行区块信息的获取
	transactionsInfo := GetTransactionInformation()
	// 创建一个区块链
	bc := BlockChain.CreateBlockChain()
	// 在添加区块的时候也必须填写一个挖矿人
	coinbase := BlockChain.CreateCoinBaseTransaction(global.DestinationAddressOfMiner, global.NameOfMiner)
	transactions := []*BlockChain.Transaction{coinbase}
	// 通过transactionsInfo进行交易的创建和插入
	// 这里创建了两个交易，两个交易都是普通的交易
	for _, txInfo := range transactionsInfo {
		// 按照信息进行普通的交易的创建,这里其实可以进行交易的合并
		// 可能一个源头存在多个不同的目的地址，这里就需要进行合并，
		// 否则如果一个两条交易会出现两个找零钱情况，这样就凭空生出来钱了。
		transaction := BlockChain.CreateNormalTransaction(txInfo.src, txInfo.dest, txInfo.cost, bc)
		transactions = append(transactions, transaction)
	}
	// 除了进行coinbase transaction的插入,我们还需要进行其他的交易的插入
	bc.CreateBlock(transactions, false)
	// 不能在Create完成之后就关闭了，以后还是需要进行使用的
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("close db failed, err: %v", err)
		}
	}(bc.Dao.DB)
}

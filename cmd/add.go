// Package cmd /*
package cmd

import (
	"block_chain/demo/block_chain"
	"block_chain/demo/tools"
	"block_chain/global"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: `Add block_chain into existed block_chain chain`,
	Long:  `Add block_chain into existed block_chain chain`,
	Run: func(cmd *cobra.Command, args []string) {
		AddBlock()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&global.DestinationAddressOfMiner, global.DestinationAddressOfMinerFlag, "a", "", "address of the miner")
	addCmd.Flags().StringVarP(&global.NameOfMiner, global.NameOfMinerFlag, "n", "", "name of the miner")
	tools.MarkRequiredFlag(addCmd, global.DestinationAddressOfMinerFlag)
	tools.MarkRequiredFlag(addCmd, global.NameOfMinerFlag)
}

func AddBlock() {
	// 创建一个区块链
	bc := block_chain.CreateBlockChain()
	// 在创建创世区块的时候就会进行tail的初始化
	// 在添加区块的时候也必须填写一个挖矿人
	coinbase := block_chain.CreateCoinBaseTransaction(global.DestinationAddressOfMiner, global.NameOfMiner)
	transactions := []*block_chain.Transaction{coinbase}
	bc.CreateBlock(transactions, false)
	// 不能在Create完成之后就关闭了，以后还是需要进行使用的
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("close db failed, err: %v", err)
		}
	}(bc.Dao.DB)
}

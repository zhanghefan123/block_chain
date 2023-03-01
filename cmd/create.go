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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: `Create Block Chain`,
	Long:  `Create Block Chain`,
	Run: func(cmd *cobra.Command, args []string) {
		CreateBlockChainWithAddress()
	},
}

func CreateBlockChainWithAddress() {
	// 创建一个区块链
	bc := block_chain.CreateBlockChain()
	// 在创建创世区块的时候就会进行tail的初始化
	coinbase := block_chain.CreateCoinBaseTransaction(global.DestinationAddressOfMiner, "Genesis Block")
	transactions := []*block_chain.Transaction{coinbase}
	bc.CreateBlock(transactions, true)
	// 不能在Create完成之后就关闭了，以后还是需要进行使用的
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("close db failed, err: %v", err)
		}
	}(bc.Dao.DB)
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&global.DestinationAddressOfMiner, global.DestinationAddressOfMinerFlag, "a", "", "address of the miner")
	tools.MarkRequiredFlag(createCmd, global.DestinationAddressOfMinerFlag)
}

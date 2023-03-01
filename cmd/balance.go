// Package cmd /*
package cmd

import (
	BlockChain "block_chain/demo/block_chain"
	"block_chain/demo/tools"
	"block_chain/global"
	"fmt"
	"github.com/spf13/cobra"
)

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: `get balance of a specified address in the block_chain chain`,
	Long:  `get balance of a specified address in the block_chain chain`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("用户 ", global.DestinationAddressOfMiner, " 的比特币的余额为:", calculateBalance(), "个")
	},
}

func calculateBalance() float64 {
	// 创建一个区块链
	bc := BlockChain.CreateBlockChain()
	UTXOs := bc.FindUTXOs(global.DestinationAddressOfMiner)
	var sum float64 = 0.0
	for _, utxo := range UTXOs {
		sum += utxo.Value
	}
	return sum
}

func init() {
	balanceCmd.Flags().StringVarP(&global.DestinationAddressOfMiner, global.DestinationAddressOfMinerFlag, "a", "", "address of the miner")
	tools.MarkRequiredFlag(balanceCmd, global.DestinationAddressOfMinerFlag)
}

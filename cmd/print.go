// Package cmd /*
package cmd

import (
	"block_chain/demo/block_chain"
	"fmt"
	"github.com/boltdb/bolt"

	"github.com/spf13/cobra"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "print the existed block_chain chain",
	Long:  `print the existed block_chain chain`,
	Run: func(cmd *cobra.Command, args []string) {
		// 拿到一个区块链
		bc := block_chain.GetBlockChain()
		if bc == nil {
			panic("block_chain chain has not been created")
		}
		// 不能在Create完成之后就关闭了，以后还是需要进行使用的
		defer func(db *bolt.DB) {
			err := db.Close()
			if err != nil {
				fmt.Printf("close db failed, err: %v", err)
			}
		}(bc.Dao.DB)
		// 打印所有的区块
		iterator := bc.CreateBlockChainIterator()
		for {
			block := iterator.Next()
			if block == nil {
				break
			} else {
				fmt.Println(block.Str())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}

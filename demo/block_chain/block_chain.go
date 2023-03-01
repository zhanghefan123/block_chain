package block_chain

import (
	"block_chain/demo/tools"
	"fmt"
)

// 进行常量的定义
// 数据库的名称
const nameOfDb = "test.db"

// 桶的名称
const nameOfBucket = "test_bucket"

// BlockChain 区块链结构
type BlockChain struct {
	Dao *BlockChainDao // 存储区块链的数据库
	// 存储的内容
	// 第一条数据: blockHash1 --> block1 最老的插入区块
	// 第二条数据: blockHash2 --> block2
	// 第三条数据: blockHash3 --> block3
	// ...
	// 最后一条数据: blockHashN --> blockN 最新的插入区块
	// 最后一条数据: l --> blockHashN
	Tail []byte // 最后一个区块的哈希
}

// CreateBlock 创建一个区块
func (bc *BlockChain) CreateBlock(txs []*Transaction, isGenesis bool) {
	if isGenesis {
		bc.Tail = []byte{}
	} else {
		bc.Tail = bc.Dao.GetLastHash()
	}
	tempBlock := &Block{
		Version:       0,
		PrevBlockHash: bc.Tail,
		TimeStamp:     tools.GetTimeStamp(),
		MerkleRoot:    []byte{},
		Difficulty:    0,
		Nonce:         0,
		Transactions:  txs,      // 这里为什么报错，
		Hash:          []byte{}, // 先填充空数据, 后面再计算
	}
	// 根据传入的交易进行梅克尔根的计算
	tempBlock.SetMerkleRoot()
	// ======================== 挖矿 ========================
	pow := CreateProofOfWork(tempBlock, 2)
	pow.Run()
	// ======================== 挖矿 ========================
	bc.Dao.WriteBlock(tempBlock.Hash, tempBlock)
	bc.Tail = tempBlock.Hash
}

func (bc *BlockChain) FindNeedUTXOs(srcAddr string, cost float64) (map[string][]int64, float64) {
	var needUtxos = make(map[string][]int64)
	var usedTxOutputs = make(map[string][]int64)
	var total float64
	// <========================================================>
	iter := bc.CreateBlockChainIterator()
	for {
		currentBlock := iter.Next()
		if currentBlock != nil {
			// 进行交易的遍历,找到自己的花费
			for _, tx := range currentBlock.Transactions {
				// 进行此交易之中的所有输出的遍历
				for index, output := range tx.TxOutputs {
					// 如果usedTxOutputs之中关于此条交易有记录
					if usedTxOutputs[string(tx.TXID)] != nil {
						for _, usedIndex := range usedTxOutputs[string(tx.TXID)] {
							if int64(index) == usedIndex {
								continue
							} else {
								if output.PubKeyHash == srcAddr {
									// 如果小于继续累加
									if total < cost {
										// 下面是未花费的UTXO
										// 找到自己需要的最少的UTXO
										total += output.Value
										needUtxos[string(tx.TXID)] = append(needUtxos[string(tx.TXID)], int64(index))
										// 加完之后满足条件
										if total >= cost {
											return needUtxos, total
										}
									}
								}
							}
						}
					} else { // 如果没有记录的话，说明此交易之中的所有的output都是未花费的
						if output.PubKeyHash == srcAddr {
							// 下面是未花费的UTXO
							// UTXOs = append(UTXOs, output)
							// 如果小于继续累加
							if total < cost {
								// 下面是未花费的UTXO
								// 找到自己需要的最少的UTXO
								total += output.Value
								needUtxos[string(tx.TXID)] = append(needUtxos[string(tx.TXID)], int64(index))
								// 加完之后满足条件
								if total >= cost {
									return needUtxos, total
								}
							}
						}
					}
				}
				// 如果当前交易是挖矿交易的时候，直接进行跳过即可
				if tx.IsCoinBaseTransaction() {
					continue
				}
				// 为什么要遍历input呢?
				// 我们的区块链是从后向前进行遍历的，所以我们首先看最后一层的输入，依赖于上一层的输出，
				// 上一层的输出一旦用过了，我们就不需要进行统计了,否则将出现双花的情况
				for index, input := range tx.TXInputs {
					// 如果是挖矿交易的话，那么就不需要进行统计了。因为这个交易不存在上一个交易的输出
					// 说明已经被使用过了
					if input.ScriptSig == srcAddr {
						// 这里的TxID是对应于这个input的上一个交易的输出的ID
						usedTxOutputs[string(input.TxID)] = append(usedTxOutputs[string(input.TxID)], int64(index))
					}
				}
			}
		} else {
			break
		}
	}

	// <========================================================>
	return needUtxos, total
}

// FindUTXOs 查找指定地址的所有的未花费的交易输出
func (bc *BlockChain) FindUTXOs(address string) []TxOutput {
	var UTXOs []TxOutput
	// 1. 遍历区块链
	// 2. 遍历区块中的交易
	usedTxOutputs := make(map[string][]int64)
	// 首先拿到iterator
	iter := bc.CreateBlockChainIterator()
	for {
		currentBlock := iter.Next()
		if currentBlock != nil {
			// 进行交易的遍历,找到自己的花费
			for _, tx := range currentBlock.Transactions {
				for index, output := range tx.TxOutputs {
					// 在这里做一个过滤,如果已经被使用过了,那么就不要再添加到集合中了
					if usedTxOutputs[string(tx.TXID)] != nil {
						for _, usedIndex := range usedTxOutputs[string(tx.TXID)] {
							if int64(index) == usedIndex {
								continue
							} else {
								if output.PubKeyHash == address {
									fmt.Printf("txid:%x,index:%d,value:%f\n", tx.TXID, index, output.Value)
									UTXOs = append(UTXOs, output)
								}
							}
						}
					} else {
						if output.PubKeyHash == address {
							fmt.Printf("txid:%x,index:%d,value:%f\n", tx.TXID, index, output.Value)
							UTXOs = append(UTXOs, output)
						}
					}
				}
				// 如果当前交易是挖矿交易的时候，直接进行跳过即可
				if tx.IsCoinBaseTransaction() {
					continue
				}
				// 为什么要遍历input呢?
				// 我们的区块链是从后向前进行遍历的，所以我们首先看最后一层的输入，依赖于上一层的输出，
				// 上一层的输出一旦用过了，我们就不需要进行统计了,否则将出现双花的情况
				for _, input := range tx.TXInputs {
					// 如果是挖矿交易的话，那么就不需要进行统计了。因为这个交易不存在上一个交易的输出
					// 说明已经被使用过了
					if input.ScriptSig == address {
						fmt.Printf("Used txid:%x,index:%d, ScriptSig: %s\n", input.TxID, input.Index, input.ScriptSig)
						// 这里的TxID是对应于这个input的上一个交易的输出的ID
						usedTxOutputs[string(input.TxID)] = append(usedTxOutputs[string(input.TxID)], input.Index)
					}
				}
			}
		} else {
			break
		}
	}

	return UTXOs
}

// CreateBlockChain 创建一个区块链
func CreateBlockChain() *BlockChain {
	// 首先创建一个
	bc := &BlockChain{
		Dao: CreateBlockChainDao(nameOfDb, nameOfBucket, false),
	}
	// 返回区块链对象
	return bc
}

func GetBlockChain() *BlockChain {
	return &BlockChain{
		Dao: CreateBlockChainDao(nameOfDb, nameOfBucket, true),
	}
}

// BlockChainIterator 定义区块链的迭代器方便进行遍历
type BlockChainIterator struct {
	Dao         *BlockChainDao
	CurrentHash []byte
}

// CreateBlockChainIterator 创建一个区块链的迭代器
func (bc *BlockChain) CreateBlockChainIterator() *BlockChainIterator {
	return &BlockChainIterator{
		Dao:         bc.Dao,
		CurrentHash: bc.Dao.GetLastHash(),
	}
}

// Next 区块链迭代器的方法，用于获取下一个区块
func (bci *BlockChainIterator) Next() *Block {
	if len(bci.CurrentHash) == 0 {
		return nil
	} else {
		currentBlock := bci.Dao.GetBlockByHash(bci.CurrentHash)
		bci.CurrentHash = currentBlock.PrevBlockHash
		return currentBlock
	}
}

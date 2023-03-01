package block_chain

import (
	"block_chain/demo/tools"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// ProofOfWork 定义一个工作量证明的结构
type ProofOfWork struct {
	// 1. 区块
	Block *Block
	// 2. 目标值,也叫做难度值,big.Int类似于java中的BigInteger
	Target *big.Int
}

// CreateProofOfWork 创建一个工作量证明对象
func CreateProofOfWork(block *Block, difficulty int) *ProofOfWork {
	// 创建一个指定长度的全部为0的字符串
	allZeroStr := tools.CreateFixedLengthBytes([]byte("0")[0], 64)
	tools.SetSpecifiedIndexInBytes(allZeroStr, difficulty, []byte("1")[0])
	// 将目标阈值字符串转换为字节数组
	tmpInt := big.Int{}
	// 将字符串转换为大整数，第二个参数为字符串的进制
	tmpInt.SetString(string(*allZeroStr), 16)
	return &ProofOfWork{
		Block:  block,
		Target: &tmpInt,
	}
}

// Run 开始挖矿,第一个返回值为挖矿成功后的区块哈希，第二个返回值为挖矿成功后的随机数
func (pow *ProofOfWork) Run() {
	for {
		// 1. 拼装数据 只需要对头进行哈希，因为体的变化由MerkleRoot来控制
		// 在将整体序列化到数据库的时候才需要对于整体进行哈希的操作
		blockInfo := tools.JoinBytes([][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.MerkleRoot, // 由transactions连接成一个大byte计算得到
			tools.ConvertUint64ToBytes(pow.Block.Version),
			tools.ConvertUint64ToBytes(pow.Block.TimeStamp),
			tools.ConvertUint64ToBytes(pow.Block.Difficulty),
			tools.ConvertUint64ToBytes(pow.Block.Nonce),
		}, []byte{})
		// 2. 生成哈希
		myHash := sha256.New()
		myHash.Write(blockInfo)
		result := myHash.Sum(nil)
		// 3. 将哈希值转换为大整数和目标值进行比较
		currentInt := big.Int{}
		currentInt.SetBytes(result)
		// 4. 如果当前的哈希值小于目标值，则表示挖矿成功
		if currentInt.Cmp(pow.Target) == -1 {
			fmt.Printf("挖矿成功,当前区块的哈希值为：%x 随机数为：%d\n", result, pow.Block.Nonce)
			pow.Block.Hash = result
			return
		} else {
			pow.Block.Nonce = pow.Block.Nonce + 1
		}
	}
}

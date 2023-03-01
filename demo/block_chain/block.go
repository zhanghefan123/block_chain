package block_chain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"strings"
)

// Block 区块结构
type Block struct {
	// 定义版本号
	Version uint64
	// 前一个区块的哈希
	PrevBlockHash []byte
	// 梅克尔根,根哈希值
	MerkleRoot []byte
	// 时间戳
	TimeStamp uint64
	// 难度值
	Difficulty uint64
	// 随机数,也叫做挖矿值
	Nonce uint64
	// 交易数据
	// Data []byte
	Transactions []*Transaction

	// 注意：当前区块的哈希,正常比特币区块之中没有当前区块的哈希,我们这里为了方便进行了简化
	Hash []byte
}

func (b *Block) Str() string {
	builder := strings.Builder{}
	sperator := fmt.Sprintf("<==============================================================================>\n")
	versionStr := fmt.Sprintf("Version: %d\n", b.Version)
	prevBlockHashStr := fmt.Sprintf("PrevBlockHash: %x\n", b.PrevBlockHash)
	merkleRootStr := fmt.Sprintf("MerkleRoot: %x\n", b.MerkleRoot)
	timeStampStr := fmt.Sprintf("TimeStamp: %d\n", b.TimeStamp)
	difficultyStr := fmt.Sprintf("Difficulty: %d\n", b.Difficulty)
	nonceStr := fmt.Sprintf("Nonce: %d\n", b.Nonce)

	hashStr := fmt.Sprintf("Hash: %x\n", b.Hash)
	builder.WriteString(sperator)
	builder.WriteString(versionStr)
	builder.WriteString(prevBlockHashStr)
	builder.WriteString(merkleRootStr)
	builder.WriteString(timeStampStr)
	builder.WriteString(difficultyStr)
	builder.WriteString(nonceStr)
	builder.WriteString(hashStr)
	// 将全部的交易的信息进行输出
	for index, tx := range b.Transactions {
		txStr := fmt.Sprintf("Transaction %d: %s\n", index, tx.TXInputs[0].ScriptSig)
		builder.WriteString(txStr)
	}
	builder.WriteString(sperator)
	return builder.String()
}

// Serialize 序列化
func (b *Block) Serialize() []byte {
	// 利用gob进行序列化
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	} else {
		return buffer.Bytes()
	}
}

// Deserialize 反序列化
func Deserialize(byteArray []byte) *Block {
	// 利用gob进行反序列化
	decoder := gob.NewDecoder(bytes.NewReader(byteArray))
	decodeBlock := Block{}
	err := decoder.Decode(&decodeBlock)
	if err != nil {
		panic(err)
	} else {
		return &decodeBlock
	}
}

func (b *Block) SetMerkleRoot() {
	var bigBytes []byte
	for _, tx := range b.Transactions {
		bigBytes = append(bigBytes, tx.TXID...)
	}
	// 将整体进行哈希运算
	myHash := sha256.New()
	myHash.Write(bigBytes)
	result := myHash.Sum(nil)
	b.MerkleRoot = result
}

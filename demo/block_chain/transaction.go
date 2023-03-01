package block_chain

import (
	"block_chain/global"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

// Transaction 交易结构的创建
type Transaction struct {
	TXID      []byte
	TXInputs  []TxInput
	TxOutputs []TxOutput
}

// TxInput 交易输入结构的定义
type TxInput struct {
	// 引用的交易的ID
	TxID []byte
	// 引用的某个交易之中输出的索引
	Index int64
	// 解锁脚本, 我们使用目的地址来进行模拟，实际上是需要公钥+签名
	ScriptSig string
}

// TxOutput 交易输出结构的定义
type TxOutput struct {
	// 转账金额
	Value float64
	// 锁定脚本, 我们使用目的地址来进行模拟
	PubKeyHash string
}

// SetTXID 使用交易的哈希值来设置交易的ID
func (tx *Transaction) SetTXID() {
	// 利用gob进行序列化
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	} else {
		// 将序列化之后的数据进行哈希计算
		myHash := sha256.New()
		myHash.Write(buffer.Bytes())
		result := myHash.Sum(nil)
		tx.TXID = result
	}
}

// CreateCoinBaseTransaction 创建挖矿产生的交易
// address 用来指定挖矿的人
func CreateCoinBaseTransaction(address string, data string) *Transaction {
	// 挖矿交易只有一个input
	// 无需引用交易ID
	// 无需引用index
	// 矿工由于挖矿的时候无需指定签名，所以这里的ScriptSig字段可以由矿工自由填写,一般是填写矿池的名字
	input := TxInput{[]byte{}, -1, data}
	// 挖矿交易只有一个output
	output := TxOutput{global.CurrentBlockReward, address}
	// 创建交易
	tx := Transaction{nil, []TxInput{input}, []TxOutput{output}}
	tx.SetTXID()
	// 返回交易
	return &tx
}

// CreateNormalTransaction 创建一笔普通的转账交易,这里先完成多个目的地址的转账
func CreateNormalTransaction(srcAddr string, destAddr []string, cost []float64, bc *BlockChain) *Transaction {
	fmt.Println("发起的源地址:", srcAddr)
	// 首先计算总的要求的输出量
	totalNeededOutput := 0.0
	for _, value := range cost {
		totalNeededOutput += value
	}
	// 0.创建 input 以及 output 数组
	var inputs []TxInput
	var outputs []TxOutput
	// 1.找到最合理的UTXOs集合
	bestUTXOs, resValue := bc.FindNeedUTXOs(srcAddr, totalNeededOutput)
	// 2.将这些UTXO逐一转换为input
	for id, array := range bestUTXOs {
		for _, value := range array {
			fmt.Printf("%x :%d %s\n", []byte(id), value, srcAddr)
			input := TxInput{[]byte(id), value, srcAddr}
			inputs = append(inputs, input)
		}
	}
	// 3.进行output的创建，这里可能能够创建多个output
	for i, addr := range destAddr {
		output := TxOutput{cost[i], addr}
		outputs = append(outputs, output)
	}
	// 4.如果有零钱，需要找零
	if resValue > totalNeededOutput {
		// 找0
		output := TxOutput{resValue - totalNeededOutput, srcAddr}
		outputs = append(outputs, output)
		fmt.Printf("找零：%f 个比特币\n", resValue-totalNeededOutput)
	} else if resValue < totalNeededOutput {
		fmt.Printf("余额不足,您的余额还剩 %f 个比特币, 您的目标金额 %f", resValue, cost)
	}
	// 创建的交易之中的inputs是转换而来的
	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetTXID()
	return &tx
}

// IsCoinBaseTransaction 判断当前的交易是否是挖矿交易
func (tx *Transaction) IsCoinBaseTransaction() bool {
	// 挖矿交易的特点
	// 1. 只有一个input
	// 2. input中的TxID为nil
	// 3. input中的index为-1
	if len(tx.TXInputs) == 1 && tx.TXInputs[0].TxID == nil && tx.TXInputs[0].Index == -1 {
		return true
	}
	return false
}

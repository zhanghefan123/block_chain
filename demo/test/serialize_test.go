package test

import (
	"block_chain/demo/block_chain"
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
)

func TestSerializeAndDeserializeWithGob(t *testing.T) {
	// 创建一个区块
	createdBlock := struct {
		name string
		age  int
	}{
		"zhf",
		20,
	}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(createdBlock)
	if err != nil {
		t.Errorf("序列化失败，错误信息为：%v\n", err)
	} else {
		fmt.Printf("序列化成功，序列化后的数据为：%x\n", buffer.Bytes())
	}
	// 反序列化
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	decodeBlock := block_chain.Block{}
	err = decoder.Decode(&decodeBlock)
	if err != nil {
		t.Errorf("反序列化失败，错误信息为：%v\n", err)
	}
	// 打印反序列化的结果
	fmt.Println(decodeBlock.Str())

}

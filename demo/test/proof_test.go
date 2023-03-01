package test

import (
	"block_chain/demo/tools"
	"testing"
)

func TestProof(t *testing.T) {
	// 选择一个字符，将字符转换为byte
	str := "0"
	character := []byte(str)[0]
	t.Log(str, "转换为byte后的值为：", character)
	// 创建一个64位的二进制字符串
	result := tools.CreateFixedLengthBytes(character, 64)
	tools.SetSpecifiedIndexInBytes(result, 0, 65)
	t.Log("创建的二进制字符串为：", string(*result))
}

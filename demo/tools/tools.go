package tools

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

func ConvertUint64ToBytes(number uint64) []byte {
	// 将uint64转换为[]byte
	var buffer bytes.Buffer
	// 写入buffer，为什么使用大端序？
	// 网络传输过程采用的是大端序，所以我们这里也采用大端序
	err := binary.Write(&buffer, binary.BigEndian, number)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buffer.Bytes()
}

// GetTimeStamp 获取当前时间戳
func GetTimeStamp() uint64 {
	// 获取当前时间戳
	return uint64(time.Now().Unix())
}

// CreateFixedLengthBytes 创建一个固定长度的字符串
func CreateFixedLengthBytes(c byte, length int) *[]byte {
	var byteArray []byte
	for i := 0; i < length; i++ {
		byteArray = append(byteArray, c)
	}
	return &byteArray
}

func SetSpecifiedIndexInBytes(target *[]byte, index int, value byte) {
	(*target)[index] = value
}

// JoinBytes 将[][]byte之中的每一个byte数组，以sep为分隔，拼接为[]byte
func JoinBytes(source [][]byte, sep []byte) []byte {
	return bytes.Join(source, sep)
}

// MarkRequiredFlag 将一个选项标识为必须的
func MarkRequiredFlag(cmd *cobra.Command, flag string) {
	err := cmd.MarkFlagRequired(flag)
	if err != nil {
		panic(err)
	}
}

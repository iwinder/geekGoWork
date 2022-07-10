package week09

import (
	"bytes"
	"encoding/binary"

	"github.com/golang/glog"
)

const (
	DefaultHeader = "Context-Lenght"
	// DefaultBody         = "Body"
	DefaultHeaderLength  = 14
	DefaultBodyLenLength = 4
)

//封包
func Packet(message []byte) []byte {

	return append(append([]byte(DefaultHeader), intToBytes(len(message))...), message...)
}

func UnPacket(messageObj []byte, readerChannel chan []byte) []byte {

	length := len(messageObj)
	i := 0
	for ; i < length; i++ {
		if length < i+DefaultHeaderLength+DefaultBodyLenLength { // 如果当前位置+默认字段占位总长度小于当前内容长度，此时跳过
			break
		}
		if string(messageObj[i:i+DefaultHeaderLength]) == DefaultHeader { // 此时为新消息头开始
			messageLen := bytesToInt(messageObj[i+DefaultHeaderLength : i+DefaultHeaderLength+DefaultBodyLenLength]) // 消息体长度
			if length < i+DefaultHeaderLength+DefaultBodyLenLength+messageLen {                                      // 目前字符串长度小于时，跳过
				break
			}
			data := messageObj[i+DefaultHeaderLength+DefaultBodyLenLength : i+DefaultHeaderLength+DefaultBodyLenLength+messageLen] // 信息
			glog.Infoln("UnPacket data", string(data))
			readerChannel <- data
			i += DefaultHeaderLength + DefaultBodyLenLength + messageLen - 1 // end index
		}
	}
	if i == length { // 没有找到我们约定的头部信息, 返回空切片值
		return make([]byte, 0)
	}
	return messageObj[i:] // 返回消息
}

//initToBytes  int转大端 []byte ，如转小端，需要将 binary.BigEndian->binary.LittleEndian
func intToBytes(n int) []byte {
	i := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, i)
	return bytesBuffer.Bytes()
}

func bytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var i int32
	binary.Read(bytesBuffer, binary.BigEndian, &i)

	return int(i)
}

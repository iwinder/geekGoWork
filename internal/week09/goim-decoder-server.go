package week09

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"io"
	"log"
	"net"

	"github.com/golang/glog"
	decoder "github.com/iwinder/geekGoWork/internal/week09/decoder"
)

// 启动一个TCP服务，根据 type 服务端实现
func RunGoimDecoder(addr string) {
	// 日志打印配置
	flag.Set("v", "5")
	flag.Set("logtostderr", "true")
	flag.Parse()
	defer glog.Flush()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatalln("TCP Server 启动失败...", err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			glog.Fatalln("TCP Server Accept Fail...", err)
		}
		go handle(conn)

	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		peek, err := reader.Peek(decoder.PackageLen)
		if err != nil {
			if err != io.EOF {
				glog.Warning(" 链接被对方关闭", err)
				break
			} else {
				glog.Errorln("链接接收数据异常：", err)
			}
			break
		}
		// 先获取包大小
		buffer := bytes.NewBuffer(peek)
		var size int32
		if err := binary.Read(buffer, binary.BigEndian, &size); err != nil {
			glog.Errorln("链接接收数据异常：", err)
		}
		if int32(reader.Buffered()) < size { // 当数据小于包的设定大小时跳过
			continue
		}
		// 获取包中的内容
		data := make([]byte, size) // 按大小读取包中内容
		if _, err := reader.Read(data); err != nil {
			glog.Errorln("链接接收数据异常：", err)
			continue
		}

		content, err := decoder.Decoder(data) // 解析包的内容
		if err != nil {
			log.Println(err.Error())
			continue
		}
		glog.Infoln("解析数据内容数据：", string(content.Body))
	}
}

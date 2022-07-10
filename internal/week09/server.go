package week09

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"net"

	constant "github.com/iwinder/geekGoWork/internal/week09/constant"
	protocol "github.com/iwinder/geekGoWork/internal/week09/protocol"

	"github.com/golang/glog"
)

// TcpFixLength fix length 服务端实现
func TcpFixLength(conn net.Conn) {
	glog.Infoln("Fix Length Server 接收到数据开始>>>>>>")
	defer conn.Close()
	var buf []byte
	for {
		buf = make([]byte, constant.BYPE_LENGTH)
		_, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				glog.Warning("FixLength 链接被对方关闭", err)
			} else {
				glog.Errorln("FixLength 接收数据异常：", err)
			}

			return
		}
		glog.Infoln("FixLength 接收到数据：", string(buf))
	}
}

// TcpDelimiterBased delimiter based 服务端实现
func TcpDelimiterBased(conn net.Conn) {
	glog.Infoln("Delimiter Based  Server 接收到数据开始>>>>>>")
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		slice, err := reader.ReadSlice('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				glog.Warning("DelimiterBased 链接被对方关闭", err)
			} else {
				glog.Errorln("DelimiterBased 接收数据异常：", err)
			}
			return
		}
		glog.Infoln("DelimiterBased 接收到数据：", string(slice))
	}

}

// length field based frame decoder
func TcpLenField(conn net.Conn) {
	glog.Infoln("Length field based frame decoder Server 接收到数据开始>>>>>>")
	defer conn.Close()
	var readerChannel = make(chan []byte, 32)
	go func() {
		for { // 为快速实现，for循环阻塞，获取所有数据
			select {
			case data := <-readerChannel:
				glog.Infoln("Length field based frame decoder readerChannel 接收到数据：", string(data))
			}
		}
	}()
	buf := make([]byte, 0)
	inputBuffer := make([]byte, 1024)
	for {
		n, err := conn.Read(inputBuffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				glog.Warning("Length field based frame decoder 链接被对方关闭", err)
			} else {
				glog.Errorln("Length field based frame decoder 接收数据异常：", err)
			}
			return
		}
		protocol.UnPacket(append(buf, inputBuffer[:n]...), readerChannel)
	}

}

// 启动一个TCP服务，根据 type 服务端实现
func RunTcpFixLength(addr string, types int) {
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
		switch types {
		case 1:
			go TcpFixLength(conn)
		case 2:
			go TcpDelimiterBased(conn)
		case 3:
			go TcpLenField(conn)
		}

	}
}

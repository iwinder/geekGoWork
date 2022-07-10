package week09

import (
	"flag"
	"math"
	"net"

	"github.com/golang/glog"

	constant "github.com/iwinder/geekGoWork/internal/week09/constant"
	protocol "github.com/iwinder/geekGoWork/internal/week09/protocol"
)

func ClientTcpFixLength(conn net.Conn) {
	glog.Infoln("Fix Length Client 发送数据开始>>>>>>")
	sendBytes := make([]byte, constant.BYPE_LENGTH)
	sendMsg := "{\"userName\":Fix Length,\"userAge\",10}"
	for i := 0; i < 5; i++ {
		tmpBytes := []byte(sendMsg) // 组装本次需要发送的数据
		// 按固定值读取消息并发送
		blen := len(tmpBytes)
		page := round(float64(blen), float64(constant.BYPE_LENGTH)) // 获取一条信息的按 constant.BYPE_LENGTH 拆分次数
		for k := 0; k < page; k++ {
			kdx := k * constant.BYPE_LENGTH
			idx := kdx
			for j := 0; idx < blen && j < constant.BYPE_LENGTH; j++ {
				idx = kdx + j
				sendBytes[j] = tmpBytes[idx]
				idx++
			}
			glog.Infoln("Fix Length Client 数据>>>>>>", string(sendBytes))
			_, err := conn.Write(sendBytes)
			if err != nil {
				glog.Errorln("FixLength 发送数据异常：", err)
				return
			}
		}
		glog.Infoln("Fix Length Client 完成一次>>>>>>")
	}
}

func ClientTcpDelimiterBased(conn net.Conn) {
	glog.Infoln("Delimiter Based  Client 发送数据开始>>>>>>")
	var sendMsgs string
	sendMsg := "{\"userName\":Delimiter Based,\"userAge\",12}\n"
	for i := 0; i < 5; i++ {
		sendMsgs += sendMsg
		_, err := conn.Write([]byte(sendMsgs))
		if err != nil {
			glog.Errorln("Delimiter Based 发送数据异常：", err)
			return
		}
		glog.Infoln("Delimiter Based Client 完成一次>>>>>>")
	}
}
func ClientTcpLenField(conn net.Conn) {
	glog.Infoln("Length field based frame decoder 发送数据开始>>>>>>")
	for i := 0; i < 1000; i++ {
		sendMsg := "{\"userName\":Length field,\"userAge\",22}"
		_, err := conn.Write(protocol.Packet([]byte(sendMsg)))
		if err != nil {
			glog.Errorln("Length field based frame decoder 发送数据异常：", err)
			return
		}
		glog.Infoln("Delimiter Based Client 完成一次>>>>>>")
	}
}

func RunClientTcpFixLength(addr string, types int) {
	// 日志打印配置
	flag.Set("v", "5")
	flag.Set("logtostderr", "true")
	flag.Parse()
	defer glog.Flush()
	lis, err := net.Dial("tcp", addr)
	if err != nil {
		glog.Fatalln("TCP Cilent 启动失败...", err)
	}

	defer lis.Close()
	switch types {
	case 1:
		ClientTcpFixLength(lis)
	case 2:
		ClientTcpDelimiterBased(lis)
	case 3:
		ClientTcpLenField(lis)
	}

}

func round(x float64, y float64) int {
	return int(math.Ceil(x / y))
}

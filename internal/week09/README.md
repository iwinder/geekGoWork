# 第八周作业
## 作业一
> 1. 总结几种 socket 粘包的解包方式：fix length/delimiter based/length field based frame decoder。尝试举例其应用。    

粘包与半包：
- 正常情况服务端与客户端是一条消息一条消息的读取
- 粘包问题：当一方（如客户端）发送两个数据包（如内容是ABC的数据包和内容是DEF数据包）时，另一端（如服务端）一次接收到了这两个数据包且粘合在了一起（其包内容类似 ABCDEF）的问题
- 半包：当一方（如客户端）发送了一个数据包(如内容是ABC)，另一方接收到的是数据包内容只是该包的一部分(如AB)的情况

粘包的主要原因：
- 发送方每次写入数据 < 套接字（Socket）缓冲区大小；
- 接收方读取套接字（Socket）缓冲区数据不够及时。

半包的主要原因：
- 发送方每次写入数据 > 套接字（Socket）缓冲区大小；
- 发送的数据大于协议的 MTU (Maximum Transmission Unit，最大传输单元)，因此必须拆包。

解决方案：
1. 消息定长(即fix length)
   - 控制接收方和传递方每次传递消息的长度为固定值，长度不够时使用空字符弥补。
    - 明显可见该方式增加了不必要的传输，从而增加了网络传输的负担
2. 在包结尾增加特定分隔符（即delimiter based）
    - 双方约定好在每个消息的结尾指定一个特殊字符(如`\nn`)，每次读到该字符算一次消息。
    - 这种方式如果数据量过大，查找定界符会消耗一些性能
3. 封装请求协议，根据里面的长度读取（即基于长度字段的框架解码器，length field based frame decoder）
    - 在TCP的协议基础上再封装一次，即消息头和消息体的形式，消息头中包含消息的总长度（或消息体的总长度），接收方读取到该值后，便可知道这条消息的具体边界。
    - 这种方式意味着增加了一次封包/拆包的操作，在某种程度上造成一定的延迟等问题。

代码实现:
- [server.go](server.go) 为服务端实现
    - `func TcpFixLength(conn net.Conn)` 为对方案一的实现
        ```Go
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
        ```
    - `func TcpDelimiterBased(conn net.Conn) ` 为对方案二的实现
        ```Go
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
        ```
    - `func TcpLenField(conn net.Conn) `为对方案三的实现
        ```Go
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
        ```
    - 可通过执行[cmd/week09/server/server.go](../../cmd/week09/server/server.go) 启动服务端，通过指定参数 types=1/2/3 进行方案切换
        ```Go
            func main() {
                // // 方案一
                // week09.RunTcpFixLength("127.0.0.1:8033", 1)
                // // 方案二
                // week09.RunTcpFixLength("127.0.0.1:8033", 2)
                // 方案三
                week09.RunTcpFixLength("127.0.0.1:8033", 3)
            }
        ```
- [cilent.go](cilent.go) 为客户端实现
    - 相关实现对照服务端，前缀增加了`Client`
    -  可通过执行[cmd/week09/client/client.go](../../cmd/week09/client/client.go) 启动服务端
- [protocol/protocol.go](protocol/protocol.go) 为方案三中具体封包和拆包逻辑

## 作业二
> 2. 实现一个从 socket connection 中解码出 goim 协议的解码器。

根据视频和PPT可知协议结构为：
|Type |	Size |
| --- | --- |
| Package Length，包长度 |  4 bytes |
| Header Length，头长度 | 2 bytes |
| Protocol Version，协议版本 | 2 bytes |
| Operation，操作码 |   4 bytes |
| Sequence 请求序号 ID |    4 bytes |
| Body，包内容 | PackLen-HeaderLen |

项目结构：
- [goim-decoder-server.go](goim-decoder-server.go) 为服务端 TCP 创建与运行逻辑
- [/cmd/week09/decoder/decoder.go](../..//cmd/week09/decoder/decoder.go) 为启动脚本
- 简单解析实现方案在[/decoder/goim-decoder.go](decoder/goim-decoder.go)

```Go
func Decoder(msg []byte) (*Pack, error) {
	if len(msg) < SeqOff+1 { // 如果包的byte 无法取到Body部分，则认为是不完整的包
		return nil, ErrorOFIncomplete
	}
	// 取前4个获得 Package Length，包长度
	packageLength := binary.BigEndian.Uint32(msg[0:PackageOff])
	// // 第 4 - 6 是 Header Length，头长度
	headerLength := binary.BigEndian.Uint16(msg[PackageOff:HeaderLenOff])
	// 第 6 - 8 是版本号
	version := binary.BigEndian.Uint16(msg[HeaderLenOff:VersionOff])
	// 8 - 12 是 operation
	operation := binary.BigEndian.Uint32(msg[VersionOff:OperOff])
	// 12 - 16 是 seqId
	seqId := binary.BigEndian.Uint32(msg[OperOff:SeqOff])

	// BodyLen =  PackLenght-HeaderLen
	// backLen := int32(packageLength) - int32(headerLength)
	content := msg[SeqOff:]
	return &Pack{
		PackLenght: int32(packageLength),
		HeaderLen:  int16(headerLength),
		Ver:        int16(version),
		Op:         int32(operation),
		Seq:        int32(seqId),
		Body:       content,
	}, nil
}
```
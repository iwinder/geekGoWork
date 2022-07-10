package decoder

import (
	"encoding/binary"
	"errors"
)

const (
	PackageLen         = 4
	HeaderLen          = 2
	ProtocolVersionLen = 2
	OperationLen       = 4
	SequenceLen        = 4

	PackageOff   = PackageLen
	HeaderLenOff = PackageOff + HeaderLen
	VersionOff   = HeaderLenOff + ProtocolVersionLen
	OperOff      = VersionOff + OperationLen
	SeqOff       = OperOff + SequenceLen
)

// 定义包不完整异常
var ErrorOFIncomplete = errors.New("the  package is incomplete")

type Pack struct {
	PackLenght int32 // 这个包的长度，一条消息的长度，要考虑大小端
	HeaderLen  int16 // 协议中原生header长度
	Ver        int16
	Op         int32
	Seq        int32
	Body       []byte // 长度为 PackLenght-HeaderLen
}

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

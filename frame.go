package websocket

const (
	OPCODE_CONTINUATION = 0x0
	OPCODE_FOR_TEXT     = 0x1
	OPCODE_BINARY       = 0x2
	OPCODE_CLOSE        = 0x8
	OPCODE_PING         = 0x9
	OPCODE_PONG         = 0xA
)

type Frame struct {
	fin           byte
	opcode        byte
	payloadLength int
	mask          []byte
	payload       []byte
}

func (f Frame) Text() string {
	return string(f.payload)
}

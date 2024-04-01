package websocket

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	bufferSize = 4096
)

type Connection interface {
	Close() error
}

type Websocket struct {
	conn   Connection
	bufrw  *bufio.ReadWriter
	header http.Header
}

func New(w http.ResponseWriter, r *http.Request) (*Websocket, error) {
	hj, ok := w.(http.Hijacker)

	if !ok {
		return nil, errors.New("Hijacker를 지원하지 않습니다.")
	}

	conn, bufrw, err := hj.Hijack()

	if err != nil {
		return nil, err
	}

	return &Websocket{conn, bufrw, r.Header}, nil
}

func makeAcceptHash(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	h.Write([]byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (ws *Websocket) Handshake() error {
	hash := makeAcceptHash(ws.header.Get("Sec-WebSocket-Key"))
	lines := []string{
		"HTTP/1.1 101 Web Socket Protocol Handshake",
		"Server: acky_websocket_server",
		"Upgrade: WebSocket",
		"Connection: Upgrade",
		"Sec-WebSocket-Accept: " + hash,
		"",
		"",
	}
	return ws.write([]byte(strings.Join(lines, "\r\n")))
}

func (ws *Websocket) write(p []byte) error {
	if _, err := ws.bufrw.Write(p); err != nil {
		return err
	}

	return ws.bufrw.Flush()
}

func (ws *Websocket) read() ([]byte, error) {
	data := make([]byte, bufferSize)

	n, err := ws.bufrw.Read(data)
	if err != nil {
		return nil, err
	}

	return data[:n], nil
}

func (ws *Websocket) Recv() (Frame, error) {
	frame := Frame{}
	head, err := ws.read()
	if err != nil {
		return frame, err
	}

	length := uint64(head[1] & 0x7F)

	maskIndex := 2

	if length == 126 {
		byteArr := []byte{head[2], head[3]}
		length = uint64(binary.BigEndian.Uint16(byteArr))
		maskIndex = 4
	} else if length == 127 {
		byteArr := []byte{head[2], head[3], head[4], head[5]}
		length = uint64(binary.BigEndian.Uint64(byteArr))
		maskIndex = 6
	}

	frame.fin = head[0] & 0x80
	frame.opcode = head[0] & 0x0F
	frame.payloadLength = int(length)
	frame.mask = head[maskIndex : maskIndex+4]

	frame.payload = head[maskIndex+4:]

	for i := 0; i < len(frame.payload); i++ {
		frame.payload[i] ^= frame.mask[i%len(frame.mask)]
	}
	return frame, nil
}

func (ws *Websocket) Send(f Frame) error {
	data := make([]byte, 2)
	data[0] = 0x80 | f.opcode

	if f.payloadLength <= 125 {
		data[1] = byte(f.payloadLength)
		data = append(data, f.payload...)
	} else if f.payloadLength > 125 {
		data[1] = byte(126)
		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, uint16(f.payloadLength))
		data = append(data, size...)
		data = append(data, f.payload...)
	} else {
		data[1] = byte(127)
		size := make([]byte, 8)
		binary.BigEndian.PutUint64(size, uint64(f.payloadLength))
		data = append(data, size...)
		data = append(data, f.payload...)
	}
	fmt.Println(string(data))
	return ws.write(data)
}

func (ws *Websocket) Close() error {
	return ws.conn.Close()
}

package websocket

import (
	"fmt"
	"log"
	"net/http"
)

func Run() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws, err := New(w, r)
		if err != nil {
			fmt.Println(err)
		}

		// Handshake로 ws 연결유지
		err = ws.Handshake()
		if err != nil {
			fmt.Println(err)
		}

		for {
			frame, err := ws.Recv()

			if err != nil {
				// 웹소켓에서 페이로드를 못받아옴
				fmt.Println(err)
				break
			}

			switch frame.opcode {
			case OPCODE_CLOSE:
				return
			case OPCODE_PING:
				frame.opcode = 10
				fallthrough
			case OPCODE_BINARY, OPCODE_CONTINUATION, OPCODE_FOR_TEXT:
				if err = ws.Send(frame); err != nil {
					log.Println(err)
					return
				}
			}
		}

	})

	http.ListenAndServe(":5002", nil)
}

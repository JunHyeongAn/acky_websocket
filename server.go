package websocket

import (
	"fmt"
	"log"
	"net/http"
)

func Run(path string, port string) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		ws, err := New(w, r)
		if err != nil {
			fmt.Println(err)
		}

		err = ws.Handshake()
		if err != nil {
			fmt.Println(err)
		}

		for {
			frame, err := ws.Recv()

			if err != nil {
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

	http.ListenAndServe(":"+port, nil)
}

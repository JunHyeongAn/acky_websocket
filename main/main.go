package main

import (
	"fmt"

	"github.com/JunHyeongAn/websocket"
)

func main() {
	websocket.Run("/", "5050", func(ws *websocket.Websocket) {
		for {
			frame, err := ws.Recv()

			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Println(frame.Text())
		}
	})
}

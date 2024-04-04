package websocket

import (
	"fmt"
	"net/http"
)

func Run(path string, port string, socketHandler func(ws *Websocket)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		ws, err := New(w, r)
		if err != nil {
			fmt.Println(err)
		}
		socketHandler(ws)
	})

	http.ListenAndServe(":"+port, nil)
}

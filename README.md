
# Acky_Websocket
[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
Acky_Websocket은 GO로 WebSocker을 구현한 오픈소스입니다.




## Installation

Install my-project with npm

```bash
  go get github.com/JunHyeongAn/acky_websocket
```
    
## Usage/Examples

```Go
import (
	"fmt"

	"github.com/JunHyeongAn/websocket"
)

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
```


// websocket
// @author xiangqian
// @date 23:33 2022/12/19
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// https://stackoverflow.com/questions/59294893/runtime-error-while-trying-to-connect-to-websocket-with-the-gorilla-library
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(pReq *http.Request) bool { return true },
}

func Ws(pContext *gin.Context) {
	// 服务升级，http请求升级为websocket请求
	pConn, err := upgrader.Upgrade(pContext.Writer, pContext.Request, nil)
	if err != nil {
		panic(err)
	}
	defer pConn.Close()

	for {
		messageType, buf, err := pConn.ReadMessage()
		if err != nil {
			log.Printf("%v\n", err)
			return
		}

		switch messageType {
		case websocket.TextMessage:
			log.Printf("recv: %s", string(buf))
			if err = pConn.WriteMessage(websocket.TextMessage, buf); err != nil {
				log.Printf("%v \n", err)
				break
			}

		case websocket.BinaryMessage:
			log.Printf("Binary messages are not supported.\n")

		default:
			log.Printf("%v messages are not supported.\n", messageType)
		}
	}
}

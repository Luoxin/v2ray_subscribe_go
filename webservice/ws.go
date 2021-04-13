package webservice

import (
	"compress/flate"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/websocket/v2"
	log "github.com/sirupsen/logrus"
)

type WsClient struct {
	c *websocket.Conn

	clientId string
}

func NewWsClient(c *websocket.Conn) *WsClient {
	return &WsClient{
		c: c,
	}
}

func (p *WsClient) Init(clientId string) error {
	p.clientId = clientId

	c := p.c

	err := c.SetCompressionLevel(flate.BestCompression)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	c.SetCloseHandler(func(code int, text string) error {
		log.Warnf("client %s coon closed, code:%d, msg:%s", clientId, code, text)
		wsClientDispatch.OnClose(clientId)
		return nil
	})

	c.SetPingHandler(func(appData string) error {
		return p.SendMsg(map[string]interface{}{
			"ping": time.Now().UnixNano(),
		})
	})

	c.SetPongHandler(func(appData string) error {
		return p.SendMsg(map[string]interface{}{
			"pong": time.Now().UnixNano(),
		})
	})

	return nil
}

func (p *WsClient) SendMsg(msg interface{}) error {
	return p.c.WriteJSON(msg)
}

func (p *WsClient) RecvMsgForever() {
	c := p.c
	for {
		_, recvData, err := c.ReadMessage()
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

		log.Infof("rcv data:%s", string(recvData))
	}
}

type WsClientDispatch struct {
	clientsMap map[string]*WsClient
	clientLock sync.RWMutex
}

func NewWsClientDispatch() *WsClientDispatch {
	return &WsClientDispatch{
		clientsMap: make(map[string]*WsClient),
	}
}

var wsClientDispatch = NewWsClientDispatch()

func (p *WsClientDispatch) Add(c *websocket.Conn) error {
	clientId := utils.UUIDv4()

	w := NewWsClient(c)

	err := w.Init(clientId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	p.clientLock.Lock()
	p.clientsMap[clientId] = w
	p.clientLock.Unlock()

	w.RecvMsgForever()

	return nil
}

func (p *WsClientDispatch) OnClose(clientId string) {
	p.clientLock.Lock()
	delete(p.clientsMap, clientId)
	p.clientLock.Unlock()
}

func InitWs(app *fiber.App) error {
	app.Use("ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	app.Get("ws/", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index

		err := wsClientDispatch.Add(c)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

		// for {
		// 	if mt, msg, err = c.ReadMessage(); err != nil {
		// 		log.Println("read:", err)
		// 		break
		// 	}
		// 	log.Printf("recv: %s", msg)
		//
		// 	if err = c.WriteMessage(mt, msg); err != nil {
		// 		log.Println("write:", err)
		// 		break
		// 	}
		// }

	},
		websocket.Config{
			Filter:            nil,
			HandshakeTimeout:  time.Second * 5,
			Subprotocols:      nil,
			Origins:           nil,
			EnableCompression: true,
		},
	))

	return nil
}

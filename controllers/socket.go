package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type MyWebSocketController struct {
	beego.Controller
}

var upgrader = websocket.Upgrader{}

type Message struct {
	Message string `json:"message"`
}
type SendMsg struct {
	Ws *websocket.Conn
	Msg Message
}
var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan SendMsg,10)
)

func init() {
	go handleMessages()
}

//广播发送至页面
func handleMessages() {
	for {
		msg := <-broadcast
		//fmt.Println("clients len ", len(clients))
		err := msg.Ws.WriteJSON(msg.Msg)
		beego.Info("ws",msg.Ws)
		if err != nil {
				beego.Error("client.WriteJSON error: %v", err)
				msg.Ws.Close()
				delete(clients, msg.Ws)
		}
		//for client := range clients {
		//	beego.Info("\nclients len ", len(clients))
		//	beego.Info("client  RemoteAddr", client.RemoteAddr())
		//	beego.Info("client  LocalAddr\n", client.LocalAddr())
		//	err := client.WriteJSON(msg)
		//	if err != nil {
		//		beego.Error("client.WriteJSON error: %v", err)
		//		client.Close()
		//		delete(clients, client)
		//	}
		//}
	}
}

func (c *MyWebSocketController) Get() {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		beego.Error(err)
	}
	//  defer ws.Close()

	clients[ws] = true

	beego.Info("ws",ws)
	//不断的广播发送到页面上
	for {
		//目前存在问题 定时效果不好 需要在业务代码替换时改为beego toolbox中的定时器
		time.Sleep(time.Second * 3)

		msg := Message{Message: "这是向页面发送的数据 " + time.Now().Format("2006-01-02 15:04:05")}
		ss := SendMsg{Ws:ws,Msg:msg}
		broadcast <- ss
	}
}
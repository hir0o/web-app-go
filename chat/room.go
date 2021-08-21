package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	forward chan []byte // クライアントに送信するメッセージを保持するチャネル
 	join chan *client // チャットルームに参加しようとしているクライアントのチャネル
	leave chan *client // チャットルームから退室しようとしているクライアントのチャネル
	clients map[*client]bool // 在室しているクライアントが保持される
}

func newRoom() *room {
	return &room {
		forward: make(chan []byte),
		join: make(chan *clinet),
		leave: make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
  for {
    select {
    case client := <- r.join:
      // 参加
      r.clients[client] = true
    case client := <- r.leave:
      // 退室
      delete(r.clients, client)
      close(client.send)
    case msg := <- r.forward:
      // すべてのクライアントにメッセージを送信
      for client := range r.clients {
        select {
          case client.send <- msg:
            // メッセージを送信
          default:
            // 送信に失敗
            delete(r.clients, client)
            close(client.send)
        }
      }
    }
  }
}

const (
  socketBufferSize = 1024
  messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
  ReadBufferSize: socketBufferSize,
  WriteBufferSize: socketBufferSize,
}

func(r * room) ServHTTP(w http.ResponseWriter, req *http.Request) {
  // Upgradeでhttpリクエストから、WebSocketコネクションを取得する
  socket, err := upgrader.Upgrade(w, req, nil)
  if err != nil {
    log.Fatal("ServeHTTP:", err)
    return
  }
  // socketが取得できたらclientを生成する
  client := &client{
    socket: socket,
    send: make(chan []byte, messageBufferSize),
    room: r,
  }
  // joinチャネルに渡す
  r.joim <- client
  // clietoの終了時にleaveチャネルに渡す
  defer func() { r.leave <- client }()
  // 別のスレッドで呼び出す
  go client.write()
  // 接続を保持するため
  client.read()
}


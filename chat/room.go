package main

import (
	"log"
	"net/http"
	"fmt"
	"github.com/gorilla/websocket"
  "../trace"
)

type room struct {
	forward chan []byte // クライアントに送信するメッセージを保持するチャネル
 	join chan *client // チャットルームに参加しようとしているクライアントのチャネル
	leave chan *client // チャットルームから退室しようとしているクライアントのチャネル
	clients map[*client]bool // 在室しているクライアントが保持される
  tracer trace.Tracer
}

func newRoom() *room {
	return &room {
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
    tracer: trace.Off(),
	}
}

func (r *room) run() {
  for {
    select {
    case client := <- r.join:
      // 参加
	  fmt.Println("参加した")
      r.clients[client] = true
      r.tracer.Trace("新しいクライアントが参加しました")
    case client := <- r.leave:
      // 退室
	  fmt.Println("退室した")
      delete(r.clients, client)
      close(client.send)
    case msg := <- r.forward:
      // すべてのクライアントにメッセージを送信
	  fmt.Println("メッセージの送信")
      for client := range r.clients {
        select {
          case client.send <- msg:
			fmt.Println("メッセージの送信")
            // メッセージを送信
          default:
			fmt.Println("送信に失敗")
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

func(r * room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
  r.join <- client
  // clietoの終了時にleaveチャネルに渡す
  defer func() { r.leave <- client }()
  // 別のスレッドで呼び出す
  go client.write()
  // 接続を保持するため
  client.read()
}


package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn // websocketへの参照
	send chan []byte // バッファつきチャネル.受信したメッセージが蓄積され、WebSocketでユーザーに送られるのを待機している.
	room *room  // ルームへの参照ルーム全体に送信するときに使用
}

// ReadMessageを使って、データの読み込みに使われる
func (c *client) read() {
  for {
    if _, msg, err := c.socket.ReadMessage(); err == nil {
      // 受け取ったメッセージはすぐにroomのforewardチャネルに送られる
      c.room.forward <- msg
    } else {
      // WebSocketで異常が起きたら、closeする
      break
    }
  }
  c.socket.Close()
}

func (c *client) write() {
  for msg := range c.send {
    if err := c.socket.WriteMessage(websocket.TextMessage, msg);
      err != nil {
        break
      }
  }
  c.socket.Close()
}

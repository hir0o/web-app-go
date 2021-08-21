package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
    "flag"
)

type templateHandler struct {
    once         sync.Once
    filename     string
    templ        *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    t.once.Do(func() {
        t.templ = 
            template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
    })
    t.templ.Execute(w, r)
}

func main() {
    var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
    flag.Parse() // フラグを解釈
    r := newRoom()
    // ルート
    // テンプレートのコンパイルは一度だけしか行われない。
    http.Handle("/", &templateHandler{filename: "chat.html"})
    http.Handle("/room", r)
    // チャットルームの開始
    go r.run()
    // webサーバの開始
    log.Println("Webサーバーを開始します。ポート: ", *addr)
    if err := http.ListenAndServe(*addr, nil); err != nil {
        log.Fatal("ListenAndServer:", err)
    }
}

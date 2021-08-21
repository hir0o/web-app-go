package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
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
    t.templ.Execute(w, nil)
}

func main() {
    r := newRoom()
    // ルート
    // テンプレートのコンパイルは一度だけしか行われない。
    http.Handle("/", &templateHandler{filename: "chat.html"})
    http.Hndle("/room", r)
    // チャットルームの開始
    go r.run()
    // webサーバの開始
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndServer:", err)
    }
}

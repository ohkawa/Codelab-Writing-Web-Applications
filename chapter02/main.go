// Codelab: Webアプリケーションを書いてみよう - golang.jp
// http://golang.jp/codelab-wiki
//
// サーバ立ち上げ
// $ go run main.go
//
// ブラウザでアクセス
// http://localhost:8080/monkeys

package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

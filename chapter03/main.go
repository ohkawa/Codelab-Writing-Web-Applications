// Codelab: Webアプリケーションを書いてみよう - golang.jp
// http://golang.jp/codelab-wiki
//
// test.txtがなければ実行前に作成
// echo "Hello world" > test.txt
//
// サイトには実行コマンド'$ 8g wiki.go'とありますが、現在は以下で実行可能。
// $ go run main.go
//
// ブラウザでアクセス
// http://localhost:8080/view/test

package main

import (
	"io/ioutil"
	"fmt"
	"net/http"
)

type Page struct {
	Title	string
	Body	[]byte
}

const lenPath = len("/view/")

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.ListenAndServe(":8080", nil)
}

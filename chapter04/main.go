// Codelab: Webアプリケーションを書いてみよう - golang.jp
// http://golang.jp/codelab-wiki
// 「wikiページを動かすためにhttpを使用する」 ~ 「templateパッケージ」まで
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
	"net/http"
	"html/template"
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
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[lenPath:]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl+".html") // サイトでは'ParseFile'になっているので注意!(sが足りない)

	t.Execute(w, p)
}
func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
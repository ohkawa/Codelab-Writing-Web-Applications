// Codelab: Webアプリケーションを書いてみよう - golang.jp
// http://golang.jp/codelab-wiki
// 「関数リテラルとクロージャの手ほどき」
// http.HandlerFuncを使用したリファクタリングのため、動作には変更なし
//
// 実行
// $ go run main.go
//
// ブラウザでアクセス
// http://localhost:8080/view/ANewPage
// http://localhost:8080/edit/ANewPage
//

package main

import (
	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
	"errors"
)

type Page struct {
	Title	string
	Body	[]byte
}

const lenPath = len("/view/")

var templates = make(map[string]*template.Template)
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func init() {
	for _, tmpl := range []string{"edit", "view"} {
		t := template.Must(template.ParseFiles(tmpl + ".html"))  // サイトでは'ParseFile'になっているので注意!(sが足りない)
		templates[tmpl] = t
	}
}

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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates[tmpl].Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 'err.Strings()'を'err.Error()'へ書き換え
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError) // 'err.Strings()'を'err.Error()'へ書き換え
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[lenPath:]
		if !titleValidator.MatchString(title) {
			http.NotFound(w, r)
			return
		}
		fn(w, r, title)
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}
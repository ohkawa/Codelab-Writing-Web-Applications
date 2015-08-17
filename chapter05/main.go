// Codelab: Webアプリケーションを書いてみよう - golang.jp
// http://golang.jp/codelab-wiki
// 「存在しないページのハンドリング」 ~ 「入力チェック」まで
//
// 実行
// $ go run main.go
//
// ブラウザでアクセス
// http://localhost:8080/view/test
// http://localhost:8080/edit/test
//
// http://localhost:8080/ は 404 not found(エラーハンドリング)

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

func getTitle(w http.ResponseWriter, r *http.Request) (title string, err error) { // 'os.Error'を'error'へ書き換え
	title = r.URL.Path[lenPath:]
	if !titleValidator.MatchString(title) {
		http.NotFound(w, r)
		err = errors.New("Invalid Page Title") // 'os.NewError'を'errors.New'へ書き換え
	}
	return
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

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError) // 'err.Strings()'を'err.Error()'へ書き換え
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
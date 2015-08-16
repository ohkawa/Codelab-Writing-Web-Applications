// Codelab: Webアプリケーションを書いてみよう - golang.jp
// http://golang.jp/codelab-wiki
//
// サイトには実行コマンド'$ 8g wiki.go'とありますが、現在は以下で実行可能。
// $ go run main.go
// 実行後'TestPage.txt'が生成される

package main

import (
//	"os"           // 'os.Error'が必要なくなったためosパッケージのimportも不要
	"io/ioutil"
	"fmt"
)

type Page struct {
	Title	string
	Body	[]byte
}

// 'os.Error'はなくなったため'error'へ書き換え
// func (p *Page) save() os.Error {
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// func (p *Page) save() os.Error {
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
}

// 版权 @2023 凹语言 作者。保留所有权利。

package appplay

import (
	_ "embed"
	"html/template"
	"net/http"
)

func (p *WebServer) editHandler(w http.ResponseWriter, r *http.Request) {
	snip := &Snippet{Body: []byte(edit_helloPlayground)}
	edit_Template.Execute(w, &editData{snip})
}

//go:embed _edit.tmpl.html
var edit_tmpl string

var edit_Template = template.Must(template.New("playground/index.html").Parse(edit_tmpl))

const edit_helloPlayground = `// 版权 @2019 凹语言 作者。保留所有权利。

import "fmt"
import "runtime"

global year: i32 = 2023

func main {
	println("你好，凹语言！", runtime.WAOS)
	println(add(40, 2), year)

	fmt.Println("1+1 =", 1+1)
}

func add(a: i32, b: i32) => i32 {
	return a+b
}
`

// 版权 @2023 凹语言 作者。保留所有权利。

package appplay

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

type WebServer struct {
	fs fs.FS
}

//go:embed _static
var fs_static embed.FS

func getStaticFS() fs.FS {
	if true {
		//	return fs_static
	}
	fs, err := fs.Sub(fs_static, "_static")
	if err != nil {
		panic(err)
	}
	return fs
}

func NewWebServer() *WebServer {
	p := &WebServer{
		fs: getStaticFS(),
	}

	return p
}

func (p *WebServer) Run(addr string) error {
	startTime := time.Now()
	return http.ListenAndServe(addr,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.Method, r.URL.Path)

			switch {
			case r.URL.Path == "/":
				p.editHandler(w, r)
			case r.URL.Path == "/-/play/compile":
				p.runHandler(w, r)
			case r.URL.Path == "/-/play/fmt":
				p.fmtHandler(w, r)
			case strings.HasPrefix(r.URL.Path, "/static/"):
				relpath := strings.TrimPrefix(r.URL.Path, "/static/")
				data, err := fs.ReadFile(p.fs, relpath)
				if err != nil {
					http.NotFound(w, r)
					return
				}

				http.ServeContent(w, r, r.URL.Path, startTime, bytes.NewReader(data))

			default:
				p.editHandler(w, r)
			}
		}),
	)
}

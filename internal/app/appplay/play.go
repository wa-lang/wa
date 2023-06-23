// 版权 @2023 凹语言 作者。保留所有权利。

package appplay

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func RunPlayground(addr string) error {
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}
	fmt.Printf("listen at http://%s\n", addr)

	go func() {
		time.Sleep(time.Second * 2)
		openBrowser(addr)
	}()

	s := NewWebServer()
	return s.Run(addr)
}

func openBrowser(url string) error {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}

// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appplay

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"wa-lang.org/wa/internal/3rdparty/cli"
)

var CmdPlay = &cli.Command{
	Name:  "play",
	Usage: "start Wa playground",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "http",
			Value: ":2023",
			Usage: "set http address",
		},
	},
	Action: func(c *cli.Context) error {
		err := RunPlayground(c.String("http"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}

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

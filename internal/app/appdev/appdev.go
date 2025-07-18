// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appdev

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"wa-lang.org/wa/api"
	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/config"
)

var CmdDev = &cli.Command{
	Hidden: true,
	Name:   "dev",
	Usage:  "only for dev/debug",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "hello",
		},
		&cli.BoolFlag{
			Name: "count-code-lines",
		},
	},
	Action: func(c *cli.Context) error {
		log.SetFlags(log.Llongfile)

		if c.Bool("count-code-lines") {
			RunCountCodeLines(c)
			os.Exit(0)
		}

		if c.Bool("hello") {
			_, wat, _, err := api.BuildFile(
				config.DefaultConfig(),
				"hello.wa", "func main() { println(123) }",
			)
			if err != nil {
				if len(wat) != 0 {
					fmt.Println(string(wat))
				}
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(string(wat))
			os.Exit(0)
		}

		fmt.Println("...dev...")
		return nil
	},
}

func RunCountCodeLines(c *cli.Context) {
	var dir = "."
	if c.NArg() > 0 {
		dir = c.Args().First()
	}

	total := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal("filepath.Walk: ", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatal("ioutil.ReadFile: ", err)
		}
		if needSkip(path, data) {
			return nil
		}

		n := countLine(data)
		fmt.Println(path, n)

		total += n
		return nil
	})
	fmt.Printf("total %d\n", total)
}

func needSkip(path string, data []byte) bool {
	if !hasExt(path, ".wa", ".go", ".ws") {
		return true
	}
	if strings.Contains(path, "3rdparty") {
		return true
	}
	if bytes.Contains(data, []byte("The Go Authors")) {
		return true
	}

	return false
}

func countLine(data []byte) int {
	return bytes.Count(data, []byte("\n"))
}

func hasExt(name string, extensions ...string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(name, ext) {
			return true
		}
	}
	return false
}

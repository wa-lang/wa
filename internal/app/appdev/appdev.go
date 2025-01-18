package appdev

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"wa-lang.org/wa/api"
	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/waroot/malloc"
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
			Name: "malloc",
		},
		&cli.BoolFlag{
			Name: "count-code-lines",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("count-code-lines") {
			RunCountCodeLines()
			os.Exit(0)
		}

		if c.Bool("hello") {
			_, wat, err := api.BuildFile(
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

		if c.Bool("malloc") {
			h := malloc.NewHeap(&malloc.Config{
				MemoryPages:    1,
				MemoryPagesMax: 2,
				StackPtr:       100,
				HeapBase:       1000,
				HeapLFixedCap:  3,
			})

			os.WriteFile("a.out-0.wasm", []byte(h.WasmBytes()), 0666)

			os.WriteFile("a.out-0.dot", []byte(h.DotString()), 0666)

			// 需要页扩展时失败
			p1 := h.Malloc(65536 - 10) // (malloc.KPageBytes)
			os.WriteFile("a.out-1.dot", []byte(h.DotString()), 0666)
			fmt.Println("p1:", p1)
			h.Free(p1)

			os.WriteFile("a.out-2.dot", []byte(h.DotString()), 0666)

			p2 := h.Malloc(malloc.KPageBytes)
			os.WriteFile("a.out-2.dot", []byte(h.DotString()), 0666)
			fmt.Println("p2:", p2)

			os.Exit(0)
		}

		fmt.Println("...dev...")
		return nil
	},
}

func RunCountCodeLines() {
	var dir = "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
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
		data, err := ioutil.ReadFile(path)
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

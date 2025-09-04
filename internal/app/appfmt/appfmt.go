// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appfmt

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbase"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/format"
	"wa-lang.org/wa/internal/native/abi"
	natfmt "wa-lang.org/wa/internal/native/format"
	"wa-lang.org/wa/internal/wat/watutil/watfmt"
)

var CmdFmt = &cli.Command{
	Name:      "fmt",
	Usage:     "format Wa source code file",
	ArgsUsage: "[<file.wa>|<path>|<path>/...]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "riscv",
			Usage: "set riscv cpu type for assembly code",
		},
	},
	Action: func(c *cli.Context) error {
		for _, path := range c.Args().Slice() {
			if err := Fmt(c, path); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		return nil
	},
}

func Fmt(c *cli.Context, path string) error {
	if appbase.IsNativeFile(path, ".wat") {
		return fmtWatFile(path)
	}
	if appbase.IsNativeFile(path, ".wa", ".wz") {
		_, err := fmtFile(path)
		return err
	}
	if appbase.IsNativeFile(path, ".s") {
		err := fmtNativeAsmFile(c, path)
		return err
	}

	if path == "" {
		if _, err := os.Lstat(config.WaModFile); err == nil {
			path = "./..."
		} else if _, err := os.Lstat(config.WaModFile + ".json"); err == nil {
			path = "./..."
		} else {
			path = "."
		}
	}
	if !filepath.IsAbs(path) && !strings.HasPrefix(path, ".") {
		return fmt.Errorf("%q is not valid path", path)
	}

	var waFileList []string
	switch {
	case strings.HasSuffix(path, ".wa"):
		waFileList = append(waFileList, path)
	case strings.HasSuffix(path, ".wz"):
		waFileList = append(waFileList, path)
	case strings.HasSuffix(path, "/..."):
		waFileList = getDirWaFileList(
			strings.TrimSuffix(path, "..."),
			true, ".wa", ".wz", // 包含子目录
		)
	default:
		// 不包含子目录
		waFileList = getDirWaFileList(
			path, false, ".wa", ".wz",
		)
	}

	var changedFileList []string
	for _, s := range waFileList {
		changed, err := fmtFile(s)
		if err != nil {
			return fmt.Errorf("%s: %w", s, err)
		}
		if changed {
			changedFileList = append(changedFileList, s)
		}
	}
	for _, s := range changedFileList {
		fmt.Println(s)
	}
	return nil
}

func fmtFile(path string) (changed bool, err error) {
	code, changed, err := format.File(nil, path, nil)
	if err != nil {
		return false, err
	}
	if changed {
		if err = os.WriteFile(path, code, 0666); err != nil {
			return false, err
		}
	}
	return changed, nil
}

func fmtWatFile(path string) (err error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	data, err := watfmt.Format(path, src)
	if err != nil {
		return err
	}
	os.Stdout.Write(data)
	if !bytes.HasSuffix(data, []byte("\n")) {
		os.Stdout.Write([]byte("\n"))
	}
	return nil
}

func fmtNativeAsmFile(c *cli.Context, path string) (err error) {
	if !c.Bool("riscv") {
		return fmt.Errorf("only support ricv type")
	}
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	data, err := natfmt.Format(abi.RISCV64, path, src)
	if err != nil {
		return err
	}
	os.Stdout.Write(data)
	if !bytes.HasSuffix(data, []byte("\n")) {
		os.Stdout.Write([]byte("\n"))
	}
	return nil
}

func getDirWaFileList(dir string, walkSubDir bool, extList ...string) []string {
	var waFileList []string
	if !walkSubDir {
		files, err := os.ReadDir(".")
		if err != nil {
			return nil
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			for _, ext := range extList {
				if strings.HasSuffix(file.Name(), ext) {
					waFileList = append(waFileList, filepath.Join(dir, file.Name()))
				}
			}
		}

		sort.Strings(waFileList)
		return waFileList
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, ext := range extList {
			if strings.HasSuffix(path, ext) {
				waFileList = append(waFileList, path)
				return nil
			}
		}
		return nil
	})
	sort.Strings(waFileList)
	return waFileList
}

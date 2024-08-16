// 版权 @2024 凹语言 作者。保留所有权利。

//go:build wasm
// +build wasm

// wa 命令迷你版本, 支持 WASM-wasi 版本命令行构建, 用于 web 插件环境.
package main

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appbuild"
	"wa-lang.org/wa/internal/app/appfmt"
	"wa-lang.org/wa/internal/app/appinit"
	"wa-lang.org/wa/internal/app/applsp"
	"wa-lang.org/wa/internal/app/apprun"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/version"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "Wa"
	cliApp.Usage = "Wa is a tool for managing Wa source code."
	cliApp.Copyright = "Copyright 2018 The Wa Authors. All rights reserved."
	cliApp.Version = version.Version
	cliApp.EnableBashCompletion = true
	cliApp.HideHelpCommand = true

	cliApp.CustomAppHelpTemplate = cli.AppHelpTemplate +
		"\nSee \"https://wa-lang.org\" for more information.\n"

	cliApp.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "set debug mode",
		},
		&cli.StringFlag{
			Name:    "trace",
			Aliases: []string{"t"},
			Usage:   "set trace mode (*|app|compiler|loader)",
		},
	}

	cliApp.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			config.SetDebugMode()
		}
		if parten := c.String("trace"); parten != "" {
			config.SetEnableTrace(parten)
		}
		return nil
	}

	// 没有参数时显示 help 信息
	cliApp.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			fmt.Println("unknown command:", strings.Join(c.Args().Slice(), " "))
			os.Exit(1)
		}
		cli.ShowAppHelpAndExit(c, 0)
		return nil
	}

	cliApp.Commands = []*cli.Command{
		appinit.CmdInit,
		appbuild.CmdBuild,
		apprun.CmdRun,
		appfmt.CmdFmt,
		applsp.CmdLsp, // hidden
	}

	cliApp.Run(os.Args)
}

// 版权 @2019 凹语言 作者。保留所有权利。

//go:build !wasm
// +build !wasm

// 凹语言，The Wa Programming Language.
package app

import (
	"fmt"
	"os"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/app/appast"
	"wa-lang.org/wa/internal/app/appbuild"
	"wa-lang.org/wa/internal/app/appcir"
	"wa-lang.org/wa/internal/app/appdev"
	"wa-lang.org/wa/internal/app/appdoc"
	"wa-lang.org/wa/internal/app/appfmt"
	"wa-lang.org/wa/internal/app/appgo2wa"
	"wa-lang.org/wa/internal/app/appinit"
	"wa-lang.org/wa/internal/app/applex"
	"wa-lang.org/wa/internal/app/applogo"
	"wa-lang.org/wa/internal/app/applsp"
	"wa-lang.org/wa/internal/app/appplay"
	"wa-lang.org/wa/internal/app/apprun"
	"wa-lang.org/wa/internal/app/appssa"
	"wa-lang.org/wa/internal/app/apptest"
	"wa-lang.org/wa/internal/app/appyacc"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/waroot"
)

func InitWaroot() {
	waroot.InitWarootDir()
}

func Main() {
	cliApp := cli.NewApp()
	cliApp.Name = "Wa"
	cliApp.Usage = "Wa is a tool for managing Wa source code."
	cliApp.Copyright = "Copyright 2018 The Wa Authors. All rights reserved."
	cliApp.Version = waroot.GetVersion()
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
		// 用于调试(隐藏)
		appdev.CmdDev,

		// 用户可见子命令
		appplay.CmdPlay,
		appinit.CmdInit,
		appbuild.CmdBuild,
		apprun.CmdRun,
		appfmt.CmdFmt,
		apptest.CmdTest,
		appyacc.CmdYacc,
		applogo.CmdLogo,

		// 开发者用于测试(隐藏)
		applex.CmdLex,
		appast.CmdAst,
		appssa.CmdSsa,

		// 待完善的子命令(隐藏)
		appgo2wa.CmdGo2wa,
		appcir.CmdCir,
		appdoc.CmdDoc,
		applsp.CmdLsp,
	}

	cliApp.Run(os.Args)
}

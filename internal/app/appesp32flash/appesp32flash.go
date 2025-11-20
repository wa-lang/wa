// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appesp32flash

import (
	"fmt"
	"os"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/native/espflash"
)

var CmdESP32Flash = &cli.Command{
	Hidden: true,
	Name:   "esp32flash",
	Usage:  "upload esp32 image file to board",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "set uart port",
		},
		&cli.IntFlag{
			Name:    "baudrate",
			Aliases: []string{"b"},
			Usage:   "set uart baudrate",
			Value:   115200,
		},
		&cli.IntFlag{
			Name:    "addr",
			Aliases: []string{"a"},
			Usage:   "set flash address",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		portName := c.String("port")
		baudRate := c.Int("baudrate")
		flashAddr := c.Int("addr")
		filePath := c.Args().First()

		client, err := espflash.Open(portName, baudRate)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer client.Close()

		err = client.FlashFile(uint32(flashAddr), filePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Done")
		return nil
	},
}

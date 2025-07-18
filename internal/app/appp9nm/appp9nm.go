// Copyright (C) 2024 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package appp9nm

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"wa-lang.org/wa/internal/3rdparty/cli"
	"wa-lang.org/wa/internal/p9asm/objfile"
)

var CmdP9Nm = &cli.Command{
	Hidden: true,
	Name:   "p9nm",
	Usage:  "lists the symbols defined or used by an object file, archive, or executable",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "sort",
			Aliases: []string{"n"},
			Usage:   "sort output in the given order (address|name|none|size)",
			Value:   "name",
		},
		&cli.BoolFlag{
			Name:  "size",
			Usage: "print symbol size in decimal between address and type",
		},
		&cli.BoolFlag{
			Name:  "type",
			Usage: "print symbol type after name",
		},
	},
	Action: func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "no input file")
			os.Exit(1)
		}

		sortOrder := c.String("sort")
		printSize := c.Bool("size")
		printType := c.Bool("type")
		filePrefix := c.NArg() > 1

		switch sortOrder {
		case "address", "name", "none", "size":
			// ok
		default:
			fmt.Fprintf(os.Stderr, "nm: unknown sort order %q\n", sortOrder)
			os.Exit(2)
		}

		var lastErr error
		for _, file := range c.Args().Slice() {
			f, err := objfile.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				lastErr = err
				continue
			}
			defer f.Close()

			syms, err := f.Symbols()
			if err != nil {
				fmt.Fprintf(os.Stderr, "reading %s: %v\n", file, err)
				lastErr = err
				continue
			}
			if len(syms) == 0 {
				fmt.Printf("reading %s: no symbols", file)
				continue
			}

			switch sortOrder {
			case "address":
				sort.Slice(syms, func(i, j int) bool {
					return syms[i].Addr < syms[j].Addr
				})
			case "name":
				sort.Slice(syms, func(i, j int) bool {
					return syms[i].Name < syms[j].Name
				})
			case "size":
				sort.Slice(syms, func(i, j int) bool {
					return syms[i].Size < syms[j].Size
				})
			}

			w := bufio.NewWriter(os.Stdout)
			for _, sym := range syms {
				if filePrefix {
					fmt.Fprintf(w, "%s:\t", file)
				}
				if sym.Code == 'U' {
					fmt.Fprintf(w, "%8s", "")
				} else {
					fmt.Fprintf(w, "%8x", sym.Addr)
				}
				if printSize {
					fmt.Fprintf(w, " %10d", sym.Size)
				}
				fmt.Fprintf(w, " %c %s", sym.Code, sym.Name)
				if printType && sym.Type != "" {
					fmt.Fprintf(w, " %s", sym.Type)
				}
				fmt.Fprintf(w, "\n")
			}
			w.Flush()
		}

		if lastErr != nil {
			os.Exit(1)
		}
		return nil
	},
}

// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"fmt"
	"log"
	"strings"
)

// 将单条指令规则合并到解码决策树中。
//
// 统一判定顺序: 基础编码 -> 64位模式 -> REX/前缀 -> 地址/数据位宽 -> 寄存器/内存判定.
// 最终树会合并 "any" 分支以消除歧义, 确保每条指令有唯一的识别路径
func (root *Prog) add(text, opcode, valid32, valid64, cpuid, tags string) {
	// These are not real instructions: they are either
	// prefixes for other instructions, composite instructions
	// built from multiple individual instructions, or alternate
	// mnemonics of other encodings.
	// Discard for disassembly, because we want a unique decoding.
	if strings.Contains(tags, "pseudo") {
		return
	}

	// Treat REX.W + opcode as being like having an "operand64" tag.
	// The REX.W flag sets the operand size to 64 bits; in this way it is
	// not much different than the 66 prefix that inverts 32 vs 16 bits.
	if strings.Contains(opcode, "REX.W") {
		if !strings.Contains(tags, "operand64") {
			if tags != "" {
				tags += ","
			}
			tags += "operand64"
		}
	}

	// If there is more than one operand size given, we need to do
	// a separate add for each size, because we need multiple
	// keys to be added in the operand size branch, and the code makes
	// a linear pass through the tree adding just one key to each node.
	// We would need to do the same for any other possible repeated tag
	// (for example, if an instruction could have multiple address sizes)
	// but so far operand size is the only tag we have needed to repeat.
	if strings.Count(tags, "operand") > 1 {
		f := strings.Split(tags, ",")
		var ops []string
		w := 0
		for _, tag := range f {
			if strings.HasPrefix(tag, "operand") {
				ops = append(ops, tag)
			} else {
				if strings.Contains(tag, "operand") {
					panic(fmt.Sprintf("unknown tag %q", tag))
				}
				f[w] = tag
				w++
			}
		}
		f = f[:w]
		for _, op := range ops {
			root.add(text, opcode, valid32, valid64, cpuid, strings.Join(append(f, op), ","))
		}
		return
	}

	p := root
	walk := func(action, item string) {
		p = p.walk(action, item, text, opcode)
	}

	// Ignore VEX instructions for now.
	if strings.HasPrefix(opcode, "VEX") {
		if !strings.HasPrefix(text, "VMOVNTDQ") &&
			!strings.HasPrefix(text, "VMOVDQA") &&
			!strings.HasPrefix(text, "VMOVDQU") &&
			!strings.HasPrefix(text, "VZEROUPPER") {
			return
		}
		if !strings.HasPrefix(opcode, "VEX.256") && !strings.HasPrefix(text, "VZEROUPPER") {
			return
		}
		if !strings.Contains(tags, "VEXC4") {
			root.add(text, opcode, valid32, valid64, cpuid, tags+",VEXC4")
		}
		encoding := strings.Fields(opcode)
		walk("decode", encoding[1])
		walk("is64", "any")
		if strings.Contains(tags, "VEXC4") {
			walk("prefix", "C4")
		} else {
			walk("prefix", "C5")
		}
		for _, pref := range strings.Split(encoding[0], ".") {
			if isVexEncodablePrefix[pref] {
				walk("prefix", pref)
			}
		}
	}

	var rex, prefix string
	encoding := strings.Fields(opcode)
	if len(encoding) > 0 && strings.HasPrefix(encoding[0], "REX") {
		rex = encoding[0]
		encoding = encoding[1:]
		if len(encoding) > 0 && encoding[0] == "+" {
			encoding = encoding[1:]
		}
	}
	if len(encoding) > 0 && isPrefix[encoding[0]] {
		prefix = encoding[0]
		encoding = encoding[1:]
	}
	if rex == "" && len(encoding) > 0 && strings.HasPrefix(encoding[0], "REX") {
		rex = encoding[0]
		if rex == "REX" {
			log.Printf("REX without REX.W: %s %s", text, opcode)
		}
		encoding = encoding[1:]
		if len(encoding) > 0 && encoding[0] == "+" {
			encoding = encoding[1:]
		}
	}
	if len(encoding) > 0 && isPrefix[encoding[0]] {
		log.Printf("%s %s: too many prefixes", text, opcode)
		return
	}

	var haveModRM, havePlus bool
	var usedReg string
	for len(encoding) > 0 && (isHex(encoding[0]) || isSlashNum(encoding[0])) {
		key := encoding[0]
		if isSlashNum(key) {
			if usedReg != "" {
				log.Printf("%s %s: multiple modrm checks", text, opcode)
			}
			haveModRM = true
			usedReg = key
		}
		if i := strings.Index(key, "+"); i >= 0 {
			key = key[:i+1]
			havePlus = true
		}
		walk("decode", key)
		encoding = encoding[1:]
	}

	if valid32 != "V" {
		walk("is64", "1")
	} else if valid64 != "V" {
		walk("is64", "0")
	} else {
		walk("is64", "any")
	}

	if prefix == "" {
		prefix = "0"
	}
	walk("prefix", prefix)

	if strings.Contains(tags, "address16") {
		walk("addrsize", "16")
	} else if strings.Contains(tags, "address32") {
		walk("addrsize", "32")
	} else if strings.Contains(tags, "address64") {
		walk("addrsize", "64")
	} else {
		walk("addrsize", "any")
	}

	if strings.Contains(tags, "operand16") {
		walk("datasize", "16")
	} else if strings.Contains(tags, "operand32") {
		walk("datasize", "32")
	} else if strings.Contains(tags, "operand64") {
		walk("datasize", "64")
	} else {
		walk("datasize", "any")
	}

	if len(encoding) > 0 && encoding[0] == "/r" {
		haveModRM = true
	}
	if haveModRM {
		if strings.Contains(tags, "modrm_regonly") {
			walk("ismem", "0")
		} else if strings.Contains(tags, "modrm_memonly") {
			walk("ismem", "1")
		} else {
			walk("ismem", "any")
		}
	}

	walk("op", strings.Fields(text)[0])

	if len(encoding) > 0 && strings.HasPrefix(encoding[0], "VEX") {
		for _, field := range encoding[2:] {
			walk("read", field)
		}
	} else {
		for _, field := range encoding {
			walk("read", field)
		}
	}

	var usedRM string
	for _, arg := range strings.Fields(text)[1:] {
		arg = strings.TrimRight(arg, ",")
		if usesReg[arg] && !haveModRM && !havePlus {
			log.Printf("%s %s: no modrm field to use for %s", text, opcode, arg)
			continue
		}
		if usesRM[arg] && !haveModRM {
			log.Printf("%s %s: no modrm field to use for %s", text, opcode, arg)
			continue
		}
		if usesReg[arg] {
			if usedReg != "" {
				log.Printf("%s %s: modrm reg field used by both %s and %s", text, opcode, usedReg, arg)
				continue
			}
			usedReg = arg
		}
		if usesRM[arg] {
			if usedRM != "" {
				log.Printf("%s %s: modrm r/m field used by both %s and %s", text, opcode, usedRM, arg)
				continue
			}
			usedRM = arg
		}
		walk("arg", arg)
	}

	walk("match", "!")
}

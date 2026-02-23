// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package wat2la

import (
	"fmt"
	"io"
)

const (
	kFuncStart = "_start"
	kFuncMain  = "main"
	kFuncMain2 = "_main"

	kAppArgcName = ".Wa.App.argc"
	kAppArgvName = ".Wa.App.argv"
	kAppEnvpName = ".Wa.App.envp"
)

// 启动函数
func (p *wat2laWorker) buildEntryFunc(w io.Writer) error {
	// 生成命令行参数
	p.gasGlobal(w, kAppArgcName)
	p.gasDefI64(w, kAppArgcName, 0)
	p.gasGlobal(w, kAppArgvName)
	p.gasDefI64(w, kAppArgvName, 0)
	p.gasGlobal(w, kAppEnvpName)
	p.gasDefI64(w, kAppEnvpName, 0)

	switch p.entryFuncName {
	case kFuncStart:
		return p.buildEntryFunc_start(w)
	case kFuncMain, kFuncMain2:
		return p.buildEntryFunc_x(w, kFuncMain)
	default:
		if p.entryFuncName != "" {
			return p.buildEntryFunc_x(w, p.entryFuncName)
		} else {
			return p.buildEntryFunc_start(w)
		}
	}
}

func (p *wat2laWorker) buildEntryFunc_start(w io.Writer) error {
	p.gasComment(w, "汇编程序入口函数")
	p.gasSectionTextStart(w)
	p.gasGlobal(w, kFuncStart)
	fmt.Fprintf(w, "%s:\n", kFuncStart)

	fmt.Fprintln(w, "    addi.d $sp, $sp, -16")
	fmt.Fprintln(w, "    st.d   $ra, $sp, 8")
	fmt.Fprintln(w, "    st.d   $fp, $sp, 0")
	fmt.Fprintln(w, "    addi.d $fp, $sp, 0")
	fmt.Fprintln(w, "    addi.d $sp, $sp, -32")
	fmt.Fprintln(w)

	// TODO: 补充命令行参数

	if p.m.Memory != nil {
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryInitFuncName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryInitFuncName)
		fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
	}
	if p.m.Table != nil {
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kTableInitFuncName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kTableInitFuncName)
		fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
	}

	if p.m.Start != "" {
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kFuncNamePrefix+p.m.Start)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kFuncNamePrefix+p.m.Start)
		fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
	}
	if p.m.Start != kFuncMain {
		for _, fn := range p.m.Funcs {
			switch fn.Name {
			case kFuncMain, kFuncMain2:
				fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kFuncNamePrefix+fixName(fn.Name))
				fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kFuncNamePrefix+fixName(fn.Name))
				fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
			}
		}
	}
	fmt.Fprintln(w)

	p.gasCommentInFunc(w, "runtime.exit(0)")
	fmt.Fprintf(w, "    addi.d    $a0, $zero, 0 # a0 = 0\n")
	fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kRuntimeExit)
	fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kRuntimeExit)
	fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
	fmt.Fprintln(w)

	p.gasCommentInFunc(w, "exit 后这里不会被执行, 但是依然保留")
	fmt.Fprintln(w, "    addi.d $sp, $fp, 0")
	fmt.Fprintln(w, "    ld.d   $ra, $sp, 8")
	fmt.Fprintln(w, "    ld.d   $fp, $sp, 0")
	fmt.Fprintln(w, "    addi.d $sp, $sp, 16")
	fmt.Fprintln(w, "    jirl   $zero, $ra, 0")
	fmt.Fprintln(w)

	return nil
}

func (p *wat2laWorker) buildEntryFunc_x(w io.Writer, entryFuncName string) error {
	p.gasComment(w, "汇编程序入口函数")
	p.gasSectionTextStart(w)
	p.gasGlobal(w, entryFuncName)
	fmt.Fprintf(w, "%s:\n", entryFuncName)

	fmt.Fprintln(w, "    addi.d $sp, $sp, -16")
	fmt.Fprintln(w, "    st.d   $ra, $sp, 8")
	fmt.Fprintln(w, "    st.d   $fp, $sp, 0")
	fmt.Fprintln(w, "    addi.d $fp, $sp, 0")
	fmt.Fprintln(w, "    addi.d $sp, $sp, -32")
	fmt.Fprintln(w)

	// 如果是 main 函数, 说明不是自定义的 _start, 需要补充获取命令行信息
	if entryFuncName == kFuncMain {
		// TODO: 补充命令行参数
	}

	if p.m.Memory != nil {
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kMemoryInitFuncName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kMemoryInitFuncName)
		fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
	}
	if p.m.Table != nil {
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kTableInitFuncName)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kTableInitFuncName)
		fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
	}

	if p.m.Start != "" {
		fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kFuncNamePrefix+p.m.Start)
		fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kFuncNamePrefix+p.m.Start)
		fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
	}
	for _, fn := range p.m.Funcs {
		switch fn.Name {
		case kFuncMain, kFuncMain2:
			fmt.Fprintf(w, "    pcalau12i $t0, %%pc_hi20(%s)\n", kFuncNamePrefix+fixName(fn.Name))
			fmt.Fprintf(w, "    addi.d    $t0, $t0, %%pc_lo12(%s)\n", kFuncNamePrefix+fixName(fn.Name))
			fmt.Fprintf(w, "    jirl      $ra, $t0, 0\n")
		}
	}
	fmt.Fprintln(w)

	p.gasCommentInFunc(w, "return 0")
	fmt.Fprintf(w, "    addi.d    $a0, $zero, 0 # a0 = 0\n")
	fmt.Fprintln(w)

	p.gasCommentInFunc(w, "return")
	fmt.Fprintln(w, "    addi.d $sp, $fp, 0")
	fmt.Fprintln(w, "    ld.d   $ra, $sp, 8")
	fmt.Fprintln(w, "    ld.d   $fp, $sp, 0")
	fmt.Fprintln(w, "    addi.d $sp, $sp, 16")
	fmt.Fprintln(w, "    jirl   $zero, $ra, 0")
	fmt.Fprintln(w)

	return nil
}

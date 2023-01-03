// 版权 @2019 凹语言 作者。保留所有权利。

package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/template"
	"time"

	"wa-lang.org/wa/internal/ast"
	"wa-lang.org/wa/internal/backends/compiler_c"
	"wa-lang.org/wa/internal/backends/compiler_llvm"
	"wa-lang.org/wa/internal/backends/compiler_wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/format"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/parser"
	"wa-lang.org/wa/internal/scanner"
	"wa-lang.org/wa/internal/ssa"
	"wa-lang.org/wa/internal/token"
	"wa-lang.org/wa/internal/waroot"
)

// 命令行程序对象
type App struct {
	opt  Option
	path string
	src  string
}

// 构建命令行程序对象
func NewApp(opt *Option) *App {
	logger.Tracef(&config.EnableTrace_app, "opt: %+v", opt)

	p := &App{}
	if opt != nil {
		p.opt = *opt
	}
	if p.opt.Clang == "" {
		if runtime.GOOS == "windows" {
			p.opt.Clang, _ = exec.LookPath("clang.exe")
		} else {
			p.opt.Clang, _ = exec.LookPath("clang")
		}
		if p.opt.Clang == "" {
			p.opt.Clang = "clang"
		}
	}
	if p.opt.Llc == "" {
		if runtime.GOOS == "windows" {
			p.opt.Llc, _ = exec.LookPath("llc.exe")
		} else {
			p.opt.Llc, _ = exec.LookPath("llc")
		}
		if p.opt.Llc == "" {
			p.opt.Llc = "llc"
		}
	}
	if p.opt.TargetOS == "" {
		p.opt.TargetOS = config.WaOS_Walang
	}
	if p.opt.TargetArch == "" {
		p.opt.TargetArch = config.WaArch_Wasm
	}

	return p
}

func (p *App) GetConfig() *config.Config {
	return p.opt.Config()
}

func (p *App) InitApp(name, pkgpath string, update bool) error {
	if name == "" {
		return fmt.Errorf("init failed: <%s> is empty", name)
	}

	if !update {
		if _, err := os.Lstat(name); err == nil {
			return fmt.Errorf("init failed: <%s> exists", name)
		}
	}

	var info = struct {
		Name    string
		Pkgpath string
		Year    int
	}{
		Name:    name,
		Pkgpath: pkgpath,
		Year:    time.Now().Year(),
	}

	appFS := waroot.GetExampleAppFS()
	err := fs.WalkDir(appFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		data, err := fs.ReadFile(appFS, path)
		if err != nil {
			return err
		}

		tmpl, err := template.New(path).Parse(string(data))
		if err != nil {
			return err
		}

		dstpath := filepath.Join(name, path)
		os.MkdirAll(filepath.Dir(dstpath), 0777)

		f, err := os.Create(dstpath)
		if err != nil {
			return err
		}
		defer f.Close()

		err = tmpl.Execute(f, &info)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	vendorFS := waroot.GetExampleVendorFS()
	err = fs.WalkDir(vendorFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		data, err := fs.ReadFile(vendorFS, path)
		if err != nil {
			return err
		}

		tmpl, err := template.New(path).Parse(string(data))
		if err != nil {
			return err
		}

		dstpath := filepath.Join(name, "vendor", path)
		os.MkdirAll(filepath.Dir(dstpath), 0777)

		f, err := os.Create(dstpath)
		if err != nil {
			return err
		}
		defer f.Close()

		err = tmpl.Execute(f, &info)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *App) Fmt(path string) error {
	if path == "" {
		path, _ = os.Getwd()
	}

	if strings.HasSuffix(path, "...") {
		panic("TODO: fmt dir/...")
	}

	fi, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		panic("TODO: fmt dir")
	}

	code, err := format.File(nil, path, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(path, code, 0666)
}

func (p *App) Lex(filename string) error {
	src, err := p.readSource(filename, nil)
	if err != nil {
		return err
	}

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(filename, fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}

	return nil
}

func (p *App) AST(filename string) error {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(nil, fset, filename, nil, 0)
	if err != nil {
		return err
	}

	return ast.Print(fset, f)
}

func (p *App) SSA(filename string) error {
	cfg := p.opt.Config()
	prog, err := loader.LoadProgram(cfg, filename)
	if err != nil {
		return err
	}

	prog.SSAMainPkg.WriteTo(os.Stdout)

	var funcNames []string
	for name, x := range prog.SSAMainPkg.Members {
		if _, ok := x.(*ssa.Function); ok {
			funcNames = append(funcNames, name)
		}
	}
	sort.Strings(funcNames)
	for _, s := range funcNames {
		prog.SSAMainPkg.Func(s).WriteTo(os.Stdout)
	}

	return nil
}

func (p *App) CIR(filename string) error {
	cfg := p.opt.Config()
	prog, err := loader.LoadProgram(cfg, filename)
	if err != nil {
		return err
	}

	var c compiler_c.CompilerC
	c.CompilePackage(prog.SSAMainPkg)
	print("\n\n")
	print(c.String())

	return nil
}

func (p *App) LLVM(infile string, outfile string, target string, debug bool) error {
	cfg := p.opt.Config()

	instat, err := os.Stat(infile)
	if err != nil {
		return err
	}

	// Calculate the outfile path if not given.
	if len(outfile) == 0 {
		if instat.IsDir() {
			dir := path.Base(infile)
			outfile = infile + dir + ".exe"
		} else {
			ext := path.Ext(infile)
			if len(ext) == 0 {
				outfile = infile + ".exe"
			} else {
				pos := strings.Index(infile, ext)
				outfile = infile[0:pos] + ".exe"
			}
		}
	}

	// Calculate the outfile LLVM-IR file path and the output assembly file path.
	llfile, asmfile := "", ""
	ext := path.Ext(outfile)
	if len(ext) == 0 {
		llfile = outfile + ".ll"
		asmfile = outfile + ".s"
	} else {
		pos := strings.Index(outfile, ext)
		llfile = outfile[0:pos] + ".ll"
		asmfile = outfile[0:pos] + ".s"
	}

	// Do the real compile work.
	prog, err := loader.LoadProgram(cfg, infile)
	if err != nil {
		return err
	}
	output, err := compiler_llvm.New(target, debug).Compile(prog)
	if err != nil {
		return err
	}

	// Write the outfile LLVM-IR to an intermediate .ll file.
	if err := os.WriteFile(llfile, []byte(output), 0644); err != nil {
		return err
	}

	// Invoke command `llc xxx.ll -mtriple=xxx`.
	llc := []string{llfile}
	if target != "" {
		llc = append(llc, "-mtriple", target)
	}
	// Add target specific options.
	switch target {
	case "avr":
		llc = append(llc, "-mcpu=atmega328")
	default:
	}
	cmd0 := exec.Command(p.opt.Llc, llc...)
	cmd0.Stderr = os.Stderr
	if err := cmd0.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "**** failed to invoke LLVM ****\n")
		return err
	}

	// TODO: This is a temporary solution for AVR-Arduino. We generate
	// an Arduino project instead of an ELF.
	if target == "avr" {
		// Create a new directory for the output Arduino project.
		ext, outdir := path.Ext(outfile), ""
		if len(ext) > 0 {
			pos := strings.Index(outfile, ext)
			outdir = outfile[0:pos]
		}
		if err := os.RemoveAll(outdir); err != nil {
			return err
		}
		if err := os.Mkdir(outdir, 0755); err != nil {
			return err
		}
		// Move the assembly file to the project directory.
		newAsmFile := strings.ReplaceAll(asmfile, ".s", ".S")
		if err := os.Rename(asmfile, path.Join(outdir, newAsmFile)); err != nil {
			return err
		}
		// Create the project main '.ino' file.
		inoFile := path.Join(outdir, path.Base(outdir)+".ino")
		inoStr := "void setup(void) {}\nextern \"C\" { extern void wa_main(void); }\nvoid loop(void) { wa_main(); }\n"
		if err := os.WriteFile(inoFile, []byte(inoStr), 0644); err != nil {
			return err
		}
		return nil
	}

	// Invoke command `clang xxx.s -o outfile --target=xxx`.
	clangArgs := []string{asmfile, "-static", "-o", outfile}
	if target != "" {
		clangArgs = append(clangArgs, "-target", target)
	}
	if p.opt.Debug {
		clangArgs = append(clangArgs, "-v")
	}
	cmd1 := exec.Command(p.opt.Clang, clangArgs...)
	cmd1.Stderr = os.Stderr
	cmd1.Stdout = os.Stdout
	if err := cmd1.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "**** failed to invoke CLANG ****\n")
		return err
	}

	return nil
}

func (p *App) WASM(filename string) ([]byte, error) {
	cfg := p.opt.Config()
	prog, err := loader.LoadProgram(cfg, filename)
	if err != nil {
		return nil, err
	}

	// 凹中文的源码启动函数为【启】，对应的wat函数名应当是"$0xe590af"
	main := "main"
	if strings.HasSuffix(filename, ".wz") {
		main = "$0xe590af"
	}

	output, err := compiler_wat.New().Compile(prog, main)

	if err != nil {
		return nil, err
	}

	return []byte(output), nil
}

func (p *App) readSource(filename string, src interface{}) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		case *bytes.Buffer:
			if s != nil {
				return s.Bytes(), nil
			}
		case io.Reader:
			d, err := io.ReadAll(s)
			return d, err
		}
		return nil, errors.New("invalid source")
	}

	d, err := os.ReadFile(filename)
	return d, err
}

func (p *App) isWaFile(path string) bool {
	if fi, err := os.Lstat(path); err == nil && fi.Mode().IsRegular() {
		return strings.HasSuffix(strings.ToLower(path), ".wa")
	}
	return false
}

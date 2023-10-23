// 版权 @2021 凹语言 作者。保留所有权利。

package compiler_wat

import (
	"bytes"
	_ "embed"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"wa-lang.org/wa/internal/backends/compiler_wat/wir"
	"wa-lang.org/wa/internal/backends/compiler_wat/wir/wat"
	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/loader"
	"wa-lang.org/wa/internal/logger"
	"wa-lang.org/wa/internal/ssa"
	"wa-lang.org/wa/waroot"
)

type Compiler struct {
	prog *loader.Program

	module *wir.Module
	tLib   *typeLib
}

func New() *Compiler {
	p := new(Compiler)
	p.module = wir.NewModule()
	p.module.AddGlobal("$wa.runtime.closure_data", "", p.module.GenValueType_Ref(p.module.VOID), false, nil)
	wir.SetCurrentModule(p.module)
	return p
}

func (p *Compiler) Compile(prog *loader.Program, mainFunc string) (output string, err error) {
	p.prog = prog
	p.CompileWsFiles(prog)

	p.tLib = newTypeLib(p.module, prog)

	var pkgnames []string
	for n := range prog.Pkgs {
		pkgnames = append(pkgnames, n)
	}
	sort.Strings(pkgnames)

	for i, v := range pkgnames {
		if v == "runtime" && i != 0 {
			pkgnames[i] = pkgnames[0]
			pkgnames[0] = "runtime"
		}
	}

	p.CompilePkgType(prog.Pkgs["runtime"].SSAPkg)

	for _, n := range pkgnames {
		p.CompilePkgGlobal(prog.Pkgs[n].SSAPkg)
	}
	for _, n := range pkgnames {
		p.CompilePkgFunc(prog.Pkgs[n].SSAPkg)
	}

	p.tLib.finish()

	{
		var f wir.Function
		f.InternalName, f.ExternalName = "_start", "_start"
		n, _ := wir.GetPkgMangleName(prog.SSAMainPkg.Pkg.Path())
		n += ".init"
		f.Insts = append(f.Insts, wat.NewInstCall(n))

		if mainFunc != "" {
			n, _ = wir.GetPkgMangleName(prog.SSAMainPkg.Pkg.Path())
			n += "."
			n += mainFunc
			f.Insts = append(f.Insts, wat.NewInstCall(n))
		}

		p.module.AddFunc(&f)
	}

	// p.GenJsBind()
	// p.GenJSBinding()

	return p.module.ToWatModule().String(), nil
}

func (p *Compiler) CompileWsFiles(prog *loader.Program) {
	var sb strings.Builder

	sb.WriteString(waroot.GetBaseWsCode(config.WaBackend_wat))
	sb.WriteString("\n")

	var pkgpathList = make([]string, 0, len(prog.Pkgs))
	for pkgpath := range prog.Pkgs {
		pkgpathList = append(pkgpathList, pkgpath)
	}
	sort.Strings(pkgpathList)

	var lineCommentSep = ";; -" + strings.Repeat("-", 60-4) + "\n"

	for _, pkgpath := range pkgpathList {
		pkg := prog.Pkgs[pkgpath]
		if len(pkg.WsFiles) == 0 {
			continue
		}

		func() {
			sb.WriteString(lineCommentSep)
			sb.WriteString(";; package: " + pkgpath + "\n")
			sb.WriteString(lineCommentSep)
			sb.WriteString("\n")

			for _, sf := range pkg.WsFiles {
				sb.WriteString(";; file: " + sf.Name + "\n")
				sb.WriteString("\n")

				sb.WriteString(strings.TrimSpace(sf.Code))
				sb.WriteString("\n")
			}
		}()
	}

	p.module.BaseWat = sb.String()
}

func (p *Compiler) CompilePkgType(ssaPkg *ssa.Package) {
	var memnames []string
	for name := range ssaPkg.Members {
		memnames = append(memnames, name)
	}
	sort.Strings(memnames)

	for _, name := range memnames {
		m := ssaPkg.Members[name]
		if t, ok := m.(*ssa.Type); ok {
			p.compileType(t)
		}
	}
}

func (p *Compiler) CompilePkgGlobal(ssaPkg *ssa.Package) {
	var memnames []string
	for name := range ssaPkg.Members {
		memnames = append(memnames, name)
	}
	sort.Strings(memnames)

	for _, name := range memnames {
		m := ssaPkg.Members[name]
		if global, ok := m.(*ssa.Global); ok {
			p.compileGlobal(global)
		}
	}
}

func (p *Compiler) CompilePkgFunc(ssaPkg *ssa.Package) {
	var memnames []string
	for name := range ssaPkg.Members {
		memnames = append(memnames, name)
	}
	sort.Strings(memnames)

	for _, name := range memnames {
		m := ssaPkg.Members[name]
		if fn, ok := m.(*ssa.Function); ok {
			CompileFunc(fn, p.prog, p.tLib, p.module)
		}
	}
}

func CompileFunc(f *ssa.Function, prog *loader.Program, tLib *typeLib, module *wir.Module) {
	if len(f.Blocks) < 1 {
		if f.RuntimeGetter() {
			module.AddFunc(newFunctionGenerator(prog, module, tLib).genGetter(f))
		} else if f.RuntimeSetter() {
			module.AddFunc(newFunctionGenerator(prog, module, tLib).genSetter(f))
		} else if f.RuntimeSizer() {
			module.AddFunc(newFunctionGenerator(prog, module, tLib).genSizer(f))
		} else if iname0, iname1 := f.ImportName(); len(iname0) > 0 && len(iname1) > 0 {
			var fn_name string
			if len(f.LinkName()) > 0 {
				fn_name = f.LinkName()
			} else {
				fn_name, _ = wir.GetFnMangleName(f, prog.Manifest.MainPkg)
			}

			sig := tLib.GenFnSig(f.Signature)
			module.AddImportFunc(iname0, iname1, fn_name, sig)
		}
		return
	}
	module.AddFunc(newFunctionGenerator(prog, module, tLib).genFunction(f))
}

func (p *Compiler) GenJsBind() {
	for _, g := range p.module.Globals {
		if len(g.Name_exp) == 0 {
			continue
		}

		ref_type, ok := g.Type.(*wir.Ref)
		if !ok {
			logger.Fatalf("Exported global: %s should be *T.", g.Name)
		}

		switch typ := ref_type.Base.(type) {
		case *wir.U8:
			println("Name:", g.Name_exp, ", Type:", typ.Named())

		case *wir.U16:
			println("Name:", g.Name_exp, ", Type:", typ.Named())

		case *wir.I32:
			println("Name:", g.Name_exp, ", Type:", typ.Named())

		case *wir.U32:
			println("Name:", g.Name_exp, ", Type:", typ.Named())

		case *wir.I64:
			println("Name:", g.Name_exp, ", Type:", typ.Named())

		case *wir.U64:
			println("Name:", g.Name_exp, ", Type:", typ.Named())

		case *wir.Bool:
			println("Name:", g.Name_exp, ", Type:", typ.Named())

		case *wir.Rune:
			println("Name:", g.Name_exp, ", Type:", typ.Named())

		case *wir.String:
			println("Name:", g.Name_exp, ", Type:", typ.Named())

		default:
			//非基本类型，不导出
		}
	}

	for _, f := range p.module.Funcs {
		if !f.ExplicitExported {
			continue
		}

		println("Function:", f.ExternalName)
		for i, p := range f.Params {
			switch typ := p.Type().(type) {
			case *wir.U8:
				fmt.Printf("\tParam[%d] - Name: %s, Type: %s\n", i, p.Name(), typ.Named())

			case *wir.U16:
				fmt.Printf("\tParam[%d] - Name: %s, Type: %s\n", i, p.Name(), typ.Named())

			case *wir.I32:
				fmt.Printf("\tParam[%d] - Name: %s, Type: %s\n", i, p.Name(), typ.Named())

			case *wir.U32:
				fmt.Printf("\tParam[%d] - Name: %s, Type: %s\n", i, p.Name(), typ.Named())

			case *wir.I64:
				fmt.Printf("\tParam[%d] - Name: %s, Type: %s\n", i, p.Name(), typ.Named())

			case *wir.U64:
				fmt.Printf("\tParam[%d] - Name: %s, Type: %s\n", i, p.Name(), typ.Named())

			case *wir.Bool:
				fmt.Printf("\tParam[%d] - Name: %s, Type: %s\n", i, p.Name(), typ.Named())

			case *wir.Rune:
				fmt.Printf("\tParam[%d] - Name: %s, Type: %s\n", i, p.Name(), typ.Named())

			case *wir.String:
				fmt.Printf("\tParam[%d] - Name: %s, Type: %s\n", i, p.Name(), typ.Named())

			default:
				//非基本类型，不导出
			}
		}

		for i, r := range f.Results {
			switch typ := r.(type) {
			case *wir.U8:
				fmt.Printf("\tRet[%d] - Type: %s\n", i, typ.Named())

			case *wir.U16:
				fmt.Printf("\tRet[%d] - Type: %s\n", i, typ.Named())

			case *wir.I32:
				fmt.Printf("\tRet[%d] - Type: %s\n", i, typ.Named())

			case *wir.U32:
				fmt.Printf("\tRet[%d] - Type: %s\n", i, typ.Named())

			case *wir.I64:
				fmt.Printf("\tRet[%d] - Type: %s\n", i, typ.Named())

			case *wir.U64:
				fmt.Printf("\tRet[%d] - Type: %s\n", i, typ.Named())

			case *wir.Bool:
				fmt.Printf("\tRet[%d] - Type: %s\n", i, typ.Named())

			case *wir.Rune:
				fmt.Printf("\tRet[%d] - Type: %s\n", i, typ.Named())

			case *wir.String:
				fmt.Printf("\tRet[%d] - Type: %s\n", i, typ.Named())

			default:
				//非基本类型，不导出
			}

		}

	}
}

//go:embed js_binding_tmpl.js
var js_binding_tmpl string

type JSGlobal struct {
	Name string
	Type string
}

type JSFunc struct {
	Name       string
	Params     string
	PreCall    string
	GetResults string
	Release    string
	Return     string
}

type JSModule struct {
	Filename string
	Pkg      string
	Globals  []JSGlobal
	Funcs    []JSFunc
}

func stripNamePrefix(name string) string {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return name[i+1:]
	}
	return name
}

func (p *Compiler) globalsForJsBinding() []JSGlobal {
	var globals []JSGlobal
	for _, g := range p.module.Globals {
		if len(g.Name_exp) == 0 {
			continue
		}

		ref_type, ok := g.Type.(*wir.Ref)
		if !ok {
			logger.Fatalf("Exported global: %s should be *T.", g.Name)
		}
		switch typ := ref_type.Base.(type) {
		case *wir.U8, *wir.U16, *wir.I32, *wir.U32, *wir.I64, *wir.U64, *wir.Bool, *wir.Rune, *wir.String:
			name := stripNamePrefix(g.Name_exp)
			tp := typ.Named()
			globals = append(globals, JSGlobal{Name: name, Type: tp})
		default: //非基本类型，不导出
		}

	}
	return globals
}

func (p *Compiler) funcsForJSBinding() []JSFunc {
	var funcs []JSFunc

	// 函数
	for _, f := range p.module.Funcs {
		if !f.ExplicitExported {
			continue
		}

		fn := JSFunc{Name: stripNamePrefix(f.ExternalName)}

		// 参数名称列表
		var sb strings.Builder
		for i, p := range f.Params {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(p.Name())
		}
		fn.Params = sb.String()

		// 参数
		if len(f.Params) > 0 {
			var sb strings.Builder
			var sbr strings.Builder
			for i, p := range f.Params {
				name := p.Name()
				switch p.Type().(type) {
				case *wir.U8, *wir.U16, *wir.I32, *wir.U32, *wir.I64, *wir.U64, *wir.Bool, *wir.Rune:
					sb.WriteString(fmt.Sprintf("params.push(%s);\n", name))
				case *wir.String: // 字符串类型需要转换为[l,b,d]的形式
					sb.WriteString(fmt.Sprintf("let p%d = this._mem_util.set_string(%s);\n", i, name))
					sb.WriteString(fmt.Sprintf("params = params.concat(p%d);\n", i))
					sbr.WriteString(fmt.Sprintf("this._mem_util.block_release(p%d[0]);\n", i))
				default:
				}
			}
			fn.PreCall = sb.String()
			fn.Release = sbr.String()
		}

		// 返回值
		if len(f.Results) > 0 {
			var sb strings.Builder
			var sbr strings.Builder
			for i, r := range f.Results {
				switch tp := r.(type) {
				case *wir.U8, *wir.U16, *wir.I32, *wir.U32, *wir.I64, *wir.U64, *wir.Bool, *wir.Rune, *wir.String:
					sb.WriteString(fmt.Sprintf("let r%d = this._mem_util.extract_%s(res);\n", i, tp.Named()))
					if i > 0 {
						sbr.WriteString(",")
					}
					sbr.WriteString(fmt.Sprintf("r%d", i))
				default:
				}
			}
			fn.GetResults = sb.String()
			if len(f.Results) == 1 {
				fn.Return = "return " + sbr.String() + ";"
			} else {
				fn.Return = "return [" + sbr.String() + "];"
			}
		}

		funcs = append(funcs, fn)

	}
	return funcs
}

func (p *Compiler) GenJSBinding(wasmFilename string) string {
	// 模板
	t, err := template.New("js").Parse(js_binding_tmpl)
	if err != nil {
		logger.Fatal(err)
	}
	data := JSModule{
		Filename: wasmFilename,
		Pkg:      p.prog.Manifest.MainPkg,
		Globals:  p.globalsForJsBinding(),
		Funcs:    p.funcsForJSBinding(),
	}
	var bf bytes.Buffer
	err = t.Execute(&bf, data)
	if err != nil {
		logger.Fatal(err)
	}

	// 测试用：生成一个js文件
	// TODO: 正式版需要删除
	// os.WriteFile("js_binding_test.js", bf.Bytes(), 0644)

	return bf.String()
}

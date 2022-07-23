// 版权 @2019 凹语言 作者。保留所有权利。

package app

// 命令行选项
type Option struct {
	Debug      bool
	TargetArch string
	TargetOS   string
	Clang      string
	WasmLLC    string
	WasmLD     string
}

// 命令行程序对象
type App struct {
	opt  Option
	path string
	src  string
}

// 构建命令行程序对象
func NewApp(opt *Option) *App {
	panic("TODO")

}

func (p *App) InitApp(name, pkgpath string, update bool) error {
	panic("TODO")
}

func (p *App) Fmt(path string) error {
	panic("TODO")
}

func (p *App) Lex(filename string) error {
	return nil
}

func (p *App) AST(filename string) error {
	panic("TODO")
}

func (p *App) SSA(filename string) error {
	panic("TODO")
}

func (p *App) ASM(filename string) error {
	panic("TODO")
}

func (p *App) Build(filename string, src interface{}, outfile string) (output []byte, err error) {
	panic("TODO")
}

func (p *App) Run(filename string, src interface{}) (data []byte, err error) {
	panic("TODO")
}

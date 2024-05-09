// 版权 @2019 凹语言 作者。保留所有权利。

// 凹语言打包程序: Windows/Linux/macOS/wasm
package main

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"wa-lang.org/wa/internal/3rdparty/wabt-go"
	"wa-lang.org/wa/internal/version"
	"wa-lang.org/wa/waroot"
)

const (
	darwin  = "darwin"
	linux   = "linux"
	windows = "windows"

	amd64 = "amd64"
	arm64 = "arm64"
)

var (
	flagOutputDir = flag.String("output", "_output", "set output dir")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func main() {
	flag.Parse()

	b := &Builder{Output: *flagOutputDir}
	b.GenAll()
}

type Builder struct {
	Output string // 目标的根目录
}

func NewBuilder(outputDir string) *Builder {
	return &Builder{Output: outputDir}
}

func (p *Builder) GenAll() {
	waRoot_darwin_amd64 := p.getWarootPath(darwin, arm64)
	waRoot_darwin_arm64 := p.getWarootPath(darwin, amd64)
	waRoot_linux_amd64 := p.getWarootPath(linux, amd64)
	waRoot_windows_amd64 := p.getWarootPath(windows, amd64)

	p.genWarootFiles(waRoot_darwin_amd64)
	p.genWarootFiles(waRoot_darwin_arm64)
	p.genWarootFiles(waRoot_linux_amd64)
	p.genWarootFiles(waRoot_windows_amd64)

	p.genWat2wasmExe()
	p.genWaExe()

	p.zipDir(waRoot_darwin_amd64)
	p.zipDir(waRoot_darwin_arm64)
	p.zipDir(waRoot_linux_amd64)
	p.zipDir(waRoot_windows_amd64)

	os.RemoveAll(waRoot_darwin_amd64)
	os.RemoveAll(waRoot_darwin_arm64)
	os.RemoveAll(waRoot_linux_amd64)
	os.RemoveAll(waRoot_windows_amd64)

	p.genChecksums()
}

func (p *Builder) genChecksums() {
	paths := []string{
		p.getWarootPath(darwin, arm64),
		p.getWarootPath(darwin, amd64),
		p.getWarootPath(linux, amd64),
		p.getWarootPath(windows, amd64),
	}
	var buf bytes.Buffer
	for _, path := range paths {
		data, err := os.ReadFile(path + ".zip")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&buf, "%x MD5(%s)\n", md5.Sum(data), filepath.Base(path)+".zip")
	}

	os.WriteFile(
		filepath.Join(p.Output, "wa-"+version.Version+".checksums.txt"),
		buf.Bytes(),
		0666,
	)
}

func (p *Builder) genWaExe() {
	// darwin/arm64
	{
		waRootPath := p.getWarootPath(darwin, arm64)
		dstpath := filepath.Join(waRootPath, "bin", "wa")

		cmd := exec.Command("go", "build", "-o", dstpath, "wa-lang.org/wa")
		cmd.Env = append([]string{"GOOS=" + darwin, "GOARCH=" + arm64}, os.Environ()...)
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Print(string(output))
			panic(err)
		}
	}

	// darwin/amd64
	{
		waRootPath := p.getWarootPath(darwin, amd64)
		dstpath := filepath.Join(waRootPath, "bin", "wa")

		cmd := exec.Command("go", "build", "-o", dstpath, "wa-lang.org/wa")
		cmd.Env = append([]string{"GOOS=" + darwin, "GOARCH=" + amd64}, os.Environ()...)
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Print(string(output))
			panic(err)
		}
	}

	// ubuntu
	{
		waRootPath := p.getWarootPath(linux, amd64)
		dstpath := filepath.Join(waRootPath, "bin", "wa")

		cmd := exec.Command("go", "build", "-o", dstpath, "wa-lang.org/wa")
		cmd.Env = append([]string{"GOOS=" + linux, "GOARCH=" + amd64}, os.Environ()...)
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Print(string(output))
			panic(err)
		}
	}

	// windows
	{
		waRootPath := p.getWarootPath(windows, amd64)
		dstpath := filepath.Join(waRootPath, "bin", "wa.exe")

		cmd := exec.Command("go", "build", "-o", dstpath, "wa-lang.org/wa")
		cmd.Env = append([]string{"GOOS=" + windows, "GOARCH=" + amd64}, os.Environ()...)
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Print(string(output))
			panic(err)
		}
	}
}

func (p *Builder) genWat2wasmExe() {
	// macos/arm64
	{
		waRootPath := p.getWarootPath(darwin, arm64)
		dstpath := filepath.Join(waRootPath, "bin", wabt.Wat2WasmName)
		os.WriteFile(dstpath, []byte(wabt.Wat2wasm_macos), 0777)
	}

	// macos/amd64
	{
		waRootPath := p.getWarootPath(darwin, amd64)
		dstpath := filepath.Join(waRootPath, "bin", wabt.Wat2WasmName)
		os.WriteFile(dstpath, []byte(wabt.Wat2wasm_macos), 0777)
	}

	// ubuntu
	{
		waRootPath := p.getWarootPath(linux, amd64)
		dstpath := filepath.Join(waRootPath, "bin", wabt.Wat2WasmName)
		os.WriteFile(dstpath, []byte(wabt.Wat2wasm_ubuntu), 0777)
	}

	// windows
	{
		waRootPath := p.getWarootPath(windows, amd64)
		dstpath := filepath.Join(waRootPath, "bin", wabt.Wat2WasmName)
		os.WriteFile(dstpath, []byte(wabt.Wat2wasm_windows), 0777)
	}
}

func (p *Builder) genWarootFiles(waRootPath string) error {
	os.RemoveAll(waRootPath)
	os.MkdirAll(waRootPath, 0777)
	warootfs := waroot.GetRootFS()
	err := fs.WalkDir(waroot.GetRootFS(), ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		data, err := fs.ReadFile(warootfs, path)
		if err != nil {
			return err
		}

		dstpath := filepath.Join(waRootPath, path)
		os.MkdirAll(filepath.Dir(dstpath), 0777)

		if s := filepath.Base(path); s == "_keep" || s == ".keep" {
			return nil
		}

		f, err := os.Create(dstpath)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.Write(data); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Builder) zipDir(dir string) {
	file, err := os.Create(dir + ".zip")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		relpath, _ := filepath.Rel(dir, path)
		f, err := w.Create(filepath.Join("wa", relpath))
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (p *Builder) getWarootPath(waos, waarch string) string {
	return fmt.Sprintf("%s/wa_%s_%s-%s", p.Output, version.Version, waos, waarch)
}

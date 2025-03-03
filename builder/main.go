// 版权 @2024 凹语言 作者。保留所有权利。

// go1.21 support wasip1/wasm
// go1.24 support wasmexport

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
	"runtime"
	"strings"

	"wa-lang.org/wa/internal/3rdparty/gover"
	"wa-lang.org/wa/internal/version"
	"wa-lang.org/wa/waroot"
)

const (
	darwin  = "darwin"
	linux   = "linux"
	windows = "windows"
	wasip1  = "wasip1"

	amd64 = "amd64"
	arm64 = "arm64"
	wasm  = "wasm"
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
	waRoot_darwin_amd64 := p.getWarootPath(darwin, amd64)
	waRoot_darwin_arm64 := p.getWarootPath(darwin, arm64)
	waRoot_linux_amd64 := p.getWarootPath(linux, amd64)
	waRoot_windows_amd64 := p.getWarootPath(windows, amd64)
	waRoot_wasip1_wasm := p.getWarootPath(wasip1, wasm)

	waRoot_docker := fmt.Sprintf("%s/wa-docker-linux-amd64", p.Output)
	waRoot_wasip1 := fmt.Sprintf("%s/wa-wasip1", p.Output)

	p.genWarootFiles(waRoot_darwin_amd64)
	p.genWarootFiles(waRoot_darwin_arm64)
	p.genWarootFiles(waRoot_linux_amd64)
	p.genWarootFiles(waRoot_windows_amd64)

	if isWasip1Enabled() {
		p.genWarootFiles(waRoot_wasip1_wasm)
	}

	p.genWaExe()

	os.RemoveAll(waRoot_docker)
	cpDir(waRoot_docker, waRoot_linux_amd64)

	if isWasip1Enabled() {
		os.RemoveAll(waRoot_wasip1)
		cpDir(waRoot_wasip1, waRoot_wasip1_wasm)
	}

	p.zipDir(waRoot_darwin_amd64)
	p.zipDir(waRoot_darwin_arm64)
	p.zipDir(waRoot_linux_amd64)
	p.zipDir(waRoot_windows_amd64)

	if isWasip1Enabled() {
		p.zipDir(waRoot_wasip1_wasm)
	}

	os.RemoveAll(waRoot_darwin_amd64)
	os.RemoveAll(waRoot_darwin_arm64)
	os.RemoveAll(waRoot_linux_amd64) // keep for build docker image
	os.RemoveAll(waRoot_windows_amd64)

	if isWasip1Enabled() {
		os.RemoveAll(waRoot_wasip1_wasm)
	}

	p.genChecksums()
}

func (p *Builder) genChecksums() {
	paths := []string{
		p.getWarootPath(darwin, arm64),
		p.getWarootPath(darwin, amd64),
		p.getWarootPath(linux, amd64),
		p.getWarootPath(windows, amd64),
	}
	if isWasip1Enabled() {
		paths = append(paths, p.getWarootPath(wasip1, wasm))
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
	CGO_ENABLED := "CGO_ENABLED=0"

	// wasip1/wasm
	if isWasip1Enabled() {
		waRootPath := p.getWarootPath(wasip1, wasm)
		dstpath := filepath.Join(waRootPath, "bin", "wa.wasm")

		cmd := exec.Command("go", "build", "-o", dstpath, "wa-lang.org/wa")
		cmd.Env = append([]string{"GOOS=" + wasip1, "GOARCH=" + wasm, CGO_ENABLED}, os.Environ()...)
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Print(string(output))
			panic(err)
		}
	}

	// darwin/arm64
	{
		waRootPath := p.getWarootPath(darwin, arm64)
		dstpath := filepath.Join(waRootPath, "bin", "wa")

		cmd := exec.Command("go", "build", "-o", dstpath, "wa-lang.org/wa")
		cmd.Env = append([]string{"GOOS=" + darwin, "GOARCH=" + arm64, CGO_ENABLED}, os.Environ()...)
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
		cmd.Env = append([]string{"GOOS=" + darwin, "GOARCH=" + amd64, CGO_ENABLED}, os.Environ()...)
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
		cmd.Env = append([]string{"GOOS=" + linux, "GOARCH=" + amd64, CGO_ENABLED}, os.Environ()...)
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
		cmd.Env = append([]string{"GOOS=" + windows, "GOARCH=" + amd64, CGO_ENABLED}, os.Environ()...)
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Print(string(output))
			panic(err)
		}
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

		// 跳过忽略的文件
		if s := filepath.Base(path); s == "_keep" || s == ".keep" {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
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

		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		// 需要保留可执行等标志信息
		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}

		relpath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		header.Name = filepath.Join("wa", relpath)
		header.Method = zip.Deflate

		f, err := w.CreateHeader(header)
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

func isWasip1Enabled() bool {
	// go1.24 support wasmexport
	goversion := strings.TrimPrefix(runtime.Version(), "go")
	return gover.Compare(goversion, "1.24") >= 0
}

func cpDir(dst, src string) (total int) {
	entryList, err := os.ReadDir(src)
	if err != nil && !os.IsExist(err) {
		log.Fatal("cpDir: ", err)
	}
	for _, entry := range entryList {
		if entry.IsDir() {
			cpDir(dst+"/"+entry.Name(), src+"/"+entry.Name())
		} else {
			srcFname := filepath.Clean(src + "/" + entry.Name())
			dstFname := filepath.Clean(dst + "/" + entry.Name())

			cpFile(dstFname, srcFname)
			total++
		}
	}
	return
}

func cpFile(dst, src string) {
	err := os.MkdirAll(filepath.Dir(dst), 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal("cpFile: ", err)
	}
	fsrc, err := os.Open(src)
	if err != nil {
		log.Fatal("cpFile: ", err)
	}
	defer fsrc.Close()

	fi, err := fsrc.Stat()
	if err != nil {
		log.Fatal("cpFile: ", err)
	}

	fdst, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fi.Mode())
	if err != nil {
		log.Fatal("cpFile: ", err)
	}
	defer fdst.Close()
	if _, err = io.Copy(fdst, fsrc); err != nil {
		log.Fatal("cpFile: ", err)
	}
}

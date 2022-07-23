// 版权 @2019 凹语言 作者。保留所有权利。

package config

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// 模块文件
const WaModFile = "wa.mod.json"

// wa.json 文件结构
type Manifest struct {
	Root    string `json:"root"` // wa.json 所在目录
	MainPkg string `json:"main"` // 主包路径

	Pkg Manifest_package `json:"package"`
}

// 包基础信息
type Manifest_package struct {
	Name          string   `json:"name"`                    // 名字
	Pkgpath       string   `json:"pkgpath"`                 // 模块的导入路径
	Version       string   `json:"version"`                 // 版本
	Authors       []string `json:"authors,omitempty"`       // 作者
	Description   string   `json:"description,omitempty"`   // 一句话简介
	Documentation string   `json:"documentation,omitempty"` // 包文档链接
	Readme        string   `json:"readme,omitempty"`        // README 文件 (Markdown 格式)
	Homepage      string   `json:"homepage,omitempty"`      // 主页
	Repository    string   `json:"repository,omitempty"`    // 代码仓库
	License       string   `json:"license,omitempty"`       // 版权
	LicenseFile   string   `json:"license_file,omitempty"`  // 版权文件
	Keywords      []string `json:"keywords,omitempty"`      // 关键字
	Categories    []string `json:"categories,omitempty"`    // 领域分类
}

func (p *Manifest) Clone() *Manifest {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(p)
	if err != nil {
		panic(err)
	}
	var v Manifest
	err = gob.NewDecoder(&buf).Decode(&v)
	if err != nil {
		panic(err)
	}
	return &v
}

// 加载 wa.json 文件
func LoadManifest(appPath string) (*Manifest, error) {
	workDir := appPath
	if workDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("loader.LoadManifest: workDir is empty")
		}
		workDir = wd
	}

	// 查找 wa.json 路径
	kManifestPath, err := findManifestPath(workDir)
	if err != nil {
		return nil, fmt.Errorf("loader.LoadManifest: find wa.json failed : %w", err)
	}

	// 读取 wa.json 文件
	data, err := os.ReadFile(kManifestPath)
	if err != nil {
		return nil, fmt.Errorf("loader.LoadManifest: read %s failed: %w", kManifestPath, err)
	}

	// 解码 JSON
	p := new(Manifest)
	if err := json.Unmarshal(data, &p.Pkg); err != nil {
		return nil, fmt.Errorf("loader.LoadManifest: json.Unmarshal %s failed: %w", kManifestPath, err)
	}

	// 当前 app 默认目录
	p.Root = filepath.Dir(kManifestPath)
	p.MainPkg, _ = filepath.Rel(p.Root, workDir)

	if p.MainPkg == "" || p.MainPkg == "." {
		p.MainPkg = p.Pkg.Pkgpath
	}

	return p, nil
}

func (p *Manifest) JSONString() string {
	d, _ := json.MarshalIndent(p, "", "\t")
	return string(d)
}

// 查找 wa.json 文件路径
func findManifestPath(pkgpath string) (string, error) {
	if pkgpath == "" {
		return "", fmt.Errorf("empty pkgpath")
	}

	// 依次向上查找 wa.json
	pkgroot := pkgpath
	for pkgroot != "" {
		kManifestPath := filepath.Join(pkgroot, WaModFile)
		if fi, _ := os.Stat(kManifestPath); fi != nil {
			return kManifestPath, nil
		}

		// 是否已经到根目录
		parentDir := filepath.Dir(pkgroot)
		if parentDir == pkgroot {
			break
		}

		pkgroot = parentDir
	}

	// 查找失败
	return "", fmt.Errorf("%s: <wa.json> not found", pkgpath)
}

// 版权 @2019 凹语言 作者。保留所有权利。

package config

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// 模块文件
const WaModFile = "wa.mod.json"

// WaModFile 文件结构
type Manifest struct {
	Root    string `json:"root"` // WaModFile 所在目录
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

// 简版 Manifest
func SimpleManifest(mainPkg, appName string) *Manifest {
	p := &Manifest{
		MainPkg: mainPkg,
		Pkg: Manifest_package{
			Name:    appName,
			Pkgpath: mainPkg,
		},
	}
	return p
}

// 加载 WaModFile 文件
// 如果 vfs 为空则从本地文件系统读取
func LoadManifest(vfs fs.FS, appPath string) (p *Manifest, err error) {
	if appPath == "" {
		return nil, fmt.Errorf("loader.LoadManifest: appPath is empty")
	}

	// 查找 WaModFile 路径
	kManifestPath, err := findManifestPath(vfs, appPath)
	if err != nil {
		return nil, fmt.Errorf("loader.LoadManifest: find '%s' failed : %w", WaModFile, err)
	}

	// 读取 WaModFile 文件
	var data []byte
	if vfs != nil {
		data, err = fs.ReadFile(vfs, kManifestPath)
	} else {
		data, err = os.ReadFile(kManifestPath)
	}
	if err != nil {
		return nil, fmt.Errorf("loader.LoadManifest: read %s failed: %w", kManifestPath, err)
	}

	// 解码 JSON
	p = new(Manifest)
	if err := json.Unmarshal(data, &p.Pkg); err != nil {
		return nil, fmt.Errorf("loader.LoadManifest: json.Unmarshal %s failed: %w", kManifestPath, err)
	}

	// 当前 app 默认目录
	p.Root = filepath.Dir(kManifestPath)
	p.MainPkg, _ = filepath.Rel(p.Root, appPath)

	if p.MainPkg == "" || p.MainPkg == "." {
		p.MainPkg = p.Pkg.Pkgpath
	}

	return p, nil
}

func (p *Manifest) JSONString() string {
	d, _ := json.MarshalIndent(p, "", "\t")
	return string(d)
}

// 查找 WaModFile 文件路径
func findManifestPath(vfs fs.FS, pkgpath string) (string, error) {
	if pkgpath == "" {
		return "", fmt.Errorf("empty pkgpath")
	}

	// 依次向上查找 WaModFile
	pkgroot := pkgpath
	for pkgroot != "" {
		kManifestPath := filepath.Join(pkgroot, WaModFile)
		if vfs != nil {
			if fi, _ := fs.Stat(vfs, kManifestPath); fi != nil {
				return kManifestPath, nil
			}
		} else {
			if fi, _ := os.Stat(kManifestPath); fi != nil {
				return kManifestPath, nil
			}
		}

		// 是否已经到根目录
		parentDir := filepath.Dir(pkgroot)
		if parentDir == pkgroot {
			break
		}

		pkgroot = parentDir
	}

	// 查找失败
	return "", fmt.Errorf("%s: '%s' not found", WaModFile, pkgpath)
}

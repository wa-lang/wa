# 版权 @{{.Year}} {{.Name}} 作者。保留所有权利。

name = "{{.Name}}"
pkgpath = "{{.Pkgpath}}"
target = {{if .IsWasiApp}}"wasi"{{else if .IsWasm4App}}"wasm4"{{else if .IsArduinoApp}}"arduino"{{else}}"js"{{end}}

// 版权 @2022 凹语言 作者。保留所有权利。

//go:build !cgo
// +build !cgo

package wabt

func Wat2WasmCmd(args ...string)  {
	exe := "wat2wasm"
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}
	cmd := exec.Command(exe, args...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		if len(stdoutStderr) != 0 {
			fmt.Printf("%s\n", stdoutStderr)
		}
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", stdoutStderr)
}

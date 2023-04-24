// 版权 @2022 凹语言 作者。保留所有权利。

package main

import "wa-lang.org/wa/api"

var waError error

func waGetError() error {
	return waError
}
func waGetErrorText() string {
	if waError != nil {
		return waError.Error()
	} else {
		return ""
	}
}
func waSetError(err error) {
	waError = err
}
func waClearError() {
	waError = nil
}

func waGenerateWat(filename, code string) string {
	cfg := api.DefaultConfig()
	cfg.WaArch = api.WaArch_wasm
	cfg.WaOS = api.WaOS_chrome

	wat, err := api.BuildFile(cfg, filename, code)
	if err != nil {
		if waGetError() == nil {
			waSetError(err)
		}
		return ""
	}
	return string(wat)
}

func waFormatCode(filename, code string) string {
	newCode, err := api.FormatCode(filename, code)
	if err != nil {
		if waGetError() == nil {
			waSetError(err)
		}
		return code
	}
	return newCode
}

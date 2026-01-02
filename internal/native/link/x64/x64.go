package x64

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func AssembleAndLink(inputPath string, outputPath string) error {
	absInput, _ := filepath.Abs(inputPath)
	absOutput, _ := filepath.Abs(outputPath)

	var objPath string
	var nasmFormat string
	var gccArgs []string

	switch runtime.GOOS {
	case "windows":
		nasmFormat = "win64"
		objPath = strings.TrimSuffix(absInput, filepath.Ext(absInput)) + ".obj"
		gccArgs = []string{"-nostdlib", "-static", objPath, "-o", absOutput, "-Wl,-e,main"}
	case "linux":
		nasmFormat = "elf64"
		objPath = strings.TrimSuffix(absInput, filepath.Ext(absInput)) + ".o"
		gccArgs = []string{"-nostdlib", "-static", objPath, "-o", absOutput, "-Wl,-e,_start", "-z", "noexecstack"}
	default:
		panic("unreachable")
	}

	nasmCmd := exec.Command("nasm", "-f", nasmFormat, absInput, "-o", objPath)
	if output, err := nasmCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("NASM ERR:\n%s", string(output))
	}
	defer os.Remove(objPath)

	gccCmd := exec.Command("gcc", gccArgs...)
	if output, err := gccCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("LINK ERR:\n%s", string(output))
	}

	return nil
}

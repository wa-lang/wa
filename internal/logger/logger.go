// 版权 @2021 凹语言 作者。保留所有权利。

package logger

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var DebugMode = false

func CallerInfo(skip int) (fn, filename string, line int) {
	return callerInfo(skip + 1)
}

func CallStack() string {
	var buf bytes.Buffer
	for skip := 1; ; skip += 1 {
		fn, filename, line := callerInfo(skip)
		if filename == "" {
			break
		}
		fmt.Fprintf(&buf, "%s:%d: %s", filename, line, fn)
	}
	return buf.String()
}

func callerInfo(skip int) (fn, filename string, line int) {
	pc, filename, line, _ := runtime.Caller(skip + 1)
	fn = runtime.FuncForPC(pc).Name()
	if idx := strings.LastIndex(fn, "/"); idx >= 0 {
		fn = fn[idx+1:]
	}
	if wd, _ := os.Getwd(); wd != "" {
		if rel, err := filepath.Rel(wd, filename); err == nil {
			filename = rel
		}
	}
	return
}

func Assert(ok bool, a ...interface{}) {
	if ok {
		return
	}
	fn, filename, line := callerInfo(1)
	msg := fmt.Sprint(a...)
	if msg != "" {
		panic(fmt.Sprintf("%s:%d: %s: assert failed: %v", filename, line, fn, msg))
	} else {
		panic(fmt.Sprintf("%s:%d: %s: assert failed", filename, line, fn))
	}
}

func Assertf(ok bool, format string, a ...interface{}) {
	if ok {
		return
	}
	fn, filename, line := callerInfo(1)
	msg := fmt.Sprintf(format, a...)
	if msg != "" {
		panic(fmt.Sprintf("%s:%d: %s: assert failed: %v", filename, line, fn, msg))
	} else {
		panic(fmt.Sprintf("%s:%d: %s: assert failed", filename, line, fn))
	}
}

func AssertEQ(a, b interface{}) {
	if a == b {
		return
	}
	fn, filename, line := callerInfo(1)
	panic(fmt.Sprintf("%s:%d: %s: AssertEQ: %v != %v", filename, line, fn, a, b))
}

func Debug(a ...interface{}) {
	if DebugMode {
		fn, filename, line := callerInfo(1)
		msg := fmt.Sprint(a...)
		fmt.Printf("%s:%d: %s: %s", filename, line, fn, msg)
	}
}

func Debugln(a ...interface{}) {
	if DebugMode {
		fn, filename, line := callerInfo(1)
		msg := fmt.Sprintln(a...)
		fmt.Printf("%s:%d: %s: %s", filename, line, fn, msg)
	}
}

func Debugf(format string, a ...interface{}) {
	if DebugMode {
		fn, filename, line := callerInfo(1)
		msg := fmt.Sprintf(format, a...)
		fmt.Printf("%s:%d: %s: %s", filename, line, fn, msg)
	}
}

func Fatal(a ...interface{}) {
	_, filename, line := callerInfo(1)
	msg := fmt.Sprint(a...)
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Printf("%s:%d: %s", filename, line, msg)
	os.Exit(1)
}

func Fatalf(format string, a ...interface{}) {
	_, filename, line := callerInfo(1)
	msg := fmt.Sprintf(format, a...)
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Printf("%s:%d: %s", filename, line, msg)
	os.Exit(1)
}

func Print(a ...interface{}) {
	_, filename, line := callerInfo(1)
	msg := fmt.Sprint(a...)
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Printf("%s:%d: %s", filename, line, msg)
}

func Printf(format string, a ...interface{}) {
	_, filename, line := callerInfo(1)
	msg := fmt.Sprintf(format, a...)
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Printf("%s:%d: %s", filename, line, msg)
}

func Panic(a ...interface{}) {
	msg := fmt.Sprint(a...)
	panic(msg)
}

func Panicf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	panic(msg)
}

func Trace(flagEnabled *bool, a ...interface{}) {
	if flagEnabled == nil {
		flagEnabled = &DebugMode
	}
	if *flagEnabled {
		fn, filename, line := callerInfo(1)
		if len(a) == 0 {
			a = append(a, fn)
		}
		msg := fmt.Sprintln(a...)
		msg = withPrefix("\t", msg)
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Printf("trace %s:%d: func %s\n%s", filename, line, fn, msg)
	}
}

func Tracef(flagEnabled *bool, format string, a ...interface{}) {
	if flagEnabled == nil {
		flagEnabled = &DebugMode
	}
	if *flagEnabled {
		fn, filename, line := callerInfo(1)
		if len(a) == 0 {
			a = append(a, fn)
		}
		msg := fmt.Sprintf(format, a...)
		msg = withPrefix("\t", msg)
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Printf("trace %s:%d: func %s\n%s", filename, line, fn, msg)
	}
}

func DumpFS(flagEnabled *bool, name string, fileSystem fs.FS, root string) {
	if flagEnabled == nil {
		flagEnabled = &DebugMode
	}
	if *flagEnabled {
		fmt.Printf("dumpfs: fileSystem: %s\n", name)
		fs.WalkDir(fileSystem, root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("\t%s: %v\n", path, err)
				return err
			}
			if d.IsDir() {
				fmt.Printf("\t%s # %s\n", path, "dir")
			} else {
				fmt.Printf("\t%s # %s\n", path, "file")
			}
			return nil
		})
	}
}

func withPrefix(prefix, message string) string {
	if message == "" {
		return message
	}
	if prefix == "" {
		return message
	}
	lines := strings.Split(message, "\n")
	for i, s := range lines {
		lines[i] = prefix + s
	}
	return strings.Join(lines, "\n")
}

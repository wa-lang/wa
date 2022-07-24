package syscall

// https://opensource.apple.com/source/xnu/xnu-2782.20.48/bsd/kern/syscalls.master
type fn string

const (
	EXIT  fn = "EXIT"
	WRITE fn = "WRITE"
)

var convDarwin = map[fn]int64{
	EXIT:  0x2000001,
	WRITE: 0x2000004,
}

var convLinux = map[fn]int64{
	EXIT:  60,
	WRITE: 1,
}

func Convert(f fn, goos string) int64 {
	switch goos {
	case "darwin":
		return convDarwin[f]
	case "linux":
		return convLinux[f]
	default:
		panic("unknown goos")
	}
}

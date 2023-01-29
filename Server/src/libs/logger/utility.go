package logger

import (
	"runtime"
)

func Recover() {
	if r := recover(); r != nil {
		Fatal("[PANIC]%s", r)
		buf := make([]byte, 4096)
		runtime.Stack(buf[:], false)
		Fatal("[STACK]%s", string(buf))
	}
}


package logger

import (
	"testing"
)

//日志不分离
func TestLog1(t *testing.T) {
	Init("testlog", ".", 1000, 3, DEBUG_LEVEL, false, PUT_CONSOLE)
	for i := 0; i < 30; i++ {
		Debug("This a DENUG test log.")
		Info("This a INFO test log.")
		Error("This a ERROR test log.")
		Fatal("This a Fatal test log.")
	}
}

//日志分离
func TestLog2(t *testing.T) {
	Init("testlog", ".", 1000, 3, DEBUG_LEVEL, true, PUT_CONSOLE)
	for i := 0; i < 30; i++ {
		Debug("This a DENUG test log.")
		Info("This a INFO test log.")
		Error("This a ERROR test log.")
		Fatal("This a Fatal test log.")
	}
}

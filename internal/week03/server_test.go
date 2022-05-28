package week03

import "testing"

// 由于 os.Signal，目前报异常# runtime/cgo  _cgo_export.c:3:10: fatal error: stdlib.h: No such file or directory ... compilation terminated.
func TestRunServer(t *testing.T) {
	RunServer()
}

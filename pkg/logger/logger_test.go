package logger

import "testing"

func Test_DefaultInit(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Error("Shouldn't panic.", e)
		}
	}()
	Info("some body")
}

package ast

import (
	"testing"
)

func Test_FuncAddCall(t *testing.T) {
	f := ParseFile("./example_test.go")
	f.FuncAddCall("TT", "func", "测试", "a = 2", true)
	s := f.lineNode.String()
	t.Log(s)
}

func Test_StructAddField(t *testing.T) {
	f := ParseFile("./example_test.go")
	f.StructAddField([]string{"Config", "Sub", "Child"}, "注释", "Age int")
	f.StructAddField([]string{"Config", "Sub", "Child"}, "注释", "Ath struct")
	s := f.lineNode.String()
	t.Log(s)
}

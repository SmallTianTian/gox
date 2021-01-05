package ast_util_test

type _struct1 struct {
	A string
}

type _struct2 struct {
	Inner struct {
	}
}

var _var string

const num = 1

func Func() {}

func F1(name string) {
	F1(name)
}

func F2(num int, s _struct1) {
}

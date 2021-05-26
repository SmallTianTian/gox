package utils

import (
	"testing"
)

func TestFirstUp(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, ""},
		{"blank", args{" "}, " "},
		{"blank line", args{"\n"}, "\n"},
		{"sigle word", args{"a"}, "A"},
		{"word", args{"abc"}, "Abc"},
		{"blank with sigle word", args{" a"}, " A"},
		{"blank with word", args{" abc"}, " Abc"},
		{"sigle word with blank", args{"a "}, "A "},
		{"word with blank", args{"abc "}, "Abc "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FirstUp(tt.args.word); got != tt.want {
				t.Errorf("FirstUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegalVarName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "正常小写开头",
			args: args{"a"},
			want: true,
		},
		{
			name: "正常大写开头",
			args: args{"As"},
			want: true,
		},
		{
			name: "正常带下划线",
			args: args{"a_c"},
			want: true,
		},
		{
			name: "正常下划线结尾",
			args: args{"a_"},
			want: true,
		},
		{
			name: "正常带数字",
			args: args{"a1a"},
			want: true,
		},
		{
			name: "正常带数字结尾",
			args: args{"a1"},
			want: true,
		},
		{
			name: "正常全家福",
			args: args{"a1b_c"},
			want: true,
		},
		{
			name: "非正常下划线开头",
			args: args{"_a"},
			want: false,
		},
		{
			name: "非正常数字开头",
			args: args{"1a"},
			want: false,
		},
		{
			name: "非正常带中划线",
			args: args{"a-b"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LegalVarName(tt.args.s); got != tt.want {
				t.Errorf("LegalVarName() = %v, want %v", got, tt.want)
			}
		})
	}
}

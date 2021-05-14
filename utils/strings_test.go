package utils

import "testing"

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

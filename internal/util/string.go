package util

import (
	"strings"
	"unicode"
)

func FirstUp(word string) string {
	if strings.TrimSpace(word) == "" {
		return word
	}
	sb := strings.Builder{}
	for i, r := range word {
		if unicode.IsSpace(r) {
			sb.WriteRune(r)
			continue
		}
		sb.WriteRune(unicode.ToUpper(r))
		if i < len(word)-1 {
			sb.WriteString(word[i+1:])
		}
		break
	}
	return sb.String()
}

func LegalVarName(s string) bool {
	if len(s) == 0 {
		return false
	}
	isChar := func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
	}
	isNum := func(r rune) bool {
		return (r >= '0' && r <= '9')
	}
	isUnder := func(r rune) bool {
		return r == '_'
	}

	notFirstLegal := func(r rune) bool {
		return isChar(r) || isNum(r) || isUnder(r)
	}

	for i, r := range s {
		if i == 0 {
			if !isChar(r) {
				return false
			}
		}
		if !notFirstLegal(r) {
			return false
		}
	}
	return true
}

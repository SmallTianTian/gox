package utils

import (
	"bytes"
	"strings"
	"text/template"
	"unicode"
)

func StringFormat(text string, v interface{}) (string, error) {
	var bs bytes.Buffer
	t, err := template.New("").Parse(text)
	MustNotError(err)
	err = t.Execute(&bs, v)
	return bs.String(), err
}

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

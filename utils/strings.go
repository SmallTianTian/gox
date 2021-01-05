package utils

import (
	"bytes"
	"text/template"
)

func StringFormat(text string, v interface{}) (string, error) {
	var bs bytes.Buffer
	t, err := template.New("").Parse(text)
	MustNotError(err)
	err = t.Execute(&bs, v)
	return bs.String(), err
}

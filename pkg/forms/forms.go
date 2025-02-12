package forms

import (
	"net/url"
	"strings"
)

type Form struct {
	url.Values
}

func New(data url.Values) *Form {
	return &Form{
		data,
	}
}

func (f *Form) Required(fields ...string) bool {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			return false
		}
	}

	return true
}

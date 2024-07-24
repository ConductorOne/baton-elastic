package connector

import "strings"

type Anonymous struct {
	Data string
}

func (e *Anonymous) TrimPrefix(prefix string) *Anonymous {
	return &Anonymous{
		Data: strings.TrimPrefix(e.Data, prefix),
	}
}

func (e *Anonymous) TrimSuffix(suffix string) *Anonymous {
	return &Anonymous{
		Data: strings.TrimSuffix(e.Data, suffix),
	}
}

func (e *Anonymous) Split(sep string) []string {
	return strings.Split(e.Data, sep)
}

func (e *Anonymous) ToString() string {
	return e.Data
}

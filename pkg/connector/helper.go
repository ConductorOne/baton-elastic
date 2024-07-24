package connector

import "strings"

type Utility struct {
	Data string
}

func (e *Utility) TrimPrefix(prefix string) *Utility {
	return &Utility{
		Data: strings.TrimPrefix(e.Data, prefix),
	}
}

func (e *Utility) TrimSuffix(suffix string) *Utility {
	return &Utility{
		Data: strings.TrimSuffix(e.Data, suffix),
	}
}

func (e *Utility) Split(sep string) []string {
	return strings.Split(e.Data, sep)
}

func (e *Utility) ToString() string {
	return e.Data
}

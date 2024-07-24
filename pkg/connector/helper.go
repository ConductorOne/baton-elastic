package connector

import "strings"

type Utility struct {
	Data string
}

func (u *Utility) TrimPrefix(prefix string) *Utility {
	return &Utility{
		Data: strings.TrimPrefix(u.Data, prefix),
	}
}

func (u *Utility) TrimSuffix(suffix string) *Utility {
	return &Utility{
		Data: strings.TrimSuffix(u.Data, suffix),
	}
}

func (u *Utility) Split(sep string) []string {
	return strings.Split(u.Data, sep)
}

func (u *Utility) ToString() string {
	return u.Data
}

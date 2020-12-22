package mildew

import (
	"github.com/caffix/stringset"
)

type Mildew struct {
	Subs stringset.Set
}

func NewMildew() *Mildew {
	return &Mildew{
		Subs: stringset.New(),
	}
}

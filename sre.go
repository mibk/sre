package sre

import "github.com/mibk/sre/internal/syntax"

type Rx struct {
	gr *syntax.Group
}

func MustCompile(str string) *Rx {
	gr, err := syntax.Compile(str)
	if err != nil {
		panic(err)
	}
	return &Rx{gr}
}

func (rx *Rx) Match(b []byte) bool { return rx.gr.Match(b) }

package sre

import "github.com/mibk/sre/internal/syntax"

type Rx struct {
	prog *syntax.Prog
}

func MustCompile(str string) *Rx {
	prog, err := syntax.Compile(str)
	if err != nil {
		panic(err)
	}
	return &Rx{prog}
}

func (rx *Rx) Match(b []byte) bool { return rx.prog.Match(b) }

package sre

import "github.com/mibk/sre/internal/syntax"

type Rx struct {
	expr syntax.Expr
}

func MustCompile(str string) *Rx {
	expr, err := syntax.Compile(str)
	if err != nil {
		panic(err)
	}
	return &Rx{expr}
}

func (rx *Rx) Match(b []byte) bool {
	n, ok := rx.expr.Consume(b)
	return ok && n == len(b)
}

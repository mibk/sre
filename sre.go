package sre

import "github.com/mibk/sre/internal/syntax"

type Rx struct {
	expr syntax.Expr
}

func Compile(str string) (*Rx, error) {
	expr, err := syntax.Compile(str)
	return &Rx{expr}, err
}

func MustCompile(str string) *Rx {
	rx, err := Compile(str)
	if err != nil {
		panic(err)
	}
	return rx
}

func (rx *Rx) Match(b []byte) bool {
	n, ok := rx.expr.Consume(b)
	return ok && n == len(b)
}

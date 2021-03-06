package syntax

import (
	"errors"
	"unicode/utf8"
)

type Expr interface {
	Consume(b []byte) (n int, ok bool)
}

type Group struct {
	Exprs []Expr
}

func (g *Group) add(expr Expr) {
	g.Exprs = append(g.Exprs, expr)
}

func (g *Group) repeatLastExpr(min, max int) error {
	if len(g.Exprs) == 0 {
		return errors.New("missing argument to repetition")
	}
	g.Exprs[len(g.Exprs)-1] = NewRepetition(g.Exprs[len(g.Exprs)-1], min, max)
	return nil
}

func (g *Group) Consume(b []byte) (n int, ok bool) {
	var m int
	for _, e := range g.Exprs {
		n, ok := e.Consume(b)
		if !ok {
			return 0, false
		}
		b = b[n:]
		m += n
	}
	return m, true
}

type OrGroup struct {
	Lhs Expr
	Rhs Expr
}

func (og OrGroup) Consume(b []byte) (n int, ok bool) {
	if n, ok = og.Lhs.Consume(b); ok {
		return n, ok
	}
	return og.Rhs.Consume(b)
}

type Char rune

func (c Char) Consume(b []byte) (n int, ok bool) {
	r, size := utf8.DecodeRune(b)
	if rune(c) == r {
		return size, true
	}
	return 0, false
}

type Any struct{}

func (Any) Consume(b []byte) (n int, ok bool) {
	_, size := utf8.DecodeRune(b)
	return size, size > 0
}

type CharSet struct {
	Neg   bool
	Exprs []Expr
}

var any = Any{}

func (cs CharSet) Consume(b []byte) (n int, ok bool) {
	for _, e := range cs.Exprs {
		if n, ok := e.Consume(b); ok {
			if cs.Neg {
				return 0, false
			}
			return n, true
		}
	}
	if cs.Neg {
		return any.Consume(b)
	}
	return 0, false
}

type Range struct {
	From, To rune
}

func (ra Range) Consume(b []byte) (n int, ok bool) {
	r, size := utf8.DecodeRune(b)
	if r >= ra.From && r <= ra.To {
		return size, true
	}
	return 0, false
}

const Unlimited = -1

type Repetition struct {
	Expr     Expr
	Min, Max int
}

func NewRepetition(e Expr, min, max int) Repetition {
	return Repetition{e, min, max}
}

func (r Repetition) Consume(b []byte) (n int, ok bool) {
	var i int
	var m int
	for ; ; i++ {
		n, ok := r.Expr.Consume(b)
		if !ok {
			break
		}
		b = b[n:]
		m += n
	}
	if (r.Min == Unlimited || i >= r.Min) &&
		(r.Max == Unlimited || i <= r.Max) {
		return m, true
	}
	return 0, false
}

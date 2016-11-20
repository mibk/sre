package syntax

import (
	"errors"
	"unicode/utf8"
)

type Group struct {
	Exprs []Expr
}

func (g *Group) add(expr Expr) {
	g.Exprs = append(g.Exprs, expr)
}

func (g *Group) last() (*Expr, error) {
	if len(g.Exprs) == 0 {
		return nil, errors.New("missing argument to quantifier")
	}
	return &g.Exprs[len(g.Exprs)-1], nil
}

func (g *Group) Match(b []byte) bool {
	for _, e := range g.Exprs {
		n, ok := e.Consume(b)
		if !ok {
			return false
		}
		b = b[n:]
	}
	return len(b) == 0
}

type Expr interface {
	Consume(b []byte) (n int, ok bool)
}

type Char rune

func (c Char) Consume(b []byte) (n int, ok bool) {
	r, size := utf8.DecodeRune(b)
	if rune(c) == rune(r) {
		return size, true
	}
	return 0, false
}

type Any struct{}

func (a Any) Consume(b []byte) (n int, ok bool) {
	_, size := utf8.DecodeRune(b)
	return size, size > 0
}

type CharSet struct {
	Chars []Char
}

func (cs CharSet) Consume(b []byte) (n int, ok bool) {
	for _, c := range cs.Chars {
		if n, ok := c.Consume(b); ok {
			return n, ok
		}
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

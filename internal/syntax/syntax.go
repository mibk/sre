package syntax

import (
	"errors"
	"unicode/utf8"
)

type Group struct {
	Exprs []Expr
}

func Compile(str string) (*Group, error) {
	var gr Group
	l := newLexer(str)
	for {
		r := l.Next()
		if r == EOF {
			break
		} else if r == Error {
			return nil, l.Err()
		}

		switch r {
		case QuestMark, Mul, Plus:
			last, err := gr.last()
			if err != nil {
				return nil, err
			}
			min, max := 0, 1
			switch r {
			case Plus:
				min = 1
				fallthrough
			case Mul:
				max = Unlimited
			}
			*last = NewQuantifier(*last, min, max)
		case Dot:
			gr.add(Any{})
		default:
			gr.add(Char(r))
		}
	}
	return &gr, nil
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

const Unlimited = -1

type Quantifier struct {
	Expr     Expr
	Min, Max int
}

func NewQuantifier(e Expr, min, max int) Quantifier {
	return Quantifier{e, min, max}
}

func (q Quantifier) Consume(b []byte) (n int, ok bool) {
	var i int
	var m int
	for ; ; i++ {
		n, ok := q.Expr.Consume(b)
		if !ok {
			break
		}
		b = b[n:]
		m += n
	}
	if (q.Min == Unlimited || i >= q.Min) &&
		(q.Max == Unlimited || i <= q.Max) {
		return m, true
	}
	return 0, false
}

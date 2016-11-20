package syntax

import "unicode/utf8"

type Prog struct {
	exprs []Consumer
}

func Compile(str string) (*Prog, error) {
	var prog Prog
	var lexpr Consumer

	l := newLexer(str)
	for {
		r := l.Next()
		if r == EOF {
			break
		} else if r == Error {
			return nil, l.Err()
		}

		var expr Consumer
		switch r {
		case QuestMark:
			lexpr = NewQuantifier(lexpr, 0, 1)
			continue
		case Mul:
			lexpr = NewQuantifier(lexpr, 0, Unlimited)
			continue
		case Plus:
			lexpr = NewQuantifier(lexpr, 1, Unlimited)
			continue
		case Dot:
			expr = Any{}
		default:
			expr = Char(r)
		}
		if lexpr != nil {
			prog.exprs = append(prog.exprs, lexpr)
		}
		lexpr = expr
	}
	if lexpr != nil {
		prog.exprs = append(prog.exprs, lexpr)
	}
	return &prog, nil
}

func (p *Prog) Match(b []byte) bool {
	for _, e := range p.exprs {
		n, ok := e.Consume(b)
		if !ok {
			return false
		}
		b = b[n:]
	}
	return len(b) == 0
}

type Consumer interface {
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

const (
	Unlimited = -1
)

type Quantifier struct {
	expr     Consumer
	min, max int
}

func NewQuantifier(e Consumer, min, max int) Quantifier {
	return Quantifier{e, min, max}
}

func (q Quantifier) Consume(b []byte) (n int, ok bool) {
	var i int
	var m int
	for ; ; i++ {
		n, ok := q.expr.Consume(b)
		if !ok {
			break
		}
		b = b[n:]
		m += n
	}
	if (q.min == Unlimited || i >= q.min) &&
		(q.max == Unlimited || i <= q.max) {
		return m, true
	}
	return 0, false
}

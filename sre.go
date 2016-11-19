package sre

type Rx struct {
	exprs []Consumer
}

func MustCompile(str string) *Rx {
	var rx Rx
	var lexpr Consumer
	for _, r := range str {
		var expr Consumer
		switch r {
		case '+':
			lexpr = NewQuantifier(lexpr, 1, Unlimited)
			continue
		default:
			expr = Lit(r)
		}
		if lexpr != nil {
			rx.exprs = append(rx.exprs, lexpr)
		}
		lexpr = expr
	}
	if lexpr != nil {
		rx.exprs = append(rx.exprs, lexpr)
	}
	return &rx
}

func (rx *Rx) Match(b []byte) bool {
	// TODO: not ^$
	for _, e := range rx.exprs {
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

type Lit rune

func (l Lit) Consume(b []byte) (n int, ok bool) {
	// TODO: utf8
	if len(b) > 0 && rune(l) == rune(b[0]) {
		return 1, true
	}
	return 0, false
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

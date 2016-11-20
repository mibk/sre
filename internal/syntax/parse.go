package syntax

import "errors"

type parser struct {
	lex    *lexer
	err    error
	peeked *rune
}

func Compile(str string) (Expr, error) {
	l := newLexer(str)
	p := &parser{lex: l}
	expr := p.parseGroup(false)
	if p.err != nil {
		return nil, p.err
	}
	return expr, nil
}

func (p *parser) peek() rune {
	if p.err != nil {
		return EOF
	}
	if r := p.peeked; r != nil {
		return *r
	}
	r := p.lex.Next()
	if r == Error {
		p.err = p.lex.Err()
		return EOF
	}
	p.peeked = &r
	return r
}

func (p *parser) next() rune {
	r := p.peek()
	p.peeked = nil
	return r
}

func (p *parser) saveErr(err error) {
	if p.err == nil {
		p.err = err
	}
}

func (p *parser) parseGroup(sub bool) Expr {
	gr := new(Group)
	var expr Expr = gr
Loop:
	for {
		switch r := p.next(); r {
		case EOF:
			break Loop
		case QuestMark, Mul, Plus:
			last, err := gr.last()
			if err != nil {
				p.saveErr(err)
				return nil
			}
			min, max := 0, 1
			switch r {
			case Plus:
				min = 1
				fallthrough
			case Mul:
				max = Unlimited
			}
			*last = NewRepetition(*last, min, max)
		case Dot:
			gr.add(Any{})
		case LBracket:
			gr.add(p.parseCharSet())
		case Or:
			gr = new(Group)
			expr = &OrGroup{expr, gr}
		case LParen:
			gr.add(p.parseGroup(true))
		case RParen:
			if sub {
				break Loop
			}
			p.saveErr(errors.New("unexpected )"))
			return nil
		default:
			gr.add(Char(r))
		}
	}
	return expr
}

func (p *parser) parseCharSet() *CharSet {
	var cs CharSet
	if r := p.peek(); r == Caret {
		p.next()
		cs.Neg = true
	}
Loop:
	for {
		switch r := p.next(); r {
		case EOF:
			p.saveErr(errors.New("unexpected end of character set"))
			return nil
		case RBracket:
			break Loop

		case Dot, QuestMark, Mul, Plus, LParen, RParen, Or, Caret:
			r = Unescape(r)
			fallthrough
		default:
			cs.Chars = append(cs.Chars, Char(r))
		}
	}
	if len(cs.Chars) == 0 {
		p.saveErr(errors.New("empty character set"))
		return nil
	}
	return &cs
}

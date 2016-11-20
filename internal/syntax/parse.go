package syntax

import "errors"

type parser struct {
	lex *lexer
	err error
}

func Compile(str string) (*Group, error) {
	l := newLexer(str)
	p := &parser{lex: l}
	gr := p.parseGroup()
	if p.err != nil {
		return nil, p.err
	}
	return gr, nil
}

func (p *parser) next() rune {
	if p.err != nil {
		return EOF
	}
	r := p.lex.Next()
	if r == Error {
		p.err = p.lex.Err()
		return EOF
	}
	return r
}

func (p *parser) saveErr(err error) {
	if p.err != nil {
		p.err = err
	}
}

func (p *parser) parseGroup() *Group {
	var gr Group
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
		default:
			gr.add(Char(r))
		}
	}
	return &gr
}

func (p *parser) parseCharSet() *CharSet {
	var cs CharSet
Loop:
	for {
		switch r := p.next(); r {
		case EOF:
			p.saveErr(errors.New("unexpected end of character set"))
			return nil
		case RBracket:
			break Loop

		case Dot, QuestMark, Mul, Plus:
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

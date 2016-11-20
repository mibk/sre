package syntax

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
				p.err = err
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
		default:
			gr.add(Char(r))
		}
	}
	return &gr
}

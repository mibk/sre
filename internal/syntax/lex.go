package syntax

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

const (
	_ = utf8.MaxRune + iota
	EOF
	Error

	Dot
	QuestMark
	Mul
	Plus

	LBracket
	RBracket
)

var specials = map[rune]rune{
	Dot:       '.',
	QuestMark: '?',
	Mul:       '*',
	Plus:      '+',
}

func Unescape(r rune) rune {
	r2, ok := specials[r]
	if !ok {
		panic(fmt.Sprintf("unknown special character: %c", r))
	}
	return r2
}

type lexer struct {
	src string
	i   int
	err error
}

func newLexer(s string) *lexer { return &lexer{src: s} }

func (l *lexer) Err() error { return l.err }

func (l *lexer) Next() rune {
	if l.err != nil {
		return Error
	}
	switch r := l.next(); r {
	case '\\':
		switch r := l.next(); r {
		case EOF:
			l.err = errors.New("trailing backslash")
			return Error
		case '\\', '.', '?', '*', '+', '[', ']':
			return r
		default:
			l.err = fmt.Errorf(`invalid escape sequence: \%c`)
			return r
		}
	case '.':
		return Dot
	case '?':
		return QuestMark
	case '*':
		return Mul
	case '+':
		return Plus
	case '[':
		return LBracket
	case ']':
		return RBracket
	default:
		return r
	}
}

func (l *lexer) next() rune {
	if l.i >= len(l.src) {
		return EOF
	}
	r, size := utf8.DecodeRuneInString(l.src[l.i:])
	l.i += size
	if r == utf8.RuneError {
		l.err = errors.New("invalid rune")
		return Error
	}
	return r
}

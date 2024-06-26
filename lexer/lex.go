package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/tjgurwara99/ruby-lsp/token"
)

type lexer struct {
	start   int
	pos     int
	width   int
	line    int
	col     int
	input   string
	Tokens  chan token.Token
	prevCol []int
}

const eof = -1

func Lex(input string) *lexer {
	l := &lexer{
		input:  input,
		Tokens: make(chan token.Token),
		line:   1,
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.Tokens)
}

func (l *lexer) emit(t token.TokenType) {
	if t == token.EOF {
		l.Tokens <- token.Token{
			Type:    t,
			Literal: "",
			Line:    l.line,
			Pos:     l.pos,
			Col:     l.col + (l.pos - l.start),
		}
		return
	}
	l.Tokens <- token.Token{
		Type:    t,
		Literal: l.input[l.start:l.pos],
		Line:    l.line,
		Pos:     l.start,
		Col:     l.col - (l.pos - l.start),
	}
	l.start = l.pos
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		if l.col == 0 {
			l.col = l.prevCol[len(l.prevCol)-1]
			l.line -= 1
		}
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.prevCol = append(l.prevCol, l.col)
	l.col += l.width
	l.pos += l.width
	if r == '\n' {
		l.line++
		l.col = 0
	}
	return r
}

// peek is usually supposed to be used for double character operators.
// Current syntax of Monkey doesn't have these double character operators yet.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
	l.col = l.prevCol[len(l.prevCol)-1]
	l.prevCol = l.prevCol[:len(l.prevCol)-1]
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

func (l *lexer) ignore() {
	l.start = l.pos
}

type stateFn func(*lexer) stateFn

func (l *lexer) doubleCharOperator(secondChar rune, failToken token.TokenType, passToken token.TokenType) token.TokenType {
	ch := l.peek()
	if ch == secondChar {
		l.next()
		return passToken
	}
	return failToken
}

func lexText(l *lexer) stateFn {
	switch r := l.next(); {
	case r == '\n' || r == ' ' || r == '\t' || r == '\r':
		l.ignore()
	case r == '=':
		l.emit(l.doubleCharOperator('=', token.Assign, token.Eq))
	case r == '+':
		l.emit(token.Plus)
	case r == '-':
		l.emit(token.Minus)
	case r == '!':
		l.emit(l.doubleCharOperator('=', token.Bang, token.NotEq))
	case r == '*':
		l.emit(token.Asterisk)
	case r == '/':
		l.emit(token.Slash)
	case r == '<':
		l.emit(l.doubleCharOperator('=', token.LessThan, token.LessThanEq))
	case r == '>':
		l.emit(l.doubleCharOperator('=', token.GreaterThan, token.GreaterThanEq))
	case r == ',':
		l.emit(token.Comma)
	case r == ';':
		l.emit(token.SemiColon)
	case r == '(':
		l.emit(token.LeftParen)
	case r == ')':
		l.emit(token.RightParen)
	case r == '{':
		l.emit(token.LeftBrace)
	case r == '}':
		l.emit(token.RightBrace)
	case r == '|':
		l.emit(l.doubleCharOperator('|', token.Or, token.Illegal))
	case r == '&':
		l.emit(l.doubleCharOperator('&', token.And, token.Illegal))
	case isIdent(r):
		l.backup()
		return lexIdent
	case r == eof:
		l.pos += 1
		l.col += 1
		l.emit(token.EOF)
		return nil
	}
	return lexText
}

func isIdent(r rune) bool {
	return r == '_' || unicode.IsDigit(r) || unicode.IsLetter(r)
}

func lexIdent(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isIdent(r):
			// absorb
		default:
			l.backup()
			l.emit(token.TokenForIdent(l.input[l.start:l.pos]))
			break Loop
		}
	}
	return lexText
}

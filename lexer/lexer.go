package lexer

import "github.com/alvindcastro/carlang/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipIgnored()

	var tok token.Token

	switch l.ch {
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	case 'R':
		tok = newToken(token.INVOKE, l.ch)
		l.readChar()
	case 'D':
		tok = newToken(token.CASTD, l.ch)
		l.readChar()
	case 'F':
		tok = newToken(token.CASTF, l.ch)
		l.readChar()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
		l.readChar()
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
		l.readChar()
	default:
		if isChantChar(l.ch) {
			literal := l.readChant()
			tok = token.Token{Type: token.CHANT, Literal: literal}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
			l.readChar()
		}
	}

	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipIgnored() {
	for {
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			l.readChar()
		}

		if l.ch == '/' && l.peekChar() == '/' {
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		if l.ch == '#' {
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		return
	}
}

func (l *Lexer) readChant() string {
	position := l.position
	for isChantChar(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isChantChar(ch byte) bool {
	return ch == 'Q' || ch == 'W' || ch == 'E'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

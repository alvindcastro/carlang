package lexer

import (
	"testing"

	"github.com/alvindcastro/carlang/token"
)

func TestNextToken(t *testing.T) {
	input := `
		// helper
		# shell-style helper comment
		[ WWE R D ] QQ R
		QEW QQW QQ EQE
		?
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LBRACKET, "["},
		{token.CHANT, "WWE"},
		{token.INVOKE, "R"},
		{token.CASTD, "D"},
		{token.RBRACKET, "]"},
		{token.CHANT, "QQ"},
		{token.INVOKE, "R"},
		{token.CHANT, "QEW"},
		{token.CHANT, "QQW"},
		{token.CHANT, "QQ"},
		{token.CHANT, "EQE"},
		{token.ILLEGAL, "?"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] token type wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

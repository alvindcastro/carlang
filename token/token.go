package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	CHANT  TokenType = "CHANT" // one or more Q/W/E chars
	INVOKE TokenType = "R"
	CASTD  TokenType = "D"
	CASTF  TokenType = "F"

	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"
)

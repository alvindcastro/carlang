package parser

import (
	"fmt"

	"github.com/alvindcastro/carlang/ast"
	"github.com/alvindcastro/carlang/lexer"
	"github.com/alvindcastro/carlang/token"
)

type Parser struct {
	tokens []token.Token
	pos    int
	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{}
	for {
		tok := l.NextToken()
		p.tokens = append(p.tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = p.parseStatements(token.EOF)
	return program
}

func (p *Parser) parseStatements(stop token.TokenType) []ast.Statement {
	statements := []ast.Statement{}

	for !p.currentIs(token.EOF) && !p.currentIs(stop) {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		} else {
			// Advance to avoid getting stuck after an error.
			p.pos++
		}
	}

	return statements
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.current().Type {
	case token.LBRACKET:
		return p.parseBlockOrDefinition()
	case token.CHANT:
		return p.parseChantStatement()
	case token.CASTD:
		p.pos++
		return &ast.CastStatement{Slot: 'D'}
	case token.CASTF:
		p.pos++
		return &ast.CastStatement{Slot: 'F'}
	case token.ILLEGAL:
		p.errors = append(p.errors, fmt.Sprintf("illegal token %q", p.current().Literal))
		return nil
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token %q", p.current().Literal))
		return nil
	}
}

func (p *Parser) parseBlockOrDefinition() ast.Statement {
	p.expect(token.LBRACKET)
	p.pos++ // consume '['

	block := &ast.BlockLiteral{Statements: p.parseStatements(token.RBRACKET)}

	if !p.currentIs(token.RBRACKET) {
		p.errors = append(p.errors, "expected closing ']' for block")
		return nil
	}
	p.pos++ // consume ']'

	// Definition shape: [ ... ] QWE R
	if p.currentIs(token.CHANT) && p.peekIs(token.INVOKE) {
		name := p.current().Literal
		p.pos++ // consume name
		p.pos++ // consume R
		return &ast.DefineComboStatement{Name: name, Block: block}
	}

	return &ast.PushBlockStatement{Block: block}
}

func (p *Parser) parseChantStatement() ast.Statement {
	name := p.current().Literal
	p.pos++ // consume chant

	if p.currentIs(token.INVOKE) {
		p.pos++ // consume R
		return &ast.InvokeStatement{Recipe: name}
	}

	return &ast.CallComboStatement{Name: name}
}

func (p *Parser) current() token.Token {
	if p.pos >= len(p.tokens) {
		return token.Token{Type: token.EOF, Literal: ""}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peek() token.Token {
	if p.pos+1 >= len(p.tokens) {
		return token.Token{Type: token.EOF, Literal: ""}
	}
	return p.tokens[p.pos+1]
}

func (p *Parser) currentIs(t token.TokenType) bool {
	return p.current().Type == t
}

func (p *Parser) peekIs(t token.TokenType) bool {
	return p.peek().Type == t
}

func (p *Parser) expect(t token.TokenType) bool {
	if !p.currentIs(t) {
		p.errors = append(p.errors, fmt.Sprintf("expected token %q, got %q", t, p.current().Type))
		return false
	}
	return true
}

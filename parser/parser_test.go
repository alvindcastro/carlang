package parser

import (
	"testing"

	"github.com/alvindcastro/carlang/ast"
	"github.com/alvindcastro/carlang/lexer"
)

func TestParseDefinitionsAndCalls(t *testing.T) {
	input := `[ WWE R D ] QQ R QQQ R D QQ`
	p := New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("program should have 4 statements, got=%d", len(program.Statements))
	}

	def, ok := program.Statements[0].(*ast.DefineComboStatement)
	if !ok {
		t.Fatalf("stmt[0] not DefineComboStatement. got=%T", program.Statements[0])
	}
	if def.Name != "QQ" {
		t.Fatalf("definition name wrong. got=%q", def.Name)
	}

	if _, ok := program.Statements[1].(*ast.InvokeStatement); !ok {
		t.Fatalf("stmt[1] not InvokeStatement. got=%T", program.Statements[1])
	}
	if _, ok := program.Statements[2].(*ast.CastStatement); !ok {
		t.Fatalf("stmt[2] not CastStatement. got=%T", program.Statements[2])
	}
	if _, ok := program.Statements[3].(*ast.CallComboStatement); !ok {
		t.Fatalf("stmt[3] not CallComboStatement. got=%T", program.Statements[3])
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Fatalf("parser had errors: %v", errors)
}

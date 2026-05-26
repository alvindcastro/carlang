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

func TestParsePushBlocksAndCasts(t *testing.T) {
	input := `[ QQQ R D ] D F`
	p := New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program should have 3 statements, got=%d", len(program.Statements))
	}

	block, ok := program.Statements[0].(*ast.PushBlockStatement)
	if !ok {
		t.Fatalf("stmt[0] not PushBlockStatement. got=%T", program.Statements[0])
	}
	if len(block.Block.Statements) != 2 {
		t.Fatalf("block should have 2 statements, got=%d", len(block.Block.Statements))
	}

	castD, ok := program.Statements[1].(*ast.CastStatement)
	if !ok || castD.Slot != 'D' {
		t.Fatalf("stmt[1] = %#v, want D CastStatement", program.Statements[1])
	}
	castF, ok := program.Statements[2].(*ast.CastStatement)
	if !ok || castF.Slot != 'F' {
		t.Fatalf("stmt[2] = %#v, want F CastStatement", program.Statements[2])
	}
}

func TestParserReportsErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"missing block close", "[ QQQ R D"},
		{"illegal token", "?"},
		{"unexpected invoke", "R"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(lexer.New(tt.input))
			p.ParseProgram()
			if len(p.Errors()) == 0 {
				t.Fatal("expected parser errors, got none")
			}
		})
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Fatalf("parser had errors: %v", errors)
}

package compiler

import (
	"testing"

	"github.com/alvindcastro/carlang/code"
	"github.com/alvindcastro/carlang/lexer"
	"github.com/alvindcastro/carlang/parser"
)

func TestCompileSimpleProgram(t *testing.T) {
	input := `QQQ R D QQE R D WWE R D EEE R D`
	compiler := compileInput(t, input)

	expected := code.Instructions{}
	expected = append(expected, code.Make(code.OpPushOne)...)
	expected = append(expected, code.Make(code.OpDup)...)
	expected = append(expected, code.Make(code.OpAdd)...)
	expected = append(expected, code.Make(code.OpPrintChar)...)

	if got := compiler.Bytecode().Instructions; string(got) != string(expected) {
		t.Fatalf("wrong instructions. got=%q, want=%q", got.String(), expected.String())
	}
}

func compileInput(t *testing.T, input string) *Compiler {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	compiler := New()
	if err := compiler.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}
	return compiler
}

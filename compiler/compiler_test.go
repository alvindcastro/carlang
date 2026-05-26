package compiler

import (
	"errors"
	"strings"
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

func TestCompileExpandsCombos(t *testing.T) {
	compiler := compileInput(t, `[ QQQ R D QQQ R D WWE R D ] QQ R QQ`)

	expected := code.Instructions{}
	expected = append(expected, code.Make(code.OpPushOne)...)
	expected = append(expected, code.Make(code.OpPushOne)...)
	expected = append(expected, code.Make(code.OpAdd)...)

	if got := compiler.Bytecode().Instructions; string(got) != string(expected) {
		t.Fatalf("wrong instructions. got=%q, want=%q", got.String(), expected.String())
	}
}

func TestCompileUsesFSlot(t *testing.T) {
	compiler := compileInput(t, `QQQ R QQE R F D`)

	expected := code.Instructions{}
	expected = append(expected, code.Make(code.OpPushOne)...)
	expected = append(expected, code.Make(code.OpDup)...)

	if got := compiler.Bytecode().Instructions; string(got) != string(expected) {
		t.Fatalf("wrong instructions. got=%q, want=%q", got.String(), expected.String())
	}
}

func TestCompileReportsUnsupportedPrograms(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError string
	}{
		{"undefined combo", `QQ`, "undefined combo"},
		{"recursive combo", `[ QQ ] QQ R QQ`, "recursive combo"},
		{"invalid recipe", `QQ R`, "must contain exactly three"},
		{"empty slot", `D`, "empty D slot"},
		{"anonymous block", `[ QQQ R D ]`, "anonymous block"},
		{"unsupported spell", `QEE R D`, "does not yet support spell"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := compileInputWithError(t, tt.input)
			if err == nil {
				t.Fatal("expected compile error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tt.wantError)
			}
		})
	}
}

func compileInput(t *testing.T, input string) *Compiler {
	t.Helper()
	compiler, err := parseAndCompile(input)
	if err != nil {
		t.Fatalf("compiler error: %v", err)
	}
	return compiler
}

func compileInputWithError(t *testing.T, input string) error {
	t.Helper()
	_, err := parseAndCompile(input)
	return err
}

func parseAndCompile(input string) (*Compiler, error) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		return nil, errors.New(strings.Join(p.Errors(), "; "))
	}
	compiler := New()
	if err := compiler.Compile(program); err != nil {
		return nil, err
	}
	return compiler, nil
}

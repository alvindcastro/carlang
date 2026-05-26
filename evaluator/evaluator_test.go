package evaluator

import (
	"os"
	"testing"

	"github.com/alvindcastro/carlang/lexer"
	"github.com/alvindcastro/carlang/parser"
)

func TestEvalAddsTwo(t *testing.T) {
	vm := evalInput(t, `QQQ R D QQQ R D WWE R D`)
	if len(vm.Mana) != 1 || vm.Mana[0] != 2 {
		t.Fatalf("expected mana stack [2], got=%v", vm.Mana)
	}
}

func TestEvalHelloWorld(t *testing.T) {
	input, err := os.ReadFile("../examples/hello_world.carl")
	if err != nil {
		t.Fatal(err)
	}
	vm := evalInput(t, string(input))
	if got, want := vm.Output.String(), "Hello, World!\n"; got != want {
		t.Fatalf("wrong output. got=%q, want=%q", got, want)
	}
}

func evalInput(t *testing.T, input string) *VM {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	vm := NewVM()
	if err := Eval(program, vm); err != nil {
		t.Fatalf("eval error: %v", err)
	}
	return vm
}

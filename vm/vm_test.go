package vm

import (
	"os"
	"testing"

	"github.com/alvindcastro/carlang/compiler"
	"github.com/alvindcastro/carlang/lexer"
	"github.com/alvindcastro/carlang/parser"
)

func TestVMHelloWorld(t *testing.T) {
	input, err := os.ReadFile("../examples/hello_world.carl")
	if err != nil {
		t.Fatal(err)
	}

	l := lexer.New(string(input))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	machine := New(comp.Bytecode().Instructions)
	if err := machine.Run(); err != nil {
		t.Fatalf("vm error: %v", err)
	}

	if got, want := machine.Output.String(), "Hello, World!\n"; got != want {
		t.Fatalf("wrong output. got=%q, want=%q", got, want)
	}
}

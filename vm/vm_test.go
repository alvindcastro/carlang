package vm

import (
	"os"
	"strings"
	"testing"

	"github.com/alvindcastro/carlang/code"
	"github.com/alvindcastro/carlang/compiler"
	"github.com/alvindcastro/carlang/lexer"
	"github.com/alvindcastro/carlang/parser"
)

func TestVMRunsCoreOpcodes(t *testing.T) {
	tests := []struct {
		name       string
		ops        []code.Opcode
		wantStack  []int
		wantOutput string
	}{
		{
			name:      "drop",
			ops:       []code.Opcode{code.OpPushOne, code.OpDrop},
			wantStack: []int{},
		},
		{
			name:      "dup and add",
			ops:       []code.Opcode{code.OpPushOne, code.OpDup, code.OpAdd},
			wantStack: []int{2},
		},
		{
			name:      "swap",
			ops:       []code.Opcode{code.OpPushOne, code.OpPushOne, code.OpPushOne, code.OpAdd, code.OpSwap},
			wantStack: []int{2, 1},
		},
		{
			name:      "sub",
			ops:       []code.Opcode{code.OpPushOne, code.OpPushOne, code.OpPushOne, code.OpAdd, code.OpSub},
			wantStack: []int{-1},
		},
		{
			name:      "mul",
			ops:       []code.Opcode{code.OpPushOne, code.OpPushOne, code.OpAdd, code.OpDup, code.OpMul},
			wantStack: []int{4},
		},
		{
			name:       "print",
			ops:        []code.Opcode{code.OpPushOne, code.OpPushOne, code.OpAdd, code.OpPrintChar},
			wantStack:  []int{},
			wantOutput: "\x02",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			machine := New(instructions(tt.ops...))
			if err := machine.Run(); err != nil {
				t.Fatalf("Run returned error: %v", err)
			}
			assertStack(t, machine.Stack(), tt.wantStack)
			if got := machine.Output.String(); got != tt.wantOutput {
				t.Fatalf("output = %q, want %q", got, tt.wantOutput)
			}
		})
	}
}

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

func TestVMStackReturnsCopy(t *testing.T) {
	machine := New(instructions(code.OpPushOne))
	if err := machine.Run(); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	got := machine.Stack()
	got[0] = 99
	assertStack(t, machine.Stack(), []int{1})
}

func TestVMReportsRuntimeErrors(t *testing.T) {
	tests := []struct {
		name      string
		ins       code.Instructions
		wantError string
	}{
		{"unknown opcode", code.Instructions{255}, "unknown opcode"},
		{"drop underflow", instructions(code.OpDrop), "stack underflow"},
		{"dup underflow", instructions(code.OpDup), "stack underflow"},
		{"swap underflow", instructions(code.OpPushOne, code.OpSwap), "requires two"},
		{"add underflow", instructions(code.OpPushOne, code.OpAdd), "stack underflow"},
		{"print invalid rune", instructions(code.OpPushOne, code.OpPushOne, code.OpPushOne, code.OpAdd, code.OpSub, code.OpPrintChar), "invalid rune"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.ins).Run()
			if err == nil {
				t.Fatal("expected VM error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tt.wantError)
			}
		})
	}
}

func instructions(ops ...code.Opcode) code.Instructions {
	out := code.Instructions{}
	for _, op := range ops {
		out = append(out, code.Make(op)...)
	}
	return out
}

func assertStack(t *testing.T, got []int, want []int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("stack = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("stack = %v, want %v", got, want)
		}
	}
}

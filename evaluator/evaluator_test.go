package evaluator

import (
	"errors"
	"os"
	"strings"
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

func TestEvalCoreSpells(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantMana   []int
		wantBlocks int
		wantOutput string
	}{
		{"cold snap pushes one", `QQQ R D`, []int{1}, 0, ""},
		{"ghost walk drops top mana", `QQQ R D QQW R D`, []int{}, 0, ""},
		{"ice wall duplicates top mana", `QQQ R D QQE R D`, []int{1, 1}, 0, ""},
		{"tornado swaps top two mana values", `QQQ R D QQQ R D QQQ R D WWE R D QWW R D`, []int{2, 1}, 0, ""},
		{"emp subtracts top from next", `QQQ R D QQQ R D QQQ R D WWE R D WWW R D`, []int{-1}, 0, ""},
		{"alacrity adds", `QQQ R D QQQ R D WWE R D`, []int{2}, 0, ""},
		{"chaos meteor multiplies", `QQQ R D QQQ R D WWE R D QQQ R D QQQ R D WWE R D WEE R D`, []int{4}, 0, ""},
		{"forge spirit duplicates top block", `[ [ QQQ R D ] ] QQ R QQ QEE R D`, []int{}, 2, ""},
		{"sun strike prints and pops", `QQQ R D QQQ R D WWE R D EEE R D`, []int{}, 0, "\x02"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := evalInput(t, tt.input)
			assertMana(t, vm.Mana, tt.wantMana)
			if len(vm.Grimoire) != tt.wantBlocks {
				t.Fatalf("grimoire stack length = %d, want %d", len(vm.Grimoire), tt.wantBlocks)
			}
			if got := vm.Output.String(); got != tt.wantOutput {
				t.Fatalf("output = %q, want %q", got, tt.wantOutput)
			}
		})
	}
}

func TestEvalSpellbookFSlot(t *testing.T) {
	vm := evalInput(t, `QQQ R QQE R F D`)
	assertMana(t, vm.Mana, []int{1, 1})
}

func TestEvalDeafeningBlastChoosesBranch(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantMana []int
	}{
		{
			name:     "true branch",
			input:    `[ [ QQQ R D ] [ QQQ R D QQQ R D WWE R D ] ] QQ R QQ QQQ R D QWE R D`,
			wantMana: []int{1},
		},
		{
			name:     "false branch",
			input:    `[ [ QQQ R D ] [ QQQ R D QQQ R D WWE R D ] ] QQ R QQ QQQ R D QQQ R D WWW R D QWE R D`,
			wantMana: []int{2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := evalInput(t, tt.input)
			assertMana(t, vm.Mana, tt.wantMana)
		})
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

func TestEvalReportsRuntimeErrors(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError string
	}{
		{"undefined combo", `QQ`, "undefined combo"},
		{"invalid recipe", `QQ R`, "must contain exactly three"},
		{"empty cast slot", `D`, "empty D slot"},
		{"mana underflow", `QQW R D`, "Mana Stack underflow"},
		{"forge spirit grimoire underflow", `QEE R D`, "requires one block"},
		{"deafening blast grimoire underflow", `QQQ R D QWE R D`, "Grimoire Stack underflow"},
		{"sun strike invalid rune", `QQQ R D QQQ R D QQQ R D WWE R D WWW R D EEE R D`, "invalid rune"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := evalInputWithError(t, tt.input)
			if err == nil {
				t.Fatal("expected eval error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tt.wantError)
			}
		})
	}
}

func evalInput(t *testing.T, input string) *VM {
	t.Helper()
	vm, err := parseAndEval(input)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	return vm
}

func evalInputWithError(t *testing.T, input string) error {
	t.Helper()
	_, err := parseAndEval(input)
	return err
}

func parseAndEval(input string) (*VM, error) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		return nil, errors.New(strings.Join(p.Errors(), "; "))
	}
	vm := NewVM()
	if err := Eval(program, vm); err != nil {
		return nil, err
	}
	return vm, nil
}

func assertMana(t *testing.T, got []int, want []int) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("mana stack = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("mana stack = %v, want %v", got, want)
		}
	}
}

package ast

import "testing"

func TestProgramString(t *testing.T) {
	program := &Program{Statements: []Statement{
		&DefineComboStatement{
			Name: "QQ",
			Block: &BlockLiteral{Statements: []Statement{
				&InvokeStatement{Recipe: "WWE"},
				&CastStatement{Slot: 'D'},
			}},
		},
		&PushBlockStatement{Block: &BlockLiteral{Statements: []Statement{
			&InvokeStatement{Recipe: "QQQ"},
			&CastStatement{Slot: 'F'},
		}}},
		&CallComboStatement{Name: "QQ"},
	}}

	want := "[WWE R D] QQ R\n[QQQ R F]\nQQ\n"
	if got := program.String(); got != want {
		t.Fatalf("program.String() = %q, want %q", got, want)
	}
}

package code

import "testing"

func TestInstructionsString(t *testing.T) {
	instructions := Instructions{}
	instructions = append(instructions, Make(OpPushOne)...)
	instructions = append(instructions, Make(OpDup)...)
	instructions = append(instructions, Make(OpAdd)...)
	instructions = append(instructions, Make(OpPrintChar)...)

	expected := "0000 OpPushOne\n0001 OpDup\n0002 OpAdd\n0003 OpPrintChar\n"
	if instructions.String() != expected {
		t.Fatalf("instructions wrongly formatted. got=%q", instructions.String())
	}
}

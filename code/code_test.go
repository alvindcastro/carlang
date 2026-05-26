package code

import (
	"os"
	"path/filepath"
	"testing"
)

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

func TestInstructionsStringFormatsUnknownOpcode(t *testing.T) {
	instructions := Instructions{255, byte(OpPushOne)}

	expected := "ERROR: opcode 255 undefined\n0001 OpPushOne\n"
	if got := instructions.String(); got != expected {
		t.Fatalf("instructions wrongly formatted. got=%q, want=%q", got, expected)
	}
}

func TestMakeReturnsEmptyInstructionForUnknownOpcode(t *testing.T) {
	if got := Make(Opcode(255)); len(got) != 0 {
		t.Fatalf("Make returned %v, want empty instruction", got)
	}
}

func TestReadOperandsWithTwoByteOperand(t *testing.T) {
	def := &Definition{Name: "OpTest", OperandWidths: []int{2}}
	operands, read := ReadOperands(def, Instructions{0x01, 0x00})

	if read != 2 {
		t.Fatalf("read = %d, want 2", read)
	}
	if len(operands) != 1 || operands[0] != 256 {
		t.Fatalf("operands = %v, want [256]", operands)
	}
}

func TestBytecodeFileRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "program.cbc")
	instructions := Instructions{}
	instructions = append(instructions, Make(OpPushOne)...)
	instructions = append(instructions, Make(OpPrintChar)...)

	if err := WriteFile(path, instructions); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	got, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(got) != string(instructions) {
		t.Fatalf("ReadFile returned %v, want %v", got, instructions)
	}
}

func TestReadFileRejectsInvalidMagic(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.cbc")
	if err := os.WriteFile(path, []byte("not bytecode"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := ReadFile(path); err == nil {
		t.Fatal("ReadFile returned nil error for invalid bytecode")
	}
}

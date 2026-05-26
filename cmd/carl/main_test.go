package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alvindcastro/carlang/code"
)

func TestParseFileReportsParserErrors(t *testing.T) {
	path := writeTempFile(t, "bad.carl", "?")

	if _, err := parseFile(path); err == nil {
		t.Fatal("parseFile returned nil error for invalid source")
	}
}

func TestRunCommandExecutesProgram(t *testing.T) {
	path := writeTempFile(t, "print.carl", printTwoProgram)

	got := captureStdout(t, func() {
		if err := runCommand([]string{path}); err != nil {
			t.Fatalf("runCommand returned error: %v", err)
		}
	})

	if got != "\x02" {
		t.Fatalf("stdout = %q, want %q", got, "\x02")
	}
}

func TestCompileAndVMCommandsRoundTrip(t *testing.T) {
	sourcePath := writeTempFile(t, "print.carl", printTwoProgram)
	bytecodePath := filepath.Join(t.TempDir(), "print.cbc")

	if err := compileCommand([]string{sourcePath, "-o", bytecodePath}); err != nil {
		t.Fatalf("compileCommand returned error: %v", err)
	}

	instructions, err := code.ReadFile(bytecodePath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if len(instructions) == 0 {
		t.Fatal("compiled bytecode was empty")
	}

	got := captureStdout(t, func() {
		if err := vmCommand([]string{bytecodePath}); err != nil {
			t.Fatalf("vmCommand returned error: %v", err)
		}
	})
	if got != "\x02" {
		t.Fatalf("stdout = %q, want %q", got, "\x02")
	}
}

func TestDisCommandDisassemblesSourceAndBytecode(t *testing.T) {
	sourcePath := writeTempFile(t, "print.carl", printTwoProgram)
	bytecodePath := filepath.Join(t.TempDir(), "print.cbc")
	if err := compileCommand([]string{sourcePath, "-o", bytecodePath}); err != nil {
		t.Fatalf("compileCommand returned error: %v", err)
	}

	for _, path := range []string{sourcePath, bytecodePath} {
		t.Run(filepath.Ext(path), func(t *testing.T) {
			got := captureStdout(t, func() {
				if err := disCommand([]string{path}); err != nil {
					t.Fatalf("disCommand returned error: %v", err)
				}
			})
			if !strings.Contains(got, "OpPrintChar") {
				t.Fatalf("disassembly = %q, want OpPrintChar", got)
			}
		})
	}
}

func TestCommandsReportUsageErrors(t *testing.T) {
	tests := []struct {
		name string
		run  func() error
	}{
		{"run missing path", func() error { return runCommand(nil) }},
		{"run unknown flag", func() error { return runCommand([]string{"--bad", "file.carl"}) }},
		{"compile missing path", func() error { return compileCommand(nil) }},
		{"compile missing output value", func() error { return compileCommand([]string{"source.carl", "-o"}) }},
		{"vm wrong arity", func() error { return vmCommand(nil) }},
		{"dis wrong arity", func() error { return disCommand(nil) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.run(); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func writeTempFile(t *testing.T, name string, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = old
	}()

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

const printTwoProgram = `QQQ R D QQQ R D WWE R D EEE R D`

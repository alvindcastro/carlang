package repl

import (
	"strings"
	"testing"
)

func TestStartEvaluatesLinesWithPersistentState(t *testing.T) {
	var out strings.Builder
	Start(strings.NewReader("QQQ R D\nQQQ R D WWE R D\n"), &out)

	got := out.String()
	for _, want := range []string{
		"carl> Mana: [1]\n",
		"carl> Mana: [2]\n",
		"carl> ",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("REPL output = %q, want substring %q", got, want)
		}
	}
}

func TestStartReportsParserAndEvalErrors(t *testing.T) {
	var out strings.Builder
	Start(strings.NewReader("?\nD\n"), &out)

	got := out.String()
	for _, want := range []string{
		"parser error:",
		"error: cannot cast empty D slot",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("REPL output = %q, want substring %q", got, want)
		}
	}
}

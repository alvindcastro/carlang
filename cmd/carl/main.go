package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alvindcastro/carlang/ast"
	"github.com/alvindcastro/carlang/code"
	"github.com/alvindcastro/carlang/compiler"
	"github.com/alvindcastro/carlang/evaluator"
	"github.com/alvindcastro/carlang/lexer"
	"github.com/alvindcastro/carlang/parser"
	"github.com/alvindcastro/carlang/repl"
	bytecodevm "github.com/alvindcastro/carlang/vm"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		if err := runCommand(os.Args[2:]); err != nil {
			fatal(err)
		}
	case "compile":
		if err := compileCommand(os.Args[2:]); err != nil {
			fatal(err)
		}
	case "vm":
		if err := vmCommand(os.Args[2:]); err != nil {
			fatal(err)
		}
	case "dis":
		if err := disCommand(os.Args[2:]); err != nil {
			fatal(err)
		}
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	default:
		usage()
		os.Exit(1)
	}
}

func runCommand(args []string) error {
	var debug bool
	var inputPath string

	for _, arg := range args {
		switch arg {
		case "--debug", "-debug":
			debug = true
		default:
			if strings.HasPrefix(arg, "-") {
				return fmt.Errorf("unknown run flag %q", arg)
			}
			if inputPath != "" {
				return fmt.Errorf("usage: carl run [--debug] <file.carl>")
			}
			inputPath = arg
		}
	}

	if inputPath == "" {
		return fmt.Errorf("usage: carl run [--debug] <file.carl>")
	}

	program, err := parseFile(inputPath)
	if err != nil {
		return err
	}

	machine := evaluator.NewVM()
	if err := evaluator.Eval(program, machine); err != nil {
		return err
	}

	fmt.Print(machine.Output.String())

	if debug {
		fmt.Fprintf(os.Stderr, "\nMana Stack: %v\n", machine.Mana)
		fmt.Fprintf(os.Stderr, "Grimoire Stack: %d block(s)\n", len(machine.Grimoire))
		fmt.Fprintf(os.Stderr, "D: %s\n", machine.Spellbook[0])
		fmt.Fprintf(os.Stderr, "F: %s\n", machine.Spellbook[1])
		fmt.Fprintf(os.Stderr, "Defined Combos: %d\n", len(machine.Tome))
	}

	return nil
}

func compileCommand(args []string) error {
	var outPath string
	var disassemble bool
	var inputPath string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-o", "--output":
			if i+1 >= len(args) {
				return fmt.Errorf("missing value for %s", arg)
			}
			i++
			outPath = args[i]
		case "--disassemble", "-disassemble":
			disassemble = true
		default:
			if strings.HasPrefix(arg, "-") {
				return fmt.Errorf("unknown compile flag %q", arg)
			}
			if inputPath != "" {
				return fmt.Errorf("usage: carl compile <file.carl> [-o file.cbc] [--disassemble]")
			}
			inputPath = arg
		}
	}

	if inputPath == "" {
		return fmt.Errorf("usage: carl compile <file.carl> [-o file.cbc] [--disassemble]")
	}

	program, err := parseFile(inputPath)
	if err != nil {
		return err
	}

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		return err
	}

	if outPath == "" {
		outPath = strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".cbc"
	}

	if dir := filepath.Dir(outPath); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	instructions := comp.Bytecode().Instructions
	if err := code.WriteFile(outPath, instructions); err != nil {
		return err
	}

	if disassemble {
		fmt.Print(instructions.String())
	}

	return nil
}

func vmCommand(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: carl vm <file.cbc>")
	}
	instructions, err := code.ReadFile(args[0])
	if err != nil {
		return err
	}
	machine := bytecodevm.New(instructions)
	if err := machine.Run(); err != nil {
		return err
	}
	fmt.Print(machine.Output.String())
	return nil
}

func disCommand(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: carl dis <file.carl|file.cbc>")
	}

	var instructions code.Instructions
	if strings.HasSuffix(args[0], ".cbc") {
		var err error
		instructions, err = code.ReadFile(args[0])
		if err != nil {
			return err
		}
	} else {
		program, err := parseFile(args[0])
		if err != nil {
			return err
		}
		comp := compiler.New()
		if err := comp.Compile(program); err != nil {
			return err
		}
		instructions = comp.Bytecode().Instructions
	}

	fmt.Print(instructions.String())
	return nil
}

func parseFile(path string) (*ast.Program, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	l := lexer.New(string(input))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parser errors: %s", strings.Join(p.Errors(), "; "))
	}
	return program, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "Carlang")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  carl run <file.carl>")
	fmt.Fprintln(os.Stderr, "  carl run --debug <file.carl>")
	fmt.Fprintln(os.Stderr, "  carl compile <file.carl> [-o file.cbc] [--disassemble]")
	fmt.Fprintln(os.Stderr, "  carl vm <file.cbc>")
	fmt.Fprintln(os.Stderr, "  carl dis <file.carl|file.cbc>")
	fmt.Fprintln(os.Stderr, "  carl repl")
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	os.Exit(1)
}

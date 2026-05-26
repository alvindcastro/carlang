package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/alvindcastro/carlang/evaluator"
	"github.com/alvindcastro/carlang/lexer"
	"github.com/alvindcastro/carlang/parser"
)

const Prompt = "carl> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	vm := evaluator.NewVM()

	for {
		fmt.Fprint(out, Prompt)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				fmt.Fprintf(out, "parser error: %s\n", msg)
			}
			continue
		}

		before := vm.Output.Len()
		if err := evaluator.Eval(program, vm); err != nil {
			fmt.Fprintf(out, "error: %s\n", err)
			continue
		}

		if vm.Output.Len() > before {
			output := vm.Output.String()[before:]
			fmt.Fprint(out, output)
			if output[len(output)-1] != '\n' {
				fmt.Fprintln(out)
			}
		}
		fmt.Fprintf(out, "Mana: %v\n", vm.Mana)
	}
}

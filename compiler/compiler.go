package compiler

import (
	"fmt"

	"github.com/alvindcastro/carlang/ast"
	"github.com/alvindcastro/carlang/code"
	"github.com/alvindcastro/carlang/spell"
)

type Bytecode struct {
	Instructions code.Instructions
}

type Compiler struct {
	instructions code.Instructions
	tome         map[string]*ast.BlockLiteral
	spellbook    [2]spell.Spell
	activeMacros map[string]bool
}

func New() *Compiler {
	return &Compiler{
		tome:         make(map[string]*ast.BlockLiteral),
		activeMacros: make(map[string]bool),
	}
}

func (c *Compiler) Compile(program *ast.Program) error {
	for _, stmt := range program.Statements {
		if err := c.compileStatement(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{Instructions: c.instructions}
}

func (c *Compiler) compileBlock(block *ast.BlockLiteral) error {
	for _, stmt := range block.Statements {
		if err := c.compileStatement(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) compileStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.DefineComboStatement:
		c.tome[s.Name] = s.Block
		return nil
	case *ast.CallComboStatement:
		block, ok := c.tome[s.Name]
		if !ok {
			return fmt.Errorf("undefined combo %q", s.Name)
		}
		if c.activeMacros[s.Name] {
			return fmt.Errorf("recursive combo %q cannot be compiled", s.Name)
		}
		c.activeMacros[s.Name] = true
		err := c.compileBlock(block)
		delete(c.activeMacros, s.Name)
		return err
	case *ast.InvokeStatement:
		sp, err := spell.FromRecipe(s.Recipe)
		if err != nil {
			return err
		}
		c.spellbook[1] = c.spellbook[0]
		c.spellbook[0] = sp
		return nil
	case *ast.CastStatement:
		var sp spell.Spell
		if s.Slot == 'D' {
			sp = c.spellbook[0]
		} else {
			sp = c.spellbook[1]
		}
		if sp == spell.Empty {
			return fmt.Errorf("cannot compile cast of empty %c slot", s.Slot)
		}
		op, err := opcodeForSpell(sp)
		if err != nil {
			return err
		}
		c.emit(op)
		return nil
	case *ast.PushBlockStatement:
		return fmt.Errorf("compiler does not yet support anonymous block values")
	default:
		return fmt.Errorf("unsupported statement %T", stmt)
	}
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func opcodeForSpell(sp spell.Spell) (code.Opcode, error) {
	switch sp {
	case spell.ColdSnap:
		return code.OpPushOne, nil
	case spell.GhostWalk:
		return code.OpDrop, nil
	case spell.IceWall:
		return code.OpDup, nil
	case spell.Tornado:
		return code.OpSwap, nil
	case spell.EMP:
		return code.OpSub, nil
	case spell.Alacrity:
		return code.OpAdd, nil
	case spell.ChaosMeteor:
		return code.OpMul, nil
	case spell.SunStrike:
		return code.OpPrintChar, nil
	default:
		return 0, fmt.Errorf("compiler does not yet support spell %s", sp)
	}
}

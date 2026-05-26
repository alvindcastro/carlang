package evaluator

import (
	"fmt"
	"strings"

	"github.com/alvindcastro/carlang/ast"
	"github.com/alvindcastro/carlang/spell"
)

type VM struct {
	Mana      []int
	Grimoire  []*ast.BlockLiteral
	Spellbook [2]spell.Spell // D = [0], F = [1]
	Tome      map[string]*ast.BlockLiteral
	Output    strings.Builder
}

func NewVM() *VM {
	return &VM{Tome: make(map[string]*ast.BlockLiteral)}
}

func Eval(program *ast.Program, vm *VM) error {
	for _, stmt := range program.Statements {
		if err := evalStatement(stmt, vm); err != nil {
			return err
		}
	}
	return nil
}

func evalBlock(block *ast.BlockLiteral, vm *VM) error {
	for _, stmt := range block.Statements {
		if err := evalStatement(stmt, vm); err != nil {
			return err
		}
	}
	return nil
}

func evalStatement(stmt ast.Statement, vm *VM) error {
	switch s := stmt.(type) {
	case *ast.DefineComboStatement:
		vm.Tome[s.Name] = s.Block
		return nil
	case *ast.PushBlockStatement:
		vm.Grimoire = append(vm.Grimoire, s.Block)
		return nil
	case *ast.CallComboStatement:
		block, ok := vm.Tome[s.Name]
		if !ok {
			return fmt.Errorf("undefined combo %q", s.Name)
		}
		return evalBlock(block, vm)
	case *ast.InvokeStatement:
		sp, err := spell.FromRecipe(s.Recipe)
		if err != nil {
			return err
		}
		vm.Spellbook[1] = vm.Spellbook[0]
		vm.Spellbook[0] = sp
		return nil
	case *ast.CastStatement:
		var sp spell.Spell
		if s.Slot == 'D' {
			sp = vm.Spellbook[0]
		} else {
			sp = vm.Spellbook[1]
		}
		if sp == spell.Empty {
			return fmt.Errorf("cannot cast empty %c slot", s.Slot)
		}
		return cast(sp, vm)
	default:
		return fmt.Errorf("unsupported statement %T", stmt)
	}
}

func cast(sp spell.Spell, vm *VM) error {
	switch sp {
	case spell.ColdSnap:
		vm.Mana = append(vm.Mana, 1)
	case spell.GhostWalk:
		_, err := popMana(vm)
		return err
	case spell.IceWall:
		v, err := peekMana(vm)
		if err != nil {
			return err
		}
		vm.Mana = append(vm.Mana, v)
	case spell.Tornado:
		if len(vm.Mana) < 2 {
			return fmt.Errorf("Tornado requires two mana values")
		}
		top := len(vm.Mana) - 1
		vm.Mana[top], vm.Mana[top-1] = vm.Mana[top-1], vm.Mana[top]
	case spell.EMP:
		b, err := popMana(vm)
		if err != nil {
			return err
		}
		a, err := popMana(vm)
		if err != nil {
			return err
		}
		vm.Mana = append(vm.Mana, a-b)
	case spell.Alacrity:
		b, err := popMana(vm)
		if err != nil {
			return err
		}
		a, err := popMana(vm)
		if err != nil {
			return err
		}
		vm.Mana = append(vm.Mana, a+b)
	case spell.ChaosMeteor:
		b, err := popMana(vm)
		if err != nil {
			return err
		}
		a, err := popMana(vm)
		if err != nil {
			return err
		}
		vm.Mana = append(vm.Mana, a*b)
	case spell.ForgeSpirit:
		if len(vm.Grimoire) < 1 {
			return fmt.Errorf("Forge Spirit requires one block on the Grimoire Stack")
		}
		vm.Grimoire = append(vm.Grimoire, vm.Grimoire[len(vm.Grimoire)-1])
	case spell.SunStrike:
		v, err := popMana(vm)
		if err != nil {
			return err
		}
		if v < 0 || v > 1114111 {
			return fmt.Errorf("Sun Strike cannot print invalid rune value %d", v)
		}
		vm.Output.WriteRune(rune(v))
	case spell.DeafeningBlast:
		condition, err := popMana(vm)
		if err != nil {
			return err
		}
		falseBlock, err := popBlock(vm)
		if err != nil {
			return err
		}
		trueBlock, err := popBlock(vm)
		if err != nil {
			return err
		}
		if condition != 0 {
			return evalBlock(trueBlock, vm)
		}
		return evalBlock(falseBlock, vm)
	default:
		return fmt.Errorf("unsupported spell %s", sp)
	}
	return nil
}

func popMana(vm *VM) (int, error) {
	if len(vm.Mana) == 0 {
		return 0, fmt.Errorf("Mana Stack underflow")
	}
	v := vm.Mana[len(vm.Mana)-1]
	vm.Mana = vm.Mana[:len(vm.Mana)-1]
	return v, nil
}

func peekMana(vm *VM) (int, error) {
	if len(vm.Mana) == 0 {
		return 0, fmt.Errorf("Mana Stack underflow")
	}
	return vm.Mana[len(vm.Mana)-1], nil
}

func popBlock(vm *VM) (*ast.BlockLiteral, error) {
	if len(vm.Grimoire) == 0 {
		return nil, fmt.Errorf("Grimoire Stack underflow")
	}
	b := vm.Grimoire[len(vm.Grimoire)-1]
	vm.Grimoire = vm.Grimoire[:len(vm.Grimoire)-1]
	return b, nil
}

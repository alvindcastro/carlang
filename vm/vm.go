package vm

import (
	"fmt"
	"strings"

	"github.com/alvindcastro/carlang/code"
)

type VM struct {
	instructions code.Instructions
	stack        []int
	Output       strings.Builder
}

func New(instructions code.Instructions) *VM {
	return &VM{instructions: instructions}
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpPushOne:
			vm.push(1)
		case code.OpDrop:
			if _, err := vm.pop(); err != nil {
				return err
			}
		case code.OpDup:
			v, err := vm.peek()
			if err != nil {
				return err
			}
			vm.push(v)
		case code.OpSwap:
			if len(vm.stack) < 2 {
				return fmt.Errorf("OpSwap requires two stack values")
			}
			top := len(vm.stack) - 1
			vm.stack[top], vm.stack[top-1] = vm.stack[top-1], vm.stack[top]
		case code.OpSub:
			b, err := vm.pop()
			if err != nil {
				return err
			}
			a, err := vm.pop()
			if err != nil {
				return err
			}
			vm.push(a - b)
		case code.OpAdd:
			b, err := vm.pop()
			if err != nil {
				return err
			}
			a, err := vm.pop()
			if err != nil {
				return err
			}
			vm.push(a + b)
		case code.OpMul:
			b, err := vm.pop()
			if err != nil {
				return err
			}
			a, err := vm.pop()
			if err != nil {
				return err
			}
			vm.push(a * b)
		case code.OpPrintChar:
			v, err := vm.pop()
			if err != nil {
				return err
			}
			if v < 0 || v > 1114111 {
				return fmt.Errorf("OpPrintChar cannot print invalid rune value %d", v)
			}
			vm.Output.WriteRune(rune(v))
		default:
			return fmt.Errorf("unknown opcode %d", op)
		}
	}
	return nil
}

func (vm *VM) Stack() []int {
	out := make([]int, len(vm.stack))
	copy(out, vm.stack)
	return out
}

func (vm *VM) push(v int) {
	vm.stack = append(vm.stack, v)
}

func (vm *VM) pop() (int, error) {
	if len(vm.stack) == 0 {
		return 0, fmt.Errorf("VM stack underflow")
	}
	v := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return v, nil
}

func (vm *VM) peek() (int, error) {
	if len(vm.stack) == 0 {
		return 0, fmt.Errorf("VM stack underflow")
	}
	return vm.stack[len(vm.stack)-1], nil
}

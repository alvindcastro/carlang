package ast

import "strings"

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	return out.String()
}

type BlockLiteral struct {
	Statements []Statement
}

func (b *BlockLiteral) String() string {
	var out strings.Builder
	out.WriteString("[")
	for i, s := range b.Statements {
		if i > 0 {
			out.WriteString(" ")
		}
		out.WriteString(s.String())
	}
	out.WriteString("]")
	return out.String()
}

type PushBlockStatement struct {
	Block *BlockLiteral
}

func (*PushBlockStatement) statementNode()   {}
func (s *PushBlockStatement) String() string { return s.Block.String() }

type DefineComboStatement struct {
	Name  string
	Block *BlockLiteral
}

func (*DefineComboStatement) statementNode()   {}
func (s *DefineComboStatement) String() string { return s.Block.String() + " " + s.Name + " R" }

type InvokeStatement struct {
	Recipe string
}

func (*InvokeStatement) statementNode()   {}
func (s *InvokeStatement) String() string { return s.Recipe + " R" }

type CastStatement struct {
	Slot rune // 'D' or 'F'
}

func (*CastStatement) statementNode()   {}
func (s *CastStatement) String() string { return string(s.Slot) }

type CallComboStatement struct {
	Name string
}

func (*CallComboStatement) statementNode()   {}
func (s *CallComboStatement) String() string { return s.Name }

package main

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

// ---------- LLInstructionWrapper ----------

type LLInstructionWrapper interface {
	LMCInstructions() []lmc.Instruction
	LMCDefs() []*lmc.DataInstr
	AddToLL(prog *lmc.Program)
	LLBase() []ir.Instruction
}

type LLInstructionBase struct {
	LLInstructionWrapper
	instrs []lmc.Instruction
	defs []*lmc.DataInstr
	base []ir.Instruction
}

func (base *LLInstructionBase) AddToLL(prog *lmc.Program) {
	prog.AddInstructions(base.LMCInstructions(), base.LMCDefs())
}

func (base *LLInstructionBase) LMCInstructions() []lmc.Instruction {
	return base.instrs
}

func (base *LLInstructionBase) LMCDefs() []*lmc.DataInstr {
	return base.defs
}

func (base *LLInstructionBase) LLBase() []ir.Instruction {
	return base.base
}

// ---------- LLBinaryInstr ---------

type LLBinaryInstr struct {
	LLInstructionBase
	x *lmc.Mailbox
	y *lmc.Mailbox
	dst *lmc.Mailbox
}

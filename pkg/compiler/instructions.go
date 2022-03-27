package compiler

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

// ---------- LLInstructionWrapper ----------

type LLInstructionWrapper interface {
	LMCInstructions() []lmc.Instruction
	LMCDefs() []*lmc.DataInstr
	LLBase() []ir.Instruction
}

type LLInstructionBase struct {
	LLInstructionWrapper
	base []ir.Instruction
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

func WrapBinaryInst(prog *lmc.Program, instr *ir.InstAdd) (*LLBinaryInstr, error) {
	var err error
	var x *lmc.Mailbox
	var y *lmc.Mailbox
	var dst *lmc.Mailbox

	x, err = GetValueMailbox(prog, instr.X)
	if err != nil {
		return nil, err
	}

	y, err = GetValueMailbox(prog, instr.Y)
	if err != nil {
		return nil, err
	}

	dst = prog.Memory.GetMailboxAddress(lmc.Address(instr.ID()))
	if dst == nil {
		return nil, UnknownMailboxError(lmc.Address(instr.ID()))
	}

	return &LLBinaryInstr{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		x: x,
		y: y,
		dst: dst,
	}, nil
}

// ---------- LLInstAdd ----------

type LLInstAdd struct {
	LLBinaryInstr
	instrs []lmc.Instruction
	defs []*lmc.DataInstr
}

func WrapLLInstAdd(prog *lmc.Program, instr *ir.InstAdd) (*LLInstAdd, error) {
	wrapped, err := WrapBinaryInst(prog, instr)
	if err != nil {
		return nil, err
	}

	return &LLInstAdd{
		LLBinaryInstr: *wrapped,
		instrs: []lmc.Instruction{
			lmc.NewLoadInstr(wrapped.x),
			lmc.NewAddInstr(wrapped.y),
			lmc.NewStoreInstr(wrapped.dst),
		},
		defs: nil,
	}, err
}

func (instr *LLInstAdd) LMCInstructions() []lmc.Instruction {
	return instr.instrs
}

func (instr *LLInstAdd) LMCDefs() []*lmc.DataInstr {
	return instr.defs
}

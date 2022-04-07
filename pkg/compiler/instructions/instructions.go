package instructions

import (
	c "github.com/clr1107/lmc-llvm-target/compiler"
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

// ---------- LLUnaryInstr ----------

type LLUnaryInstr struct {
	LLInstructionBase
	box *lmc.Mailbox
}

func WrapUnaryInstr(compiler *c.Compiler, instr ir.Instruction) (*LLUnaryInstr, error) {
	addr, err := c.ReflectGetLocalID(instr)
	if err != nil {
		return nil, err
	}

	box := compiler.Prog.Memory.GetMailboxAddress(addr)
	if box == nil {
		return nil, c.UnknownMailboxError(addr)
	}

	return &LLUnaryInstr{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		box: box,
	}, nil
}

// ---------- LLBinaryInstr ---------

type LLBinaryInstr struct {
	LLInstructionBase
	x *lmc.Mailbox
	y *lmc.Mailbox
	dst *lmc.Mailbox
}

func WrapBinaryInstr(compiler *c.Compiler, instr *ir.InstAdd) (*LLBinaryInstr, error) {
	var err error
	var x *lmc.Mailbox
	var y *lmc.Mailbox
	var dst *lmc.Mailbox

	x, err = c.MailboxFromLLValue(compiler, instr.X)
	if err != nil {
		return nil, err
	}

	y, err = c.MailboxFromLLValue(compiler, instr.Y)
	if err != nil {
		return nil, err
	}

	dst = compiler.Prog.Memory.GetMailboxAddress(lmc.Address(instr.ID()))
	if dst == nil {
		return nil, c.UnknownMailboxError(lmc.Address(instr.ID()))
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

package instructions

import (
	c "github.com/clr1107/lmc-llvm-target/compiler"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
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
	x *lmc.Mailbox
	dst *lmc.Mailbox
}

func WrapUnaryInstr(compiler *c.Compiler, instr ir.Instruction, X value.Value, dst *lmc.Mailbox) (*LLUnaryInstr, error) {
	x, err := c.MailboxFromLLValue(compiler, X)
	if err != nil {
		return nil, err
	}

	return &LLUnaryInstr{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		x: x,
		dst: dst,
	}, nil
}

// ---------- LLBinaryInstr ---------

type LLBinaryInstr struct {
	LLInstructionBase
	x *lmc.Mailbox
	y *lmc.Mailbox
	dst *lmc.Mailbox
}

func WrapBinaryInstr(compiler *c.Compiler, instr ir.Instruction, X value.Value, Y value.Value, dst *lmc.Mailbox) (*LLBinaryInstr, error) {
	x, err := c.MailboxFromLLValue(compiler, X)
	if err != nil {
		return nil, err
	}

	y, err := c.MailboxFromLLValue(compiler, Y)
	if err != nil {
		return nil, err
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

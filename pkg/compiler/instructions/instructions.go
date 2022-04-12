package instructions

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

// ---------- LLInstructionWrapper ----------

type LLInstructionWrapper interface {
	LMCInstructions() []lmc.Instruction
	LMCOps() []*lmc.MemoryOp
	LLBase() []ir.Instruction
}

type LLInstructionBase struct {
	LLInstructionWrapper
	base []ir.Instruction
}

func (base *LLInstructionBase) LLBase() []ir.Instruction {
	return base.base
}

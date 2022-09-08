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

// ---------- Other wrappers ----------

// ---------- WInstBitcast ----------

type WInstBitcast struct {
	LLInstructionBase
	From      *lmc.Mailbox
	To        *lmc.Mailbox
	memoryOps []*lmc.MemoryOp
}

func NewWInstBitcast(instr *ir.InstBitCast, fromBox *lmc.Mailbox, toBox *lmc.Mailbox, ops []*lmc.MemoryOp) *WInstBitcast {
	return &WInstBitcast{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		From:      fromBox,
		To:        toBox,
		memoryOps: ops,
	}
}

func (w *WInstBitcast) LMCInstructions() []lmc.Instruction {
	return []lmc.Instruction{
		lmc.NewLoadInstr(w.From),
		lmc.NewStoreInstr(w.To),
	}
}

func (w *WInstBitcast) LMCOps() []*lmc.MemoryOp {
	return w.memoryOps
}

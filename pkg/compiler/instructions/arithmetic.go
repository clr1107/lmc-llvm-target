package instructions

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

// ---------- WInstAdd ----------

type WInstAdd struct {
	LLInstructionBase
	X *lmc.Mailbox
	Y *lmc.Mailbox
	Dst       *lmc.Mailbox
	memoryOps []*lmc.MemoryOp
}

func NewWInstAdd(instr *ir.InstAdd, x *lmc.Mailbox, y *lmc.Mailbox, dst *lmc.Mailbox, ops []*lmc.MemoryOp) *WInstAdd {
	return &WInstAdd{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		X:         x,
		Y:         y,
		Dst:       dst,
		memoryOps: ops,
	}
}

func (w *WInstAdd) LMCInstructions() []lmc.Instruction {
	return []lmc.Instruction{
		lmc.NewLoadInstr(w.X),
		lmc.NewAddInstr(w.Y),
		lmc.NewStoreInstr(w.Dst),
	}
}

func (w *WInstAdd) LMCOps() []*lmc.MemoryOp {
	return w.memoryOps
}

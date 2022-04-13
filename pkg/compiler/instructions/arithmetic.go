package instructions

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

// ---------- WArithmeticInst ----------

type WArithmeticInst struct {
	LLInstructionBase
	X *lmc.Mailbox
	Y *lmc.Mailbox
	Dst *lmc.Mailbox
	memoryOps []*lmc.MemoryOp
}

func NewWArithmeticInst(instr ir.Instruction, x *lmc.Mailbox, y *lmc.Mailbox, dst *lmc.Mailbox, ops []*lmc.MemoryOp) *WArithmeticInst {
	return &WArithmeticInst{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		X: x,
		Y: y,
		Dst: dst,
		memoryOps: ops,
	}
}

// ---------- WInstAdd ----------

type WInstAdd struct {
	WArithmeticInst
}

func NewWInstAdd(inst *WArithmeticInst) *WInstAdd {
	return &WInstAdd{
		WArithmeticInst: *inst,
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

// ---------- WInstSub ----------

type WInstSub struct {
	WArithmeticInst
}

func NewWInstSub(inst *WArithmeticInst) *WInstSub {
	return &WInstSub{
		WArithmeticInst: *inst,
	}
}

func (w *WInstSub) LMCInstructions() []lmc.Instruction {
	return []lmc.Instruction{
		lmc.NewLoadInstr(w.X),
		lmc.NewSubInstr(w.Y),
		lmc.NewStoreInstr(w.Dst),
	}
}

func (w *WInstSub) LMCOps() []*lmc.MemoryOp {
	return w.memoryOps
}

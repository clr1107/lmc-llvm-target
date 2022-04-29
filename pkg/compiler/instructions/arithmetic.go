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

// ---------- WInstMul ----------

type WInstMul struct {
	WArithmeticInst
	Counter *lmc.Mailbox
	OneConst *lmc.Mailbox
	LoopLabel *lmc.Label
}

func NewWInstMul(inst *WArithmeticInst, counter *lmc.Mailbox, oneConst *lmc.Mailbox, loopLabel *lmc.Label, ops []*lmc.MemoryOp) *WInstMul {
	inst.memoryOps = append(inst.memoryOps, ops...)
	return &WInstMul{
		WArithmeticInst: *inst,
		Counter: counter,
		OneConst: oneConst,
		LoopLabel: loopLabel,
	}
}

func (w *WInstMul) LMCInstructions() []lmc.Instruction {
	return []lmc.Instruction{
		lmc.NewLoadInstr(w.X),
		lmc.NewStoreInstr(w.Counter),
		lmc.NewLabelled(w.LoopLabel, lmc.NewLoadInstr(w.Dst)),
		lmc.NewAddInstr(w.Y),
		lmc.NewStoreInstr(w.Dst),
		lmc.NewLoadInstr(w.Counter),
		lmc.NewSubInstr(w.OneConst),
		lmc.NewStoreInstr(w.Counter),
		lmc.NewBranchInstr(lmc.BRPositive, w.LoopLabel),
		lmc.NewLoadInstr(w.Dst),
		lmc.NewSubInstr(w.Y),
		lmc.NewStoreInstr(w.Dst),
	}
}

func (w *WInstMul) LMCOps() []*lmc.MemoryOp {
	return w.memoryOps
}

// ---------- WInstDiv ----------

type WInstDiv struct {
	WArithmeticInst
	Temp *lmc.Mailbox
	OneConst *lmc.Mailbox
	LoopLabel *lmc.Label
}

func NewWInstDiv(inst *WArithmeticInst, temp *lmc.Mailbox, oneConst *lmc.Mailbox, loopLabel *lmc.Label, ops []*lmc.MemoryOp) *WInstDiv {
	inst.memoryOps = append(inst.memoryOps, ops...)
	return &WInstDiv{
		WArithmeticInst: *inst,
		Temp: temp,
		OneConst: oneConst,
		LoopLabel: loopLabel,
	}
}

func (w *WInstDiv) LMCInstructions() []lmc.Instruction {
	return []lmc.Instruction{
		lmc.NewLoadInstr(w.X),
		lmc.NewStoreInstr(w.Temp),
		lmc.NewLabelled(w.LoopLabel, lmc.NewLoadInstr(w.Dst)),
		lmc.NewAddInstr(w.OneConst),
		lmc.NewStoreInstr(w.Dst),
		lmc.NewLoadInstr(w.Temp),
		lmc.NewSubInstr(w.Y),
		lmc.NewStoreInstr(w.Temp),
		lmc.NewBranchInstr(lmc.BRPositive, w.LoopLabel),
		lmc.NewLoadInstr(w.Dst),
		lmc.NewSubInstr(w.OneConst),
		lmc.NewStoreInstr(w.Dst),
	}
}

func (w *WInstDiv) LMCOps() []*lmc.MemoryOp {
	return w.memoryOps
}

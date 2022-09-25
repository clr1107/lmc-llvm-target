package instructions

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
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

type EmptyWInst struct {
	LLInstructionBase
}

func NewEmptyWInst(base []ir.Instruction) *EmptyWInst {
	return &EmptyWInst{LLInstructionBase{
		base: base,
	}}
}

func (w *EmptyWInst) LMCInstructions() []lmc.Instruction {
	return []lmc.Instruction{}
}

func (w *EmptyWInst) LMCOps() []*lmc.MemoryOp {
	return []*lmc.MemoryOp{}
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

// ---------- WInstICmp ----------

type iCmpMethod int

const (
	iCmpGtEq = iota
	iCmpLtEq
	iCmpGt
	iCmpLt
	iCmpEq
	_
	iCmpNeq
)

type WInstICmp struct {
	LLInstructionBase
	X         *lmc.Mailbox
	Y         *lmc.Mailbox
	Dst       *lmc.Mailbox
	oneConst  *lmc.Mailbox
	method    iCmpMethod
	memoryOps []*lmc.MemoryOp
}

func NewWInstICmp(instr *ir.InstICmp, x *lmc.Mailbox, y *lmc.Mailbox, dst *lmc.Mailbox, oneConst *lmc.Mailbox, ops []*lmc.MemoryOp) *WInstICmp {
	var m iCmpMethod

	switch instr.Pred {
	case enum.IPredEQ: // equals
		m = iCmpEq
	case enum.IPredNE:
		m = iCmpNeq
	case enum.IPredSGT:
		m = iCmpGt
	case enum.IPredSGE:
		m = iCmpGtEq
	case enum.IPredSLT:
		m = iCmpLt
	case enum.IPredSLE:
		m = iCmpLtEq
	default:
		panic("unsigned integer comparisons are unimplemented")
	}

	return &WInstICmp{
		X:         x,
		Y:         y,
		Dst:       dst,
		oneConst:  oneConst,
		method:    m,
		memoryOps: ops,
	}
}

func (w *WInstICmp) LMCInstructions() []lmc.Instruction {
	first := w.X
	second := w.Y

	if w.method&0x1 != 0 { // if < or <=
		first, second = second, first
	}

	instrs := []lmc.Instruction{
		lmc.NewLoadInstr(first),
		lmc.NewSubInstr(second),
	}

	if w.method&0x10 != 0 && w.method < 4 { // if >= or <=
		instrs = append(instrs, lmc.NewAddInstr(w.oneConst))
	}

	return append(instrs, lmc.NewStoreInstr(w.Dst))
}

func (w *WInstICmp) LMCOps() []*lmc.MemoryOp {
	return w.memoryOps
}

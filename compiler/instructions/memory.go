package instructions

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

// ---------- WInstAlloca ----------

type WInstAlloca struct {
	LLInstructionBase
	Box *lmc.Mailbox
	memoryOps []*lmc.MemoryOp
}

func NewWInstAlloca(instr *ir.InstAlloca, box *lmc.Mailbox, ops []*lmc.MemoryOp) *WInstAlloca {
	return &WInstAlloca{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		Box: box,
		memoryOps: ops,
	}
}

func (w *WInstAlloca) LMCInstructions() []lmc.Instruction {
	return nil
}

func (w *WInstAlloca) LMCOps() []*lmc.MemoryOp {
	return w.memoryOps
}

// ---------- WInstLoad ----------

type WInstLoad struct {
	LLInstructionBase
	X *lmc.Mailbox
	Dst *lmc.Mailbox
	memoryOps []*lmc.MemoryOp
}

func NewWInstLoad(instr *ir.InstLoad, x *lmc.Mailbox, dst *lmc.Mailbox, ops []*lmc.MemoryOp) *WInstLoad {
	return &WInstLoad{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		X: x,
		Dst: dst,
		memoryOps: ops,
	}
}

func (w *WInstLoad) LMCInstructions() []lmc.Instruction {
	return []lmc.Instruction{
		lmc.NewLoadInstr(w.X),
		lmc.NewStoreInstr(w.Dst),
	}
}

func (w *WInstLoad) LMCOps() []*lmc.MemoryOp {
	return w.memoryOps
}

// ---------- WInstStore ----------

type WInstStore struct {
	LLInstructionBase
	X *lmc.Mailbox
	Dst *lmc.Mailbox
	memoryOps []*lmc.MemoryOp
}

func NewWInstStore(instr *ir.InstStore, x *lmc.Mailbox, dst *lmc.Mailbox, ops []*lmc.MemoryOp) *WInstStore {
	return &WInstStore{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		X: x,
		Dst: dst,
		memoryOps: ops,
	}
}

func (w *WInstStore) LMCInstructions() []lmc.Instruction {
	return []lmc.Instruction{
		lmc.NewLoadInstr(w.X),
		lmc.NewStoreInstr(w.Dst),
	}
}

func (w *WInstStore) LMCOps() []*lmc.MemoryOp {
	return w.memoryOps
}

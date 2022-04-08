package instructions

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

// ---------- WInstAlloca ----------

type WInstAlloca struct {
	LLInstructionBase
	Box *lmc.Mailbox
}

func NewWInstAlloca(instr *ir.InstAlloca, box *lmc.Mailbox) *WInstAlloca {
	return &WInstAlloca{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		Box: box,
	}
}

func (w *WInstAlloca) LMCInstructions() []lmc.Instruction {
	return nil
}

func (w *WInstAlloca) LMCDefs() []*lmc.DataInstr {
	return []*lmc.DataInstr{
		lmc.NewDataInstr(0, w.Box),
	}
}

// ---------- WInstLoad ----------

type WInstLoad struct {
	LLInstructionBase
	X *lmc.Mailbox
	Dst *lmc.Mailbox
	newDstFlag bool
}

func NewWInstLoad(instr *ir.InstLoad, x *lmc.Mailbox, dst *lmc.Mailbox, newDstFlag bool) *WInstLoad {
	return &WInstLoad{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		X: x,
		Dst: dst,
		newDstFlag: newDstFlag,
	}
}

func (w *WInstLoad) LMCInstructions() []lmc.Instruction {
	return []lmc.Instruction{
		lmc.NewLoadInstr(w.X),
		lmc.NewStoreInstr(w.Dst),
	}
}

func (w *WInstLoad) LMCDefs() []*lmc.DataInstr {
	if w.newDstFlag {
		return []*lmc.DataInstr{
			lmc.NewDataInstr(0, w.Dst),
		}
	} else {
		return nil
	}
}

// ---------- WInstStore ----------

type WInstStore struct {
	LLInstructionBase
	X *lmc.Mailbox
	Dst *lmc.Mailbox
}

func NewWInstStore(instr *ir.InstStore, x *lmc.Mailbox, dst *lmc.Mailbox) *WInstStore {
	return &WInstStore{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		X: x,
		Dst: dst,
	}
}

func (w *WInstStore) LMCInstructions() []lmc.Instruction {
	return []lmc.Instruction{
		lmc.NewLoadInstr(w.X),
		lmc.NewStoreInstr(w.Dst),
	}
}

func (w *WInstStore) LMCDefs() []*lmc.DataInstr {
	return nil
}

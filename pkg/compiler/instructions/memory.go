package instructions

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
)

type WInstAlloca struct {
	LLInstructionBase
	Box *lmc.Mailbox
}

func NewWInstAlloca(box *lmc.Mailbox) *WInstAlloca {
	return &WInstAlloca{
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


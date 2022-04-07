package instructions

import (
	c "github.com/clr1107/lmc-llvm-target/compiler"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

// ---------- LLInstAlloca ----------

type LLInstAlloca struct {
	LLInstructionBase
	box *lmc.Mailbox
	defs []*lmc.DataInstr
}

func WrapInstAlloca(compiler *c.Compiler, instr *ir.InstAlloca) (*LLInstAlloca, error) {
	addr := lmc.Address(instr.ID())
	box, err := compiler.Prog.NewMailbox(addr, "")

	if err != nil {
		return nil, err
	}

	return &LLInstAlloca{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		box: box,
		defs: []*lmc.DataInstr{lmc.NewDataInstr(0, box)},
	}, nil
}

func (instr *LLInstAlloca) LMCInstructions() []lmc.Instruction {
	return nil
}

func (instr *LLInstAlloca) LMCDefs() []*lmc.DataInstr {
	return instr.defs
}

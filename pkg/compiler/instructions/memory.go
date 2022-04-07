package instructions

import (
	c "github.com/clr1107/lmc-llvm-target/compiler"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

// ---------- LLInstStore ----------

type LLInstStore struct {
	LLUnaryInstr
	instrs []lmc.Instruction
}

func WrapLLInstStore(compiler *c.Compiler, instr *ir.InstStore) (*LLInstStore, error) {
	dst, err := c.MailboxFromLLValue(compiler, instr.Dst)
	if err != nil {
		return nil, err
	}

	wrapped, err := WrapUnaryInstr(compiler, instr, instr.Src, dst)
	if err != nil {
		return nil, err
	}

	return &LLInstStore{
		LLUnaryInstr: *wrapped,
		instrs: []lmc.Instruction{
			lmc.NewLoadInstr(wrapped.x),
			lmc.NewStoreInstr(wrapped.dst),
		},
	}, nil
}

func (instr *LLInstStore) LMCInstructions() []lmc.Instruction {
	return instr.instrs
}

func (instr *LLInstStore) LMCDefs() []*lmc.DataInstr {
	return nil
}

// ---------- LLInstAlloca ----------

type LLInstAlloca struct {
	LLInstructionBase
	box *lmc.Mailbox
	defs []*lmc.DataInstr
}

func WrapInstAlloca(compiler *c.Compiler, instr *ir.InstAlloca) (*LLInstAlloca, error) {
	box := lmc.NewMailbox(0, "")
	if err := compiler.Prog.Memory.AddMailbox(box); err != nil {
		return nil, err
	}

	return &LLInstAlloca{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		box: box,
		defs: []*lmc.DataInstr{
			lmc.NewDataInstr(0, box),
		},
	}, nil
}

func (instr *LLInstAlloca) LMCInstructions() []lmc.Instruction {
	return nil
}

func (instr *LLInstAlloca) LMCDefs() []*lmc.DataInstr {
	return instr.defs
}

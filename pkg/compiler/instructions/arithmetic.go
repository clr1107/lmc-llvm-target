package instructions

import (
	c "github.com/clr1107/lmc-llvm-target/compiler"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

// ---------- LLInstAdd ----------

type LLInstAdd struct {
	LLBinaryInstr
	instrs []lmc.Instruction
}

func WrapLLInstAdd(compiler *c.Compiler, instr *ir.InstAdd) (*LLInstAdd, error) {
	dst := compiler.Prog.Memory.GetMailboxAddress(lmc.Address(instr.ID()))
	if dst == nil {
		return nil, c.UnknownMailboxError(lmc.Address(instr.ID()))
	}

	wrapped, err := WrapBinaryInstr(compiler, instr, instr.X, instr.Y, dst)
	if err != nil {
		return nil, err
	}

	return &LLInstAdd{
		LLBinaryInstr: *wrapped,
		instrs: []lmc.Instruction{
			lmc.NewLoadInstr(wrapped.x),
			lmc.NewAddInstr(wrapped.y),
			lmc.NewStoreInstr(wrapped.dst),
		},
	}, nil
}

func (instr *LLInstAdd) LMCInstructions() []lmc.Instruction {
	return instr.instrs
}

func (instr *LLInstAdd) LMCDefs() []*lmc.DataInstr {
	return nil
}

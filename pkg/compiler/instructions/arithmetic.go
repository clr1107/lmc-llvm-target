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
	defs []*lmc.DataInstr
}

func WrapLLInstAdd(compiler *c.Compiler, instr *ir.InstAdd) (*LLInstAdd, error) {
	wrapped, err := WrapBinaryInstr(compiler, instr)
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
		defs: nil,
	}, err
}

func (instr *LLInstAdd) LMCInstructions() []lmc.Instruction {
	return instr.instrs
}

func (instr *LLInstAdd) LMCDefs() []*lmc.DataInstr {
	return instr.defs
}

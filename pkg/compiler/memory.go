package compiler

import (
	"github.com/clr1107/lmc-llvm-target/compiler/instructions"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

func (compiler *Compiler) WrapLLInstAlloca(instr *ir.InstAlloca) (*instructions.WInstAlloca, error) {
	addr := lmc.Address(instr.ID())
	box := compiler.Prog.Memory.NewMailbox(addr, "")

	if err := compiler.Prog.Memory.AddMailbox(box); err != nil {
		return nil, err
	}

	return instructions.NewWInstAlloca(box), nil
}

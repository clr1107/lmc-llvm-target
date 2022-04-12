package compiler

import (
	"github.com/clr1107/lmc-llvm-target/compiler/instructions"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

func (compiler *Compiler) WrapLLInstAdd(instr *ir.InstAdd) (*instructions.WInstAdd, error) {
	var xBox *lmc.Mailbox
	var yBox *lmc.Mailbox
	var dstBox *lmc.Mailbox
	var ops []*lmc.MemoryOp
	var op *lmc.MemoryOp
	var err error

	op, err = compiler.GetMailboxFromLL(instr.X)
	if err != nil {
		return nil, err
	} else {
		xBox = op.Boxes[0].Box
		ops = append(ops, op)
	}

	op, err = compiler.GetMailboxFromLL(instr.Y)
	if err != nil {
		return nil, err
	} else {
		yBox = op.Boxes[0].Box
		ops = append(ops, op)
	}

	dstAddr := lmc.Address(instr.ID())
	dstBox = compiler.Prog.Memory.GetMailboxAddress(dstAddr)
	if dstBox == nil {
		op := compiler.Prog.Memory.NewMailbox(dstAddr, "")
		dstBox = op.Boxes[0].Box

		ops = append(ops, op)
	}

	return instructions.NewWInstAdd(instr, xBox, yBox, dstBox, ops), nil
}

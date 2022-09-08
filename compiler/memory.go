package compiler

import (
	"github.com/clr1107/lmc-llvm-target/compiler/instructions"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

func (compiler *Compiler) WrapLLInstAlloca(instr *ir.InstAlloca) *Compilation {
	addr := lmc.Address(instr.ID())
	op := compiler.Prog.Memory.NewMailbox(addr, "")

	return &Compilation{Wrapped: instructions.NewWInstAlloca(instr, op.Boxes[0].Box, []*lmc.MemoryOp{op})}
}

func (compiler *Compiler) WrapLLInstLoad(instr *ir.InstLoad) *Compilation {
	var srcBox *lmc.Mailbox
	var dstBox *lmc.Mailbox
	var ops []*lmc.MemoryOp
	var op *lmc.MemoryOp
	var err error

	op, err = compiler.GetMailboxFromLL(instr.Src)
	if err != nil {
		return &Compilation{Err: err}
	} else {
		srcBox = op.Boxes[0].Box
		ops = append(ops, op)
	}

	dstAddr := lmc.Address(instr.ID())
	dstBox = compiler.Prog.Memory.GetMailboxAddress(dstAddr)

	if dstBox == nil {
		op := compiler.Prog.Memory.NewMailbox(dstAddr, "")
		dstBox = op.Boxes[0].Box

		ops = append(ops, op)
	}

	return &Compilation{Wrapped: instructions.NewWInstLoad(instr, srcBox, dstBox, ops)}
}

func (compiler *Compiler) WrapLLInstStore(instr *ir.InstStore) *Compilation {
	var srcBox *lmc.Mailbox
	var dstBox *lmc.Mailbox
	var ops []*lmc.MemoryOp
	var op *lmc.MemoryOp
	var err error

	op, err = compiler.GetMailboxFromLL(instr.Src)
	if err != nil {
		return &Compilation{Err: err}
	} else {
		srcBox = op.Boxes[0].Box
		ops = append(ops, op)
	}

	op, err = compiler.GetMailboxFromLL(instr.Dst)
	if err != nil {
		return &Compilation{Err: err}
	} else {
		dstBox = op.Boxes[0].Box
		ops = append(ops, op)
	}

	return &Compilation{Wrapped: instructions.NewWInstStore(instr, srcBox, dstBox, ops)}
}

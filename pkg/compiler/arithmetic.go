package compiler

import (
	"github.com/clr1107/lmc-llvm-target/compiler/instructions"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

func (compiler *Compiler) WrapArithmeticLLInst(instr ir.Instruction, x value.Value, y value.Value, addr lmc.Address) (*instructions.WArithmeticInst, error) {
	var xBox *lmc.Mailbox
	var yBox *lmc.Mailbox
	var dstBox *lmc.Mailbox
	var ops []*lmc.MemoryOp
	var op *lmc.MemoryOp
	var err error

	op, err = compiler.GetMailboxFromLL(x)
	if err != nil {
		return nil, err
	} else {
		xBox = op.Boxes[0].Box
		ops = append(ops, op)
	}

	op, err = compiler.GetMailboxFromLL(y)
	if err != nil {
		return nil, err
	} else {
		yBox = op.Boxes[0].Box
		ops = append(ops, op)
	}

	dstBox = compiler.Prog.Memory.GetMailboxAddress(addr)
	if dstBox == nil {
		op := compiler.Prog.Memory.NewMailbox(addr, "")
		dstBox = op.Boxes[0].Box

		ops = append(ops, op)
	}

	return instructions.NewWArithmeticInst(instr, xBox, yBox, dstBox, ops), nil
}

func (compiler *Compiler) WrapLLInstAdd(instr *ir.InstAdd) (*instructions.WInstAdd, error) {
	var arithmetic *instructions.WArithmeticInst
	var err error

	arithmetic, err = compiler.WrapArithmeticLLInst(instr, instr.X, instr.Y, lmc.Address(instr.ID()))
	if err != nil {
		return nil, err
	}

	return instructions.NewWInstAdd(arithmetic), nil
}

func (compiler *Compiler) WrapLLInstSub(instr *ir.InstSub) (*instructions.WInstSub, error) {
	var arithmetic *instructions.WArithmeticInst
	var err error

	arithmetic, err = compiler.WrapArithmeticLLInst(instr, instr.X, instr.Y, lmc.Address(instr.ID()))
	if err != nil {
		return nil, err
	}

	return instructions.NewWInstSub(arithmetic), nil
}

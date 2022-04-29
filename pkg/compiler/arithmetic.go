package compiler

import (
	"github.com/clr1107/lmc-llvm-target/compiler/instructions"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

func (compiler *Compiler) wrapArithmeticInst(instr ir.Instruction, x value.Value, y value.Value, addr lmc.Address) (*instructions.WArithmeticInst, error) {
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
	if wrapped, err := compiler.wrapArithmeticInst(instr, instr.X, instr.Y, lmc.Address(instr.ID())); err != nil {
		return nil, err
	} else {
		return instructions.NewWInstAdd(wrapped), nil
	}
}

func (compiler *Compiler) WrapLLInstSub(instr *ir.InstSub) (*instructions.WInstSub, error) {
	if wrapped, err := compiler.wrapArithmeticInst(instr, instr.X, instr.Y, lmc.Address(instr.ID())); err != nil {
		return nil, err
	} else {
		return instructions.NewWInstSub(wrapped), nil
	}
}

func (compiler *Compiler) WrapLLInstMul(instr *ir.InstMul) (*instructions.WInstMul, error) {
	if wrapped, err := compiler.wrapArithmeticInst(instr, instr.X, instr.Y, lmc.Address(instr.ID())); err != nil {
		return nil, err
	} else {
		tempOp := compiler.GetTempBox()
		oneOp := compiler.Prog.Memory.Constant(1)
		labelOp := compiler.Prog.Memory.NewLabel("")

		return instructions.NewWInstMul(
			wrapped,
			tempOp.Boxes[0].Box,
			oneOp.Boxes[0].Box,
			labelOp.Labels[0].Label,
			[]*lmc.MemoryOp{tempOp, oneOp, labelOp},
		), nil
	}
}

func (compiler *Compiler) WrapLLInstDiv(instr ir.Instruction, X value.Value, Y value.Value, id int64) (*instructions.WInstDiv, error) {
	if wrapped, err := compiler.wrapArithmeticInst(instr, X, Y, lmc.Address(id)); err != nil {
		return nil, err
	} else {
		tempOp := compiler.GetTempBox()
		oneOp := compiler.Prog.Memory.Constant(1)
		labelOp := compiler.Prog.Memory.NewLabel("")

		return instructions.NewWInstDiv(
			wrapped,
			tempOp.Boxes[0].Box,
			oneOp.Boxes[0].Box,
			labelOp.Labels[0].Label,
			[]*lmc.MemoryOp{tempOp, oneOp, labelOp},
		), nil
	}
}

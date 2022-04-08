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
	var newDstFlag bool
	var err error

	xBox, err = compiler.GetMailboxFromLL(instr.X)
	if err != nil {
		return nil, err
	}

	yBox, err = compiler.GetMailboxFromLL(instr.Y)
	if err != nil {
		return nil, err
	}

	dstAddr := lmc.Address(instr.ID())
	dstBox = compiler.Prog.Memory.GetMailboxAddress(dstAddr)
	if dstBox == nil {
		dstBox = compiler.Prog.Memory.NewMailbox(dstAddr, "")
		err = compiler.Prog.Memory.AddMailbox(dstBox)
		if err != nil {
			return nil, err
		}

		newDstFlag = true
	}

	return instructions.NewWInstAdd(instr, xBox, yBox, dstBox, newDstFlag), nil
}

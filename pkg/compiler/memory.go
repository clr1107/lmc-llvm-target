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

	return instructions.NewWInstAlloca(instr, box), nil
}

func (compiler *Compiler) WrapLLInstLoad(instr *ir.InstLoad) (*instructions.WInstLoad, error) {
	var newDstFlag bool
	var srcBox *lmc.Mailbox
	var dstBox *lmc.Mailbox
	var err error

	srcBox, err = compiler.GetMailboxFromLL(instr.Src)
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

	return instructions.NewWInstLoad(instr, srcBox, dstBox, newDstFlag), nil
}

func (compiler *Compiler) WrapLLInstStore(instr *ir.InstStore) (*instructions.WInstStore, error) {
	var srcBox *lmc.Mailbox
	var dstBox *lmc.Mailbox
	var err error

	srcBox, err = compiler.GetMailboxFromLL(instr.Src)
	if err != nil {
		return nil, err
	}

	dstBox, err = compiler.GetMailboxFromLL(instr.Dst)
	if err != nil {
		return nil, err
	}

	return instructions.NewWInstStore(instr, srcBox, dstBox), nil
}

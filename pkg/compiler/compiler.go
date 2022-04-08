package compiler

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
)

type Compiler struct {
	Prog *lmc.Program
	tempBox *lmc.Mailbox
}

func NewCompiler(prog *lmc.Program) *Compiler {
	return &Compiler{
		Prog: prog,
		tempBox: nil,
	}
}

func (compiler *Compiler) GetTempBox() (*lmc.Mailbox, error) {
	if compiler.tempBox != nil {
		return compiler.tempBox, nil
	}

	box, err := compiler.Prog.NewMailbox(-1, "_TEMP")
	if err != nil {
		compiler.tempBox = box
	}

	return box, err
}

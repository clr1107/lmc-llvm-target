package compiler

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"reflect"
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

func (compiler *Compiler) GetMailboxFromLL(ll interface{}) (*lmc.Mailbox, error) {
	switch ll.(type) {
	case *constant.Null:
		return compiler.GetTempBox()
	case *constant.Int:
		return compiler.Prog.Constant(lmc.Value(ll.(*constant.Int).X.Int64()))
	case ir.Instruction: // last try, just use reflection lol
		id, err := ReflectGetLocalID(ll)
		if err != nil {
			return nil, err
		}

		mbox := compiler.Prog.Memory.GetMailboxAddress(id)
		if mbox == nil {
			return nil, UnknownMailboxError(id)
		}

		return mbox, nil
	default:
		return nil, InvalidLLTypeError(reflect.TypeOf(ll).String())
	}
}

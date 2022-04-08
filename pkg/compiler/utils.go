package compiler

import (
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"reflect"
)

func GetLLFunc(funcs []*ir.Func, name string) *ir.Func {
	for _, v := range funcs {
		if v.Name() == name {
			return v
		}
	}

	return nil
}

func GetLLInstrs(blocks []*ir.Block) []ir.Instruction {
	instrs := make([]ir.Instruction, 0)

	for _, block := range blocks {
		instrs = append(instrs, block.Insts...)
	}

	return instrs
}

func GetLLEntry(module *ir.Module) *ir.Func {
	f := GetLLFunc(module.Funcs, "_lmc")

	if f == nil {
		return nil
	}

	if len(f.Params) != 0 || f.Sig.RetType.String() != "void" {
		return nil
	}

	return f
}

func ReflectGetLocalID(x interface{}) (lmc.Address, error) {
	f := reflect.ValueOf(x).MethodByName("ID")
	if f.IsZero() {
		return -1, NonexistentPropertyError("#ID()")
	}

	res := f.Call(nil)
	if len(res) != 1 {
		return -1, NonexistentPropertyError("#ID() -> 1")

	}

	id := res[0]
	if id.Kind() != reflect.Int64 {
		return -1, IncorrectTypesError("ID[int64]", "ID")
	}

	return lmc.Address(id.Int()), nil
}

func GetMailboxFromLL(compiler *Compiler, ll interface{}) (*lmc.Mailbox, error) {
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

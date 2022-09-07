package compiler

import (
	"github.com/clr1107/lmc-llvm-target/compiler/errors"
	"github.com/clr1107/lmc-llvm-target/compiler/instructions"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"reflect"
)

func wrapBuiltinFunc(name string, instr *ir.InstCall, params []*lmc.Mailbox, ops []*lmc.MemoryOp) *instructions.WBuiltinCall {
	switch name {
	case "input":
		return instructions.NewWBuiltinCall(instr, instructions.NewBuiltinInput(), params, ops)
	case "output":
		return instructions.NewWBuiltinCall(instr, instructions.NewBuiltinOutput(), params, ops)
	case "_hlt":
		return instructions.NewWBuiltinCall(instr, instructions.NewBuiltinHalt(), params, ops)
	default:
		return nil
	}
}

func (compiler *Compiler) WrapLLInstCall(instr *ir.InstCall) (*instructions.WBuiltinCall, error) {
	var f *ir.Func
	var w *instructions.WBuiltinCall
	var ok bool
	var err error

	if f, ok = instr.Callee.(*ir.Func); !ok {
		return nil, errors.E_IncorrectType(nil, "*ir.Func.Callee", reflect.TypeOf(instr.Callee).String(), "*ir.Func")
	}

	var op *lmc.MemoryOp
	var ops []*lmc.MemoryOp
	var params []*lmc.Mailbox

	for _, a := range instr.Args {
		if op, err = compiler.GetMailboxFromLL(a); err != nil {
			return nil, err
		}

		ops = append(ops, op)
		params = append(params, op.Boxes[0].Box)
	}

	if w = wrapBuiltinFunc(f.Name(), instr, params, ops); w == nil {
		return nil, errors.E_UnknownBuiltin(f, nil)
	}

	if err = w.Invoke(); err != nil {
		return nil, err
	}

	return w, nil
}

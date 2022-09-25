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
	// Instruction functions
	case "_hlt":
		return instructions.NewWBuiltinCall(instr, instructions.NewBuiltinHltInstr(), params, ops)
	case "_inp":
		return instructions.NewWBuiltinCall(instr, instructions.NewBuiltinInpInstr(), params, ops)
	case "_out":
		return instructions.NewWBuiltinCall(instr, instructions.NewBuiltinOutInstr(), params, ops)
	case "_sta":
		return instructions.NewWBuiltinCall(instr, instructions.NewBuiltinStaInstr(), params, ops)
	// Other
	case "input":
		return instructions.NewWBuiltinCall(instr, instructions.NewBuiltinInput(), params, ops)
	case "output":
		return instructions.NewWBuiltinCall(instr, instructions.NewBuiltinOutput(), params, ops)
	default:
		return nil
	}
}

func (compiler *Compiler) WrapLLInstCall(instr *ir.InstCall) *Compilation {
	var f *ir.Func
	var w *instructions.WBuiltinCall
	var ok bool
	var err error

	if f, ok = instr.Callee.(*ir.Func); !ok {
		return &Compilation{
			Err: errors.E_IncorrectType(nil, "*ir.Func.Callee", reflect.TypeOf(instr.Callee).String(), "*ir.Func"),
		}
	}

	var op *lmc.MemoryOp
	var ops []*lmc.MemoryOp
	var params []*lmc.Mailbox

	if f.Name() == "__lmc_option__" { // protection against accidental pattern matching. See compOptionPattern
		return &Compilation{
			Err: errors.E_Err("pattern matched __lmc_option__ to function call not syntax", nil),
		}
	}

	for _, a := range instr.Args {
		if op, err = compiler.GetMailboxFromLL(a); err != nil {
			return &Compilation{Err: err}
		}

		ops = append(ops, op)
		params = append(params, op.Boxes[0].Box)
	}

	if w = wrapBuiltinFunc(f.Name(), instr, params, ops); w == nil {
		return &Compilation{Err: errors.E_UnknownBuiltin(f, nil)}
	}

	if err = w.Invoke(); err != nil {
		return &Compilation{Err: err}
	}

	return &Compilation{Wrapped: w}
}

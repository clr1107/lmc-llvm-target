package compiler

import (
	"github.com/clr1107/lmc-llvm-target/compiler/errors"
	"github.com/clr1107/lmc-llvm-target/compiler/instructions"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/value"
	"reflect"
)

type Compiler struct {
	Prog    *lmc.Program
	tempBox *lmc.Mailbox
}

func NewCompiler(prog *lmc.Program) *Compiler {
	return &Compiler{
		Prog:    prog,
		tempBox: nil,
	}
}

func (compiler *Compiler) GetTempBox() *lmc.MemoryOp {
	if compiler.tempBox != nil {
		return lmc.NewMemoryOpBox1(compiler.tempBox, false)
	}

	op := compiler.Prog.Memory.NewMailbox(-1, "_TEMP")
	compiler.tempBox = op.Boxes[0].Box

	return op
}

func (compiler *Compiler) GetMailboxFromLL(ll interface{}) (*lmc.MemoryOp, error) {
	switch x := ll.(type) {
	case *constant.Null:
		return compiler.GetTempBox(), nil
	case *constant.Int:
		return compiler.Prog.Memory.Constant(lmc.Value(x.X.Int64())), nil
	//case *ir.Param:
	case value.Value: // last try, just use reflection lol
		if !ValidLLType(x.Type()) {
			return nil, errors.E_InvalidLLTypes(nil, x.Type().String())
		}

		id, err := ReflectGetLocalID(ll)
		if err != nil {
			return nil, err
		}

		mbox := compiler.Prog.Memory.GetMailboxAddress(id)
		if mbox == nil {
			return nil, errors.E_UnknownMailbox(id, nil)
		}

		return lmc.NewMemoryOpBox1(mbox, false), nil
	default:
		return nil, errors.E_InvalidLLTypes(nil, reflect.TypeOf(ll).String())
	}
}

func (compiler *Compiler) CompileInst(instr ir.Instruction) (instructions.LLInstructionWrapper, error) {
	switch cast := instr.(type) {
	// arithmetic
	case *ir.InstAdd:
		return compiler.WrapLLInstAdd(cast)
	case *ir.InstSub:
		return compiler.WrapLLInstSub(cast)
	case *ir.InstMul:
		return compiler.WrapLLInstMul(cast)
	case *ir.InstSDiv:
		return compiler.WrapLLInstDiv(cast, cast.X, cast.Y, cast.ID())
	case *ir.InstUDiv:
		return compiler.WrapLLInstDiv(cast, cast.X, cast.Y, cast.ID())
	case *ir.InstSRem:
		return compiler.WrapLLInstRem(cast, cast.X, cast.Y, cast.ID())
	case *ir.InstURem:
		return compiler.WrapLLInstRem(cast, cast.X, cast.Y, cast.ID())
	// memory
	case *ir.InstAlloca:
		return compiler.WrapLLInstAlloca(cast)
	case *ir.InstLoad:
		return compiler.WrapLLInstLoad(cast)
	case *ir.InstStore:
		return compiler.WrapLLInstStore(cast)
	// other
	case *ir.InstCall:
		return compiler.WrapLLInstCall(cast)
	// unknown
	default:
		return nil, errors.E_UnknownLLInstruction(instr, nil)
	}
}

func (compiler *Compiler) AddCompiledInstruction(instr instructions.LLInstructionWrapper) error {
	var defs []*lmc.DataInstr

	// Consider in the future using *Program#AddMemoryOp
	for _, op := range instr.LMCOps() {
		for _, box := range op.GetNewBoxes() {
			if err := compiler.Prog.Memory.AddMailbox(box); err != nil {
				return errors.E_LMC("adding compiled instruction -- dev: CONSIDER *Program#AddMemoryOp", err)
			}
		}

		for _, label := range op.GetNewLabels() {
			if err := compiler.Prog.Memory.AddLabel(label); err != nil {
				return errors.E_LMC("adding compiled instr label -- dev: CONSIDER *Program#AddMemoryOp", err)
			}
		}

		defs = append(defs, op.Defs()...)
	}

	compiler.Prog.AddInstructions(instr.LMCInstructions(), defs)
	return nil
}

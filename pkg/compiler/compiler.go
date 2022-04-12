package compiler

import (
	"github.com/clr1107/lmc-llvm-target/compiler/instructions"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
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
		return lmc.NewMemoryOp1(compiler.tempBox, false)
	}

	op := compiler.Prog.Memory.NewMailbox(-1, "_TEMP")
	compiler.tempBox = op.Boxes[0].Box

	return op
}

func (compiler *Compiler) GetMailboxFromLL(ll interface{}) (*lmc.MemoryOp, error) {
	switch ll.(type) {
	case *constant.Null:
		return compiler.GetTempBox(), nil
	case *constant.Int:
		return compiler.Prog.Memory.Constant(lmc.Value(ll.(*constant.Int).X.Int64())), nil
	case ir.Instruction: // last try, just use reflection lol
		id, err := ReflectGetLocalID(ll)
		if err != nil {
			return nil, err
		}

		mbox := compiler.Prog.Memory.GetMailboxAddress(id)
		if mbox == nil {
			return nil, UnknownMailboxError(id)
		}

		return lmc.NewMemoryOp1(mbox, false), nil
	default:
		return nil, InvalidLLTypeError(reflect.TypeOf(ll).String())
	}
}

func (compiler *Compiler) CompileInst(instr ir.Instruction) (instructions.LLInstructionWrapper, error) {
	switch instr.(type) {
	// arithmetic
	case *ir.InstAdd:
		return compiler.WrapLLInstAdd(instr.(*ir.InstAdd))
	// memory
	case *ir.InstAlloca:
		return compiler.WrapLLInstAlloca(instr.(*ir.InstAlloca))
	case *ir.InstLoad:
		return compiler.WrapLLInstLoad(instr.(*ir.InstLoad))
	case *ir.InstStore:
		return compiler.WrapLLInstStore(instr.(*ir.InstStore))
	// unknown
	default:
		return nil, UnknownLLInstructionError(instr)
	}
}

func (compiler *Compiler) AddCompiledInstruction(instr instructions.LLInstructionWrapper) error {
	var defs []*lmc.DataInstr

	for _, op := range instr.LMCOps() {
		for _, box := range op.GetNew() {
			if err := compiler.Prog.Memory.AddMailbox(box); err != nil {
				return err
			}
		}

		defs = append(defs, op.Defs()...)
	}

	compiler.Prog.AddInstructions(instr.LMCInstructions(), defs)
	return nil
}

package instructions

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/compiler/errors"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
)

// ---------- Builtins structure ----------

type BuiltinId uint8

const (
	B_Inp BuiltinId = iota
	B_Out
	B_Hlt
	B_Sta
	B_Input
	B_Output
)

type BuiltinReturn struct {
	Instructions []lmc.Instruction
	Ops          []*lmc.MemoryOp
	Err          error
}

type Builtin interface {
	fmt.Stringer
	Id() BuiltinId
	Name() string
	Parameters() int
	Call(params []*lmc.Mailbox) *BuiltinReturn
}

type BuiltinBase struct {
	Builtin
	id         BuiltinId
	name       string
	parameters int
}

func (b *BuiltinBase) Id() BuiltinId {
	return b.id
}

func (b *BuiltinBase) Name() string {
	return b.name
}

func (b *BuiltinBase) Parameters() int {
	return b.parameters
}

func (b *BuiltinBase) String() string {
	return fmt.Sprintf("builtin %s(%d)", b.Name(), b.Parameters())
}

func (b *BuiltinBase) checkParams(params []*lmc.Mailbox) error {
	if len(params) != b.parameters {
		return fmt.Errorf("got %d parameters expected %d", len(params), b.parameters)
	}

	return nil
}

// ---------- Builtins definitions ----------

// ---------- Input function ----------
// As defined in compiler/lmc.h

type BuiltinInput struct {
	BuiltinBase
}

func NewBuiltinInput() *BuiltinInput {
	return &BuiltinInput{
		BuiltinBase{
			id:         B_Input,
			name:       "input",
			parameters: 1,
		},
	}
}

func (b *BuiltinInput) Call(params []*lmc.Mailbox) *BuiltinReturn {
	var ret BuiltinReturn

	if err := b.checkParams(params); err != nil {
		ret.Err = err
	} else {
		ret.Instructions = []lmc.Instruction{
			lmc.NewInputInstr(),
			lmc.NewStoreInstr(params[0]),
		}
		ret.Ops = []*lmc.MemoryOp{}
	}

	return &ret
}

// ---------- Output function ----------
// As defined in compiler/lmc.h

type BuiltinOutput struct {
	BuiltinBase
}

func NewBuiltinOutput() *BuiltinOutput {
	return &BuiltinOutput{
		BuiltinBase{
			id:         B_Output,
			name:       "output",
			parameters: 1,
		},
	}
}

func (b *BuiltinOutput) Call(params []*lmc.Mailbox) *BuiltinReturn {
	var ret BuiltinReturn

	if err := b.checkParams(params); err != nil {
		ret.Err = err
	} else {
		ret.Instructions = []lmc.Instruction{
			lmc.NewLoadInstr(params[0]),
			lmc.NewOutputInstr(),
		}
		ret.Ops = []*lmc.MemoryOp{}
	}

	return &ret
}

// ---------- `STA` function ----------
// As defined in compiler/lmc.h

type BuiltinStaInstr struct {
	BuiltinBase
}

func NewBuiltinStaInstr() *BuiltinStaInstr {
	return &BuiltinStaInstr{
		BuiltinBase{
			id:         B_Sta,
			name:       "sta",
			parameters: 1,
		},
	}
}

func (b *BuiltinStaInstr) Call(params []*lmc.Mailbox) *BuiltinReturn {
	var ret BuiltinReturn

	if err := b.checkParams(params); err != nil {
		ret.Err = err
	} else {
		ret.Instructions = []lmc.Instruction{
			lmc.NewStoreInstr(params[0]),
		}
		ret.Ops = []*lmc.MemoryOp{}
	}

	return &ret
}

// ---------- Unary instruction functions ----------

type BuiltinUnaryInstrFunc struct {
	BuiltinBase
	Instruction lmc.Instruction
}

func NewBuiltinInstrFunc(id BuiltinId, name string, instr lmc.Instruction) *BuiltinUnaryInstrFunc {
	return &BuiltinUnaryInstrFunc{
		Instruction: instr,
		BuiltinBase: BuiltinBase{
			id:         id,
			name:       name,
			parameters: 0,
		},
	}
}

func (b *BuiltinUnaryInstrFunc) Call(params []*lmc.Mailbox) *BuiltinReturn {
	var ret BuiltinReturn

	if err := b.checkParams(params); err != nil {
		ret.Err = err
	} else {
		ret.Instructions = []lmc.Instruction{
			b.Instruction,
		}
		ret.Ops = []*lmc.MemoryOp{}
	}

	return &ret
}

func NewBuiltinOutInstr() *BuiltinUnaryInstrFunc {
	return NewBuiltinInstrFunc(B_Out, "out", lmc.NewOutputInstr())
}

func NewBuiltinInpInstr() *BuiltinUnaryInstrFunc {
	return NewBuiltinInstrFunc(B_Inp, "inp", lmc.NewInputInstr())
}

func NewBuiltinHltInstr() *BuiltinUnaryInstrFunc {
	return NewBuiltinInstrFunc(B_Hlt, "hlt", lmc.NewHaltInstr())
}

// ---------- WBuiltinCall ----------

type WBuiltinCall struct {
	LLInstructionBase
	Func        Builtin
	Parameters  []*lmc.Mailbox
	originalOps []*lmc.MemoryOp
	invocation  *BuiltinReturn
	Invoked     bool
}

func NewWBuiltinCall(instr *ir.InstCall, f Builtin, params []*lmc.Mailbox, ops []*lmc.MemoryOp) *WBuiltinCall {
	return &WBuiltinCall{
		LLInstructionBase: LLInstructionBase{
			base: []ir.Instruction{instr},
		},
		Func:        f,
		Parameters:  params,
		originalOps: ops,
		Invoked:     false,
	}
}

func (w *WBuiltinCall) Invoke() error {
	w.invocation = w.Func.Call(w.Parameters)

	if w.invocation.Err != nil {
		return errors.E_BuiltinInvocation(w.Func.String(), w.invocation.Err)
	}

	w.invocation.Ops = append(w.originalOps, w.invocation.Ops...)
	w.Invoked = true

	return nil
}

func (w *WBuiltinCall) LMCInstructions() []lmc.Instruction {
	if !w.Invoked {
		panic("builtin not invoked, cannot return LMC instructions during compilation")
	}

	return w.invocation.Instructions
}

func (w *WBuiltinCall) LMCOps() []*lmc.MemoryOp {
	if !w.Invoked {
		panic("builtin not invoked, cannot return memory ops during compilation")
	}

	return w.invocation.Ops
}

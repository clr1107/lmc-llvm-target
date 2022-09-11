package compiler

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/compiler/errors"
	"github.com/clr1107/lmc-llvm-target/compiler/instructions"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/value"
	"reflect"
	"sort"
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
			return nil, errors.E_InvalidLLTypes(nil, x.Type().LLString())
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

type Compilation struct {
	Wrapped  instructions.LLInstructionWrapper
	Err      error
	Warnings []*errors.Warning
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

// ---------- Wrapping of other instructions ----------

func (compiler *Compiler) WrapLLBitcast(instr *ir.InstBitCast) *Compilation {
	var compilation Compilation

	if ValidLLType(instr.From.Type()) && ValidLLType(instr.To) {
		compilation.Warnings = []*errors.Warning{
			errors.W_Bitcast(instr.From.Type().LLString(), instr.To.LLString()),
		}
	} else {
		compilation.Err = errors.E_InvalidLLTypes(nil, instr.From.Type().LLString(), instr.To.LLString())
		return &compilation
	}

	var fromBox *lmc.Mailbox
	var toBox *lmc.Mailbox
	var ops []*lmc.MemoryOp
	var op *lmc.MemoryOp
	var err error

	if op, err = compiler.GetMailboxFromLL(instr.From); err != nil {
		compilation.Err = err
		return &compilation
	} else {
		fromBox = op.Boxes[0].Box
		ops = append(ops, op)
	}

	toBox = compiler.Prog.Memory.GetMailboxAddress(lmc.Address(instr.ID()))
	if toBox == nil {
		op := compiler.Prog.Memory.NewMailbox(lmc.Address(instr.ID()), "")
		toBox = op.Boxes[0].Box

		ops = append(ops, op)
	}

	compilation.Wrapped = instructions.NewWInstBitcast(instr, fromBox, toBox, ops)
	return &compilation
}

func (compiler *Compiler) WrapLLInstICmp(instr *ir.InstICmp, dstId lmc.Address) *Compilation {
	var compilation Compilation

	var xBox *lmc.Mailbox
	var yBox *lmc.Mailbox
	var dstBox *lmc.Mailbox
	var oneConst *lmc.Mailbox

	var ops []*lmc.MemoryOp
	var op *lmc.MemoryOp
	var err error

	if op, err = compiler.GetMailboxFromLL(instr.X); err != nil {
		compilation.Err = err
		return &compilation
	} else {
		xBox = op.Boxes[0].Box
		ops = append(ops, op)
	}

	if op, err = compiler.GetMailboxFromLL(instr.Y); err != nil {
		compilation.Err = err
		return &compilation
	} else {
		yBox = op.Boxes[0].Box
		ops = append(ops, op)
	}

	dstBox = compiler.Prog.Memory.GetMailboxAddress(dstId)
	if dstBox == nil {
		op := compiler.Prog.Memory.NewMailbox(dstId, "")
		dstBox = op.Boxes[0].Box

		ops = append(ops, op)
	}

	oneConst, err = compiler.Prog.Constant(1)
	if err != nil {
		compilation.Err = err
		return &compilation
	}

	compilation.Wrapped = instructions.NewWInstICmp(instr, xBox, yBox, dstBox, oneConst, ops)
	return &compilation
}

// ---------- Pattern matching instructions ----------

type Pattern interface {
	Match([]ir.Instruction) bool
	Find([]ir.Instruction) [][]int
	Compile([]ir.Instruction) *Compilation
	Priority() int
}

// ---------- singlePattern ----------

type singlePattern struct {
	matcher     func(ir.Instruction) bool
	wrapperFunc func(instr ir.Instruction) *Compilation
}

func (s *singlePattern) Match(i []ir.Instruction) bool {
	if len(i) != 1 {
		return false
	}

	return s.matcher(i[0])
}

func (s *singlePattern) Find(i []ir.Instruction) [][]int {
	var x [][]int

	for k := range i {
		if s.Match(i[k : k+1]) {
			x = append(x, []int{k})
		}
	}

	return x
}

func (s *singlePattern) Priority() int {
	return 0
}

func (s *singlePattern) Compile(i []ir.Instruction) *Compilation {
	if len(i) != 1 {
		panic("instructions given to compile a single-instr pattern is not of length 1")
	}

	return s.wrapperFunc(i[0])
}

func singleInstWrapper(compiler *Compiler, instr ir.Instruction) *Compilation {
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
	case *ir.InstBitCast:
		return compiler.WrapLLBitcast(cast)
	case *ir.InstICmp:
		return compiler.WrapLLInstICmp(cast, lmc.Address(cast.ID()))
	// unknown
	default:
		return &Compilation{Err: errors.E_UnknownLLInstruction(instr, nil)}
	}
}

// ---------- cmpZExtPattern ----------

type cmpZExtPattern struct {
	compiler *Compiler
}

func (c *cmpZExtPattern) Match(i []ir.Instruction) bool {
	if len(i) != 2 {
		return false
	}

	var ok bool
	var i1 *ir.InstICmp
	var i2 *ir.InstZExt

	i1, ok = i[0].(*ir.InstICmp)
	if ok {
		i2, ok = i[1].(*ir.InstZExt)
	}

	if !ok {
		return false
	}

	if i2ID, err := ReflectGetLocalID(i2.From); err != nil {
		panic(fmt.Sprintf("could not get local id via reflection from InstZExt: %s", err))
	} else {
		return lmc.Address(i1.ID()) == i2ID
	}
}

func (c *cmpZExtPattern) Find(i []ir.Instruction) [][]int {
	var x [][]int

	for j := 1; j < len(i); j++ {
		if c.Match(i[j-1 : j+1]) {
			x = append(x, []int{j - 1, j})
		}
	}

	return x
}

func (c *cmpZExtPattern) Compile(i []ir.Instruction) *Compilation {
	if len(i) != 2 {
		panic("instructions given to compiled cmp pattern is not of length 2")
	}

	i1 := i[0].(*ir.InstICmp)
	i2 := i[1].(*ir.InstZExt)

	return c.compiler.WrapLLInstICmp(i1, lmc.Address(i2.ID()))
}

func (c *cmpZExtPattern) Priority() int {
	return 10
}

func createPatterns(compiler *Compiler) []Pattern {

	// simple patterns - start

	simpleMatcherF := func(t reflect.Type) func(ir.Instruction) bool {
		return func(instr ir.Instruction) bool {
			return reflect.TypeOf(instr) == t
		}
	}

	simpleWrapperF := func(compiler *Compiler) func(ir.Instruction) *Compilation {
		return func(instr ir.Instruction) *Compilation {
			return singleInstWrapper(compiler, instr)
		}
	}

	var patterns []Pattern
	simpleWrapper := simpleWrapperF(compiler)

	for _, v := range []interface{}{
		&ir.InstAdd{}, &ir.InstSub{}, &ir.InstMul{}, &ir.InstSDiv{}, &ir.InstSRem{}, &ir.InstURem{}, &ir.InstAlloca{},
		&ir.InstLoad{}, &ir.InstStore{}, &ir.InstCall{}, &ir.InstBitCast{}, &ir.InstICmp{},
	} {
		patterns = append(patterns, &singlePattern{
			matcher:     simpleMatcherF(reflect.TypeOf(v)),
			wrapperFunc: simpleWrapper,
		})
	}

	// simple patterns - end

	// standalone patterns - start

	patterns = append(patterns, &cmpZExtPattern{compiler})

	// standalone patterns - end

	return patterns
}

// ---------- Engine ----------

type Match struct {
	Instrs     []ir.Instruction
	Pattern    Pattern
	firstInstr int
}

type Engine struct {
	compiler *Compiler
	patterns []Pattern
}

func NewEngine(compiler *Compiler) *Engine {
	var e Engine
	e.compiler = compiler

	for _, v := range createPatterns(compiler) {
		e.AddPattern(v)
	}

	return &e
}

func (e *Engine) AddPattern(p Pattern) {
	i := sort.Search(len(e.patterns), func(i int) bool {
		return e.patterns[i].Priority() >= p.Priority()
	})

	if i == len(e.patterns) {
		e.patterns = append(e.patterns, p)
	} else {
		e.patterns = append(e.patterns[:i+1], e.patterns[i:]...)
		e.patterns[i] = p
	}
}

func (e *Engine) FindAll(instrs []ir.Instruction) ([]*Match, error) {
	var err error
	c := NewOrderedSlice(0, func(a interface{}, b interface{}) bool {
		x := a.(Match)
		y := b.(Match)

		return x.firstInstr >= y.firstInstr
	})

	used := make(map[int]struct{})

	for _, p := range e.patterns {
		found := p.Find(instrs)

	foundLoop:
		for _, f := range found {
			var ii []ir.Instruction

			for _, ff := range f {
				if _, ok := used[ff]; ok {
					break foundLoop
				}

				ii = append(ii, instrs[ff])
				used[ff] = struct{}{}
			}

			c.Append(Match{
				Instrs:     ii,
				Pattern:    p,
				firstInstr: f[0],
			})
		}
	}

	for k := 0; k < len(instrs); k++ {
		if _, ok := used[k]; !ok {
			err = errors.E_UnknownLLInstruction(instrs[k], nil)
			break
		}
	}

	cc := make([]*Match, c.Len())
	for k, v := range c.Slice() {
		t := v.(Match)
		cc[k] = &t
	}

	return cc, err
}

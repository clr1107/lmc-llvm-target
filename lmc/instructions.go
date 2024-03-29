package lmc

import (
	"fmt"
	"strconv"
	"strings"
)

func formatInstrStr(name string, params []string) string {
	var buf strings.Builder
	buf.WriteString("Instr[" + name)

	if len(params) > 0 {
		buf.WriteRune('(')
		buf.WriteString(strings.Join(params, ","))
		buf.WriteRune(')')
	}

	buf.WriteRune(']')
	return buf.String()
}

// ---------- InstructionList ----------

// InstructionList holds all instructions along with, separately defined, data
// instructions for defining new mailboxes.
type InstructionList struct {
	Instructions    []Instruction
	DefInstructions []*DataInstr
}

func NewInstructionList() *InstructionList {
	return &InstructionList{
		Instructions:    make([]Instruction, 0),
		DefInstructions: make([]*DataInstr, 0),
	}
}

func (s *InstructionList) AddInstruction(instr Instruction) {
	s.Instructions = append(s.Instructions, instr)
}

func (s *InstructionList) RemoveInstruction(i int) error {
	if i >= len(s.Instructions) {
		return CannotRemoveInstructionIndexError(i, len(s.Instructions))
	}

	s.Instructions = append(s.Instructions[:i], s.Instructions[i+1:]...)
	return nil
}

func (s *InstructionList) AddDef(def *DataInstr) {
	s.DefInstructions = append(s.DefInstructions, def)
}

func (s *InstructionList) RemoveDef(identifier string) error {
	var i int
	var c int

	for _, x := range s.DefInstructions {
		if x.Box.Identifier() != identifier {
			s.DefInstructions[i] = x
			i++
		} else {
			c++
		}
	}

	if c != 0 {
		return VariableDoesNotExistError(identifier)
	} else {
		for j := i; j < len(s.DefInstructions); j++ {
			s.DefInstructions[j] = nil
		}

		s.DefInstructions = s.DefInstructions[:i]
	}

	return nil
}

// Implements LMCString by returning, as a string, the LMC instructions as a
// program.
func (s *InstructionList) LMCString() string {
	var buf strings.Builder

	// the length(longest label) + 1 is how many spaces to put before each instruction

	var longest int
	for _, v := range s.Instructions {
		if c, ok := v.(*Labelled); ok {
			if len(c.Identifier()) > longest {
				longest = len(c.Identifier())
			}
		}
	}

	if longest == 0 {
		longest = -1
	}

	for _, v := range s.Instructions {
		if c, ok := v.(*Labelled); ok {
			_, _ = fmt.Fprintf(
				&buf,
				"%s%s%s\n",
				c.Identifier(),
				strings.Repeat(" ", longest+1-len(c.Identifier())),
				c.Instruction.LMCString(),
			)
		} else {
			_, _ = fmt.Fprintf(
				&buf,
				"%s%s\n",
				strings.Repeat(" ", longest+1),
				v.LMCString(),
			)
		}
	}

	if len(s.DefInstructions) > 0 {
		buf.WriteRune('\n')

		for _, v := range s.DefInstructions {
			_, _ = fmt.Fprintf(&buf, "%s\n", v.LMCString())
		}
	}

	return buf.String()
}

func (s *InstructionList) String() string {
	return fmt.Sprintf("InstructionList[%d,%d]", len(s.Instructions), len(s.DefInstructions))
}

// ---------- Instructions base ----------

// Self-explanatory, eh?
type Instruction interface {
	LMCType
	Name() string
	Boxes() []*Mailbox
	ACC() bool
}

type InstructionBase struct {
	Instruction
	name string
}

func (i *InstructionBase) Name() string {
	return i.name
}

func (i *InstructionBase) Boxes() []*Mailbox {
	return make([]*Mailbox, 0)
}

// ---------- Data instruction ----------

// DataInstr handles the LMC data defining instruction `DAT`.
type DataInstr struct {
	InstructionBase
	Data Value
	Box  *Mailbox
}

func NewDataInstr(data Value, box *Mailbox) *DataInstr {
	return &DataInstr{
		InstructionBase: InstructionBase{
			name: "Data",
		},
		Data: data,
		Box:  box,
	}
}

func (i *DataInstr) String() string {
	return formatInstrStr(i.Name(), []string{strconv.Itoa(int(i.Data)), i.Box.Identifier()})
}

func (i *DataInstr) ACC() bool {
	return false
}

func (i *DataInstr) Boxes() []*Mailbox {
	return []*Mailbox{i.Box}
}

// Output form: `X DAT Y` where `X` is the box to be defined, and `Y` is the
// initial value, usually 0.
//
// E.g., `A DAT 0`.
func (i *DataInstr) LMCString() string {
	return fmt.Sprintf("%s DAT %d", i.Box.Identifier(), i.Data)
}

// ---------- Labelled instruction ----------

// Labelled merely holds an instruction and a label to tag to it. LMC is simple
// like that :)
type Labelled struct {
	label *Label
	Instruction
}

func NewLabelled(label *Label, instr Instruction) *Labelled {
	return &Labelled{
		label:       label,
		Instruction: instr,
	}
}

func (m *Labelled) Identifier() string {
	return m.label.Identifier()
}

// Output form: `X Y` where `X` is the identifier and `Y` is the lmc string form
// of the underlying instruction being labelled.
//
// E.g., `l_E ADD D`.
func (m *Labelled) LMCString() string {
	return fmt.Sprintf("%s %s", m.Identifier(), m.Instruction.LMCString())
}

// Depends on what it's wrapping.
func (m *Labelled) ACC() bool {
	return m.Instruction.ACC()
}

// ---------- Branch instruction ----------

type BranchType uint

// Branching types as an enumeration.
const (
	BRAlways   BranchType = iota // Branch Always
	BRPositive                   // Branch if acc +ve
	BRZero                       // Branch if acc zero
)

var (
	bMnemonics = [...]string{"BRA", "BRP", "BRZ"}
)

// BranchInstr handles the LMC concept of branching to labels.
type BranchInstr struct {
	InstructionBase
	BranchType BranchType
	label      *Label
}

func NewBranchInstr(branchType BranchType, label *Label) *BranchInstr {
	return &BranchInstr{
		BranchType: branchType,
		label:      label,
	}
}

func (b *BranchInstr) Identifier() string {
	return b.label.Identifier()
}

// Output form: `BR[A|P|Z] X` where `X` is the label to branch to.
//
// E.g., `BRZ l_A`.
func (b *BranchInstr) LMCString() string {
	return fmt.Sprintf("%s %s", bMnemonics[b.BranchType], b.label.Identifier())
}

func (b *BranchInstr) ACC() bool {
	return false
}

// ---------- Nullary instruction ----------

// NullaryInstr is a base struct for all instructions that have no parameters.
type NullaryInstr struct {
	InstructionBase
	mnemonic string
}

func (i *NullaryInstr) String() string {
	return formatInstrStr(i.Name(), nil)
}

// All nullary instructions have the output form: `X` where `X` is the mnemonic
// of the instruction.
//
// E.g., `INP`.
func (i *NullaryInstr) LMCString() string {
	return i.mnemonic
}

// ---------- Input instruction ----------

// InputInstr handles the LMC nullary instruction `INP`.
type InputInstr struct {
	NullaryInstr
}

func NewInputInstr() *InputInstr {
	return &InputInstr{
		NullaryInstr: NullaryInstr{
			InstructionBase: InstructionBase{
				name: "Input",
			},
			mnemonic: "INP",
		},
	}
}

func (i *InputInstr) ACC() bool {
	return true
}

// ---------- Output instruction ----------

// OutputInstr handles the nullary LMC instruction `OUT`.
type OutputInstr struct {
	NullaryInstr
}

func NewOutputInstr() *OutputInstr {
	return &OutputInstr{
		NullaryInstr: NullaryInstr{
			InstructionBase: InstructionBase{
				name: "Output",
			},
			mnemonic: "OUT",
		},
	}
}

func (o *OutputInstr) ACC() bool {
	return false
}

// ---------- Halt instruction ----------

// HaltInstr handles the nullary LMC instruction `HLT`.
type HaltInstr struct {
	NullaryInstr
}

func NewHaltInstr() *HaltInstr {
	return &HaltInstr{
		NullaryInstr: NullaryInstr{
			InstructionBase: InstructionBase{
				name: "Halt",
			},
			mnemonic: "HLT",
		},
	}
}

func (h *HaltInstr) ACC() bool {
	return false
}

// ---------- Unary instruction ----------

// UnaryInstr is a base struct for all instructions that have one parameter.
type UnaryInstr struct {
	InstructionBase
	Param    *Mailbox
	mnemonic string
}

func (i *UnaryInstr) Boxes() []*Mailbox {
	return []*Mailbox{i.Param}
}

func (i *UnaryInstr) String() string {
	return formatInstrStr(i.Name(), []string{i.Param.Identifier()})
}

// All unary instructions have the output form: `X Y` where `X` is the mnemonic
// of the instruction and `Y` is the parameter.
//
// E.g., `LDA A`.
func (i *UnaryInstr) LMCString() string {
	return fmt.Sprintf("%s %s", i.mnemonic, i.Param.Identifier())
}

// ---------- Add instruction ----------

// AddInstr handles the unary LMC instruction `ADD`.
type AddInstr struct {
	UnaryInstr
}

func NewAddInstr(param *Mailbox) *AddInstr {
	return &AddInstr{
		UnaryInstr: UnaryInstr{
			InstructionBase: InstructionBase{
				name: "Add",
			},
			Param:    param,
			mnemonic: "ADD",
		},
	}
}

func (a *AddInstr) ACC() bool {
	return true
}

// ---------- Subtract instruction ----------

// SubInstr handles the unary LMC instruction `SUB`.
type SubInstr struct {
	UnaryInstr
}

func NewSubInstr(param *Mailbox) *SubInstr {
	return &SubInstr{
		UnaryInstr: UnaryInstr{
			InstructionBase: InstructionBase{
				name: "Sub",
			},
			Param:    param,
			mnemonic: "SUB",
		},
	}
}

func (s *SubInstr) ACC() bool {
	return true
}

// ---------- Store instruction ----------

// StoreInstr handles the unary LMC instruction `STA`.
type StoreInstr struct {
	UnaryInstr
}

func NewStoreInstr(param *Mailbox) *StoreInstr {
	return &StoreInstr{
		UnaryInstr: UnaryInstr{
			InstructionBase: InstructionBase{
				name: "Store",
			},
			Param:    param,
			mnemonic: "STA",
		},
	}
}

func (s *StoreInstr) ACC() bool {
	return false
}

// ---------- Load instruction ----------

// LoadInstr handles the unary LMC instruction `LDA`.
type LoadInstr struct {
	UnaryInstr
}

func NewLoadInstr(param *Mailbox) *LoadInstr {
	return &LoadInstr{
		UnaryInstr: UnaryInstr{
			InstructionBase: InstructionBase{
				name: "Load",
			},
			Param:    param,
			mnemonic: "LDA",
		},
	}
}

func (l *LoadInstr) ACC() bool {
	return true
}

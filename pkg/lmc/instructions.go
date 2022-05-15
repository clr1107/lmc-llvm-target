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

// ---------- InstructionSet ----------

// InstructionSet holds all instructions along with, separately defined, data
// instructions for defining new mailboxes.
type InstructionSet struct {
	instructions    []Instruction
	defInstructions []*DataInstr
}

func NewInstructionSet() *InstructionSet {
	return &InstructionSet{
		instructions:    make([]Instruction, 0),
		defInstructions: make([]*DataInstr, 0),
	}
}

func (s *InstructionSet) GetInstructions() []Instruction {
	return s.instructions
}

func (s *InstructionSet) GetDefs() []*DataInstr {
	return s.defInstructions
}

func (s *InstructionSet) AddInstruction(instr Instruction) {
	s.instructions = append(s.instructions, instr)
}

func (s *InstructionSet) RemoveInstruction(i int) error {
	if i >= len(s.instructions) {
		return fmt.Errorf(
			"attempted to remove instruction with index %d out of %d instructions",
			i, len(s.instructions),
		)
	}

	s.instructions = append(s.instructions[:i], s.instructions[i+1:]...)
	return nil
}

func (s *InstructionSet) AddDef(def *DataInstr) {
	s.defInstructions = append(s.defInstructions, def)
}

func (s *InstructionSet) RemoveDef(name string) error {
	var i []int
	for k, v := range s.defInstructions {
		if v.name == name {
			i = append(i, k)
			break
		}
	}

	if len(i) == 0 {
		for _, ii := range i {
			s.defInstructions = append(s.defInstructions[:ii], s.defInstructions[ii+1:]...)
		}
	} else {
		return fmt.Errorf("attempted to remove all def of variable `%s` when it does not exist", name)
	}

	return nil
}

// Implements LMCString by returning, as a string, the LMC instructions as a
// program.
func (s *InstructionSet) LMCString() string {
	var buf strings.Builder

	// the length(longest label) + 1 is how many spaces to put before each instruction

	var longest int
	for _, v := range s.instructions {
		if c, ok := v.(*Labelled); ok {
			if len(c.Identifier()) > longest {
				longest = len(c.Identifier())
			}
		}
	}

	if longest == 0 {
		longest = -1
	}

	for _, v := range s.instructions {
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

	if len(s.defInstructions) > 0 {
		buf.WriteRune('\n')

		for _, v := range s.defInstructions {
			_, _ = fmt.Fprintf(&buf, "%s\n", v.LMCString())
		}
	}

	return buf.String()
}

func (s *InstructionSet) String() string {
	return fmt.Sprintf("InstructionSet[%d,%d]", len(s.instructions), len(s.defInstructions))
}

// ---------- Instructions base ----------

// Self-explanatory, eh?
type Instruction interface {
	LMCType
	Name() string
}

type InstructionBase struct {
	Instruction
	name string
}

func (i *InstructionBase) Name() string {
	return i.name
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

// ---------- Unary instruction ----------

// UnaryInstr is a base struct for all instructions that have one parameter.
type UnaryInstr struct {
	InstructionBase
	Param    *Mailbox
	mnemonic string
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

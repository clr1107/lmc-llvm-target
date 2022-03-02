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

type InstructionSet struct {
	instructions []IInstruction
	defInstructions []*DataInstr
}

func NewInstructionSet() *InstructionSet {
	return &InstructionSet{
		instructions: make([]IInstruction, 0),
		defInstructions: make([]*DataInstr, 0),
	}
}

func (s *InstructionSet) AddInstruction(instr IInstruction) {
	s.instructions = append(s.instructions, instr)
}

func (s *InstructionSet) AddDef(def *DataInstr) {
	s.defInstructions = append(s.defInstructions, def)
}

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

	for _, v := range s.instructions {
		if c, ok := v.(*Labelled); ok {
			_, _ = fmt.Fprintf(
				&buf,
				"%s%s%s\n",
				c.Identifier(),
				strings.Repeat(" ",  longest + 1 - len(c.Identifier())),
				c.IInstruction.LMCString(),
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

type IInstruction interface {
	LMCType
	Name() string
}

type InstructionBase struct {
	IInstruction
	name string
}

func (i *InstructionBase) Name() string {
	return i.name
}

// ---------- Unary instruction ----------

type UnaryInstr struct {
	InstructionBase
	mnemonic string
}

func (i *UnaryInstr) String() string {
	return formatInstrStr(i.Name(), nil)
}

func (i *UnaryInstr) LMCString() string {
	return i.mnemonic
}

// ---------- Input instruction ----------

type InputInstr struct {
	UnaryInstr
}

func NewInputInstr() *InputInstr {
	return &InputInstr{
		UnaryInstr: UnaryInstr{
			InstructionBase: InstructionBase{
				name: "Input",
			},
			mnemonic: "INP",
		},
	}
}

// ---------- Output instruction ----------

type OutputInstr struct {
	UnaryInstr
}

func NewOutputInstr() *OutputInstr {
	return &OutputInstr{
		UnaryInstr: UnaryInstr{
			InstructionBase: InstructionBase{
				name: "Output",
			},
			mnemonic: "OUT",
		},
	}
}

// ---------- Halt instruction ----------

type HaltInstr struct {
	UnaryInstr
}

func NewHaltInstr() *HaltInstr {
	return &HaltInstr{
		UnaryInstr: UnaryInstr{
			InstructionBase: InstructionBase{
				name: "Halt",
			},
			mnemonic: "HLT",
		},
	}
}

// ---------- Binary instruction ----------

type BinaryInstr struct {
	InstructionBase
	Param    *Mailbox
	mnemonic string
}

func (i *BinaryInstr) String() string {
	return formatInstrStr(i.Name(), []string{i.Param.Identifier()})
}

func (i *BinaryInstr) LMCString() string {
	return fmt.Sprintf("%s %s", i.mnemonic, i.Param.Identifier())
}

// ---------- Add instruction ----------

type AddInstr struct {
	BinaryInstr
}

func NewAddInstr(param *Mailbox) *AddInstr {
	return &AddInstr{
		BinaryInstr: BinaryInstr{
			InstructionBase: InstructionBase{
				name: "Add",
			},
			Param:    param,
			mnemonic: "ADD",
		},
	}
}

// ---------- Subtract instruction ----------

type SubInstr struct {
	BinaryInstr
}

func NewSubInstr(param *Mailbox) *SubInstr {
	return &SubInstr{
		BinaryInstr: BinaryInstr{
			InstructionBase: InstructionBase{
				name: "Sub",
			},
			Param:    param,
			mnemonic: "SUB",
		},
	}
}

// ---------- Store instruction ----------

type StoreInstr struct {
	BinaryInstr
}

func NewStoreInstr(param *Mailbox) *StoreInstr {
	return &StoreInstr{
		BinaryInstr: BinaryInstr{
			InstructionBase: InstructionBase{
				name: "Store",
			},
			Param:    param,
			mnemonic: "STA",
		},
	}
}

// ---------- Load instruction ----------

type LoadInstr struct {
	BinaryInstr
}

func NewLoadInstr(param *Mailbox) *LoadInstr {
	return &LoadInstr{
		BinaryInstr: BinaryInstr{
			InstructionBase: InstructionBase{
				name: "Load",
			},
			Param:    param,
			mnemonic: "LDA",
		},
	}
}

// ---------- Data instruction ----------

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

func (i *DataInstr) LMCString() string {
	return fmt.Sprintf("%s DAT %d", i.Box.Identifier(), i.Data)
}

// ---------- Labelled instruction ----------

type Labelled struct {
	label *Label
	IInstruction
}

func NewLabelled(label *Label, instr IInstruction) *Labelled {
	return &Labelled{
		label: label,
		IInstruction: instr,
	}
}

func (m *Labelled) Identifier() string {
	return m.label.Identifier()
}

func (m *Labelled) LMCString() string {
	return fmt.Sprintf("%s %s", m.Identifier(), m.IInstruction.LMCString())
}

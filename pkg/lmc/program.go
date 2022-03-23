package lmc

type Program struct {
	Memory       *Memory
	Instructions *InstructionSet
	constants map[Value]*Mailbox
}

func NewProgram(memory *Memory) *Program {
	return &Program{
		Memory:       memory,
		Instructions: NewInstructionSet(),
		constants: make(map[Value]*Mailbox, 0),
	}
}

func (p *Program) AddInstructions(instrs []IInstruction, defs []*DataInstr) error {
	for _, v := range instrs {
		p.Instructions.AddInstruction(v)
	}
	for _, v := range defs {
		p.Instructions.AddDef(v)
	}

	return nil
}

func (p *Program) NewMailbox(addr Address, identifier string) (*Mailbox, error) {
	mbox := p.Memory.NewMailbox(addr, identifier)
	if err := p.Memory.AddMailbox(mbox); err != nil {
		return mbox, err
	} else {
		def := NewDataInstr(0, mbox)
		return mbox, p.AddInstructions(nil, []*DataInstr{def})
	}
}

func (p *Program) NewLabel(identifier string) (*Label, error) {
	label := p.Memory.NewLabel(identifier)
	return label, p.Memory.AddLabel(label)
}

func (p *Program) Constant(value Value) (*Mailbox, error) {
	if v, ok := p.constants[value]; ok {
		return v, nil
	} else {
		identifier := "c_" + p.Memory.idGen(len(p.constants))
		mbox := p.Memory.NewMailbox(-1, identifier)

		if err := p.Memory.AddMailbox(mbox); err != nil {
			return mbox, err
		} else {
			p.constants[value] = mbox
		}

		def := NewDataInstr(value, mbox)
		return mbox, p.AddInstructions(nil, []*DataInstr{def})
	}
}

func (p *Program) String() string {
	return p.Instructions.LMCString()
}

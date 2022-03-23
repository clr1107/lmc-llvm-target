package lmc

type Program struct {
	Memory       *Memory
	constants map[Value]*Mailbox
}

func NewProgram(memory *Memory) *Program {
	return &Program{
		Memory:       memory,
		constants: make(map[Value]*Mailbox, 0),
	}
}

func (p *Program) AddInstructions(instrs []IInstruction, defs []*DataInstr) {
	for _, v := range instrs {
		p.Memory.instructions.AddInstruction(v)
	}
	for _, v := range defs {
		p.Memory.instructions.AddDef(v)
	}
}

func (p *Program) NewMailbox(addr Address, identifier string) (*Mailbox, error) {
	mbox := p.Memory.NewMailbox(addr, identifier)
	if err := p.Memory.AddMailbox(mbox); err != nil {
		return mbox, err
	} else {
		def := NewDataInstr(0, mbox)
		p.AddInstructions(nil, []*DataInstr{def})

		return mbox, nil
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
		p.AddInstructions(nil, []*DataInstr{def})

		return mbox, nil
	}
}

func (p *Program) String() string {
	return p.Memory.instructions.LMCString()
}

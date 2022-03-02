package lmc

type Program struct {
	Memory       *Memory
	Instructions *InstructionSet
}

func NewProgram() *Program {
	return &Program{
		Memory:       NewMemory(),
		Instructions: NewInstructionSet(),
	}
}

func (p *Program) AddInstructions(instrs []IInstruction, defs []*DataInstr) error {
	if len(p.Memory.mailboxes)+len(instrs)+len(defs) > 100 {
		return OutOfSpaceError
	}

	for _, v := range instrs {
		p.Instructions.AddInstruction(v)
	}
	for _, v := range defs {
		p.Instructions.AddDef(v)
	}

	return nil
}

func (p *Program) NewMailbox(addr Address, identifier string) (*Mailbox, error) {
	if len(p.Memory.mailboxes) + 2 > 100 { // quick check, first.
		return nil, OutOfSpaceError
	}

	if mbox, err := p.Memory.NewMailbox(addr, identifier); err != nil {
		return nil, err
	} else {
		def := NewDataInstr(0, mbox)
		return mbox, p.AddInstructions(nil, []*DataInstr{def})
	}
}

func (p *Program) String() string {
	return p.Instructions.LMCString()
}

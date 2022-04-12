package lmc

type Program struct {
	Memory       *Memory
}

func NewProgram(memory *Memory) *Program {
	return &Program{
		Memory:       memory,
	}
}

func (p *Program) AddInstructions(instrs []Instruction, defs []*DataInstr) {
	for _, v := range instrs {
		p.Memory.instructions.AddInstruction(v)
	}
	for _, v := range defs {
		p.Memory.instructions.AddDef(v)
	}
}

func (p *Program) NewMailbox(addr Address, identifier string) (*Mailbox, error) {
	op := p.Memory.NewMailbox(addr, identifier)
	if err := p.Memory.AddMailbox(op.GetNew()[0]); err != nil {
		return nil, err
	}

	p.AddInstructions(nil, op.Defs())
	return op.Boxes[0].Box, nil
}

func (p *Program) NewLabel(identifier string) (*Label, error) {
	label := p.Memory.NewLabel(identifier)
	return label, p.Memory.AddLabel(label)
}

func (p *Program) Constant(value Value) (*Mailbox, error) {
	op := p.Memory.Constant(value)

	n := op.GetNew()
	if len(n) > 0 {
		if err := p.Memory.AddMailbox(n[0]); err != nil {
			return nil, err
		}

		p.AddInstructions(nil, op.Defs())
		p.Memory.constants[value] = op.Boxes[0].Box
	}

	return op.Boxes[0].Box, nil
}

func (p *Program) String() string {
	return p.Memory.instructions.LMCString()
}

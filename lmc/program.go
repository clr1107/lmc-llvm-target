package lmc

// Program represents an LMC program. It holds the mailboxes and provides
// utility functions to operate on the memory.
type Program struct {
	Memory *Memory
}

func NewProgram(memory *Memory) *Program {
	return &Program{
		Memory: memory,
	}
}

// AddInstructions adds the instructions passed. Slices can be nil to add
// nothing.
func (p *Program) AddInstructions(instrs []Instruction, defs []*DataInstr) {
	for _, v := range instrs {
		p.Memory.instructions.AddInstruction(v)
	}
	for _, v := range defs {
		p.Memory.instructions.AddDef(v)
	}
}

// AddMemoryOp adds all the mailboxes, labels, and data instructions from a
// memory operation.
func (p *Program) AddMemoryOp(op *MemoryOp) error {
	for _, box := range op.GetNewBoxes() {
		if err := p.Memory.AddMailbox(box); err != nil {
			return err
		}
	}

	for _, label := range op.GetNewLabels() {
		if err := p.Memory.AddLabel(label); err != nil {
			return err
		}
	}

	p.AddInstructions(nil, op.Defs())
	return nil
}

// NewMailbox creates a new mailbox with the given address and identifier. If
// the identifier is empty one will be generated. This function will also
// attempt to add the mailbox and def instruction to memory; returning an error
// if this fails. See *Memory#NewMailbox.
//
// [Memory utility function]
func (p *Program) NewMailbox(addr Address, identifier string) (*Mailbox, error) {
	op := p.Memory.NewMailbox(addr, identifier)
	if err := p.Memory.AddMailbox(op.Boxes[0].Box); err != nil {
		return nil, err
	}

	p.AddInstructions(nil, op.Defs())
	return op.Boxes[0].Box, nil
}

// NewLabel creates a new label with the given identifier. If the identifier is
// empty one will be generated. This function will also attempt to add the label
// to memory; returning an error if this fails. See *Memory#NewLabel.
//
// [Memory utility function]
func (p *Program) NewLabel(identifier string) (*Label, error) {
	op := p.Memory.NewLabel(identifier)
	if err := p.Memory.AddLabel(op.Labels[0].Label); err != nil {
		return nil, err
	}

	return op.Labels[0].Label, nil
}

// Constant returns a mailbox storing a constant value. If one does not already
// exist it will be created, and the box and instructions will be added to
// memory, returning an error if this fails. See *Memory#Constant.
//
// [Memory utility function]
func (p *Program) Constant(value Value) (*Mailbox, error) {
	op := p.Memory.Constant(value)
	box := op.Boxes[0]

	defs := op.Defs()
	if len(defs) > 0 {
		if err := p.Memory.AddMailbox(box.Box); err != nil {
			return nil, err
		}

		p.AddInstructions(nil, defs)
	}

	return box.Box, nil
}

func (p *Program) String() string {
	return p.Memory.instructions.LMCString()
}

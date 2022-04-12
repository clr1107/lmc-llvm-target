package lmc

import (
	"fmt"
	"strings"
)

var (
	MailboxAlreadyExistsAddressError = func(addr Address) error {
		return fmt.Errorf("a mailbox with address %d already exists", addr)
	}
	MailboxAlreadyExistsIdentifierError = func(identifier string) error {
		return fmt.Errorf("a mailbox with identifier `%s' already exists", identifier)
	}
	LabelAlreadyExistsError = func(identifier string) error {
		return fmt.Errorf("a label with identifier `%s' already exists", identifier)
	}
)

// ---------- Mailbox ----------

type Mailbox struct {
	addr       Address
	identifier string
}

func NewMailbox(addr Address, identifier string) *Mailbox {
	return &Mailbox{
		addr:       addr,
		identifier: identifier,
	}
}

func (m *Mailbox) Identifier() string {
	return m.identifier
}

func (m *Mailbox) Address() Address {
	return m.addr
}

// --------- Label ----------

type Label struct {
	LMCType
	identifier string
}

func NewLabel(identifier string) *Label {
	return &Label{
		identifier: identifier,
	}
}

func (l *Label) Identifier() string {
	return l.identifier
}

func (l *Label) String() string {
	return fmt.Sprintf("Label[%s]", l.Identifier())
}

func (l *Label) LMCString() string {
	return l.Identifier()
}

func makeIdentifierGenerator() func(int) string {
	identifierSymbols := [...]rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R',
		'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	}

	return func(i int) string {
		buf := &strings.Builder{}

		for {
			s := buf.String()
			buf.Reset()

			buf.WriteRune(identifierSymbols[i%26])
			buf.WriteString(s)

			if i /= 26; i == 0 {
				break
			}
		}

		return buf.String()
	}
}

// ---------- MemoryOp ----------

type MemoryOpPair struct {
	Box *Mailbox
	New bool
	Value Value
}

type MemoryOp struct {
	Boxes []*MemoryOpPair
}

func NewMemoryOp1(box *Mailbox, new bool) *MemoryOp {
	return &MemoryOp{
		Boxes: []*MemoryOpPair{
			{Box: box, New: new},
		},
	}
}

func NewMemoryOp2(box1 *Mailbox, new1 bool, box2 *Mailbox, new2 bool) *MemoryOp {
	return &MemoryOp{
		Boxes: []*MemoryOpPair{
			{Box: box1, New: new1},
			{Box: box2, New: new2},
		},
	}
}

func NewMemoryOp3(box1 *Mailbox, new1 bool, box2 *Mailbox, new2 bool, box3 *Mailbox, new3 bool) *MemoryOp {
	return &MemoryOp{
		Boxes: []*MemoryOpPair{
			{Box: box1, New: new1},
			{Box: box2, New: new2},
			{Box: box3, New: new3},
		},
	}
}

func (m *MemoryOp) GetNew() []*Mailbox {
	var l []*Mailbox
	for _, p := range m.Boxes {
		if p.New {
			l = append(l, p.Box)
		}
	}

	return l
}

func (m *MemoryOp) Defs() []*DataInstr {
	var l []*DataInstr
	for _, p := range m.Boxes {
		if p.New {
			l = append(l, NewDataInstr(p.Value, p.Box))
		}
	}

	return l
}

// ---------- Memory ----------

type Memory struct {
	mailboxes    []*Mailbox
	instructions *InstructionSet
	labels       []*Label
	constants    map[Value]*Mailbox
	idGen        func(int) string
}

func NewMemory(idGen func(int) string) *Memory {
	return &Memory{
		mailboxes:    make([]*Mailbox, 0),
		instructions: NewInstructionSet(),
		labels:       make([]*Label, 0),
		constants:    make(map[Value]*Mailbox, 0),
		idGen:        idGen,
	}
}

func NewBasicMemory() *Memory {
	return NewMemory(makeIdentifierGenerator())
}

func (m *Memory) GetInstructions() *InstructionSet {
	return m.instructions
}

func (m *Memory) GetMailboxAddress(addr Address) *Mailbox {
	if addr >= 0 {
		for _, v := range m.mailboxes {
			if v.Address() == addr {
				return v
			}
		}
	}

	return nil
}

func (m *Memory) GetMailboxIdentifier(identifier string) *Mailbox {
	for _, v := range m.mailboxes {
		if v.Identifier() == identifier {
			return v
		}
	}

	return nil
}

func (m *Memory) GetLabel(identifier string) *Label {
	for _, v := range m.labels {
		if v.Identifier() == identifier {
			return v
		}
	}

	return nil
}

func (m *Memory) AddMailbox(mailbox *Mailbox) error {
	if m.GetMailboxAddress(mailbox.Address()) != nil {
		return MailboxAlreadyExistsAddressError(mailbox.Address())
	}

	if m.GetMailboxIdentifier(mailbox.Identifier()) != nil {
		return MailboxAlreadyExistsIdentifierError(mailbox.Identifier())
	}

	m.mailboxes = append(m.mailboxes, mailbox)

	return nil
}

func (m *Memory) AddLabel(label *Label) error {
	if m.GetLabel(label.Identifier()) != nil {
		return LabelAlreadyExistsError(label.Identifier())
	}

	m.labels = append(m.labels, label)
	return nil
}

func (m *Memory) NewMailbox(addr Address, identifier string) *MemoryOp {
	if identifier == "" {
		identifier = m.idGen(int(addr))
	}

	box := NewMailbox(addr, identifier)
	return NewMemoryOp1(box, true)
}

func (m *Memory) NewLabel(identifier string) *Label {
	if identifier == "" {
		identifier = "l_" + m.idGen(len(m.labels))
	}

	label := NewLabel(identifier)
	return label
}

func (m *Memory) Constant(value Value) *MemoryOp {
	if v, ok := m.constants[value]; ok {
		return NewMemoryOp1(v, false)
	} else {
		identifier := "c_" + m.idGen(len(m.constants))

		op := m.NewMailbox(-1, identifier)
		op.Boxes[0].Value = value

		m.constants[value] = op.Boxes[0].Box
		return op
	}
}

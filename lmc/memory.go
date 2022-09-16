package lmc

import (
	"fmt"
	"strings"
)

// ---------- Mailbox ----------

// Mailbox represents one memory location in LMC. It has an address, used only
// for internal logic, and an identifier. A -ve address is for non-user created
// boxes, whatever that means for the application it's used in.
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

// Label allows for a label for an instruction in LMC. Labels have only an
// identifier, and are not stored in memory -- they are read at runtime.
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

// MemoryOpBoxPair is a pair of a mailbox and a flag for if it is a new box.
// It also contains a value, for if this is an initialised box (i.e., non-zero
// initial value).
type MemoryOpBoxPair struct {
	Box   *Mailbox
	New   bool
	Value Value
}

// MemoryOpLabelPair is a pair of a label and a flag for if it is a new label.
type MemoryOpLabelPair struct {
	Label *Label
	New   bool
}

// MemoryOp contains all boxes and labels returned from completing a memory
// operation.
type MemoryOp struct {
	Boxes  []*MemoryOpBoxPair
	Labels []*MemoryOpLabelPair
}

// NewMemoryOpBox1 creates a memory operation with one mailbox.
func NewMemoryOpBox1(box *Mailbox, new bool) *MemoryOp {
	return &MemoryOp{
		Boxes: []*MemoryOpBoxPair{
			{Box: box, New: new},
		},
	}
}

// NewMemoryOpBox2 creates a memory operation with two mailboxes.
func NewMemoryOpBox2(box1 *Mailbox, new1 bool, box2 *Mailbox, new2 bool) *MemoryOp {
	return &MemoryOp{
		Boxes: []*MemoryOpBoxPair{
			{Box: box1, New: new1},
			{Box: box2, New: new2},
		},
	}
}

// NewMemoryOpBox3 creates a memory operation with three mailboxes.
func NewMemoryOpBox3(box1 *Mailbox, new1 bool, box2 *Mailbox, new2 bool, box3 *Mailbox, new3 bool) *MemoryOp {
	return &MemoryOp{
		Boxes: []*MemoryOpBoxPair{
			{Box: box1, New: new1},
			{Box: box2, New: new2},
			{Box: box3, New: new3},
		},
	}
}

// NewMemoryOpLabel1 creates a memory operation with one label.
func NewMemoryOpLabel1(label *Label, new bool) *MemoryOp {
	return &MemoryOp{
		Labels: []*MemoryOpLabelPair{
			{Label: label, New: new},
		},
	}
}

// NewMemoryOpLabel2 creates a memory operation with two labels.
func NewMemoryOpLabel2(label1 *Label, new1 bool, label2 *Label, new2 bool) *MemoryOp {
	return &MemoryOp{
		Labels: []*MemoryOpLabelPair{
			{Label: label1, New: new1},
			{Label: label2, New: new2},
		},
	}
}

// NewMemoryOpLabel3 creates a memory operation with three labels.
func NewMemoryOpLabel3(label1 *Label, new1 bool, label2 *Label, new2 bool, label3 *Label, new3 bool) *MemoryOp {
	return &MemoryOp{
		Labels: []*MemoryOpLabelPair{
			{Label: label1, New: new1},
			{Label: label2, New: new2},
			{Label: label3, New: new3},
		},
	}
}

// GetNewBoxes gives all mailboxes flagged as being new, in a new slice.
func (m *MemoryOp) GetNewBoxes() []*Mailbox {
	var l []*Mailbox
	for _, p := range m.Boxes {
		if p.New {
			l = append(l, p.Box)
		}
	}

	return l
}

// GetNewLabels gives all labels flagged as being new, in a new slice.
func (m *MemoryOp) GetNewLabels() []*Label {
	var l []*Label
	for _, p := range m.Labels {
		if p.New {
			l = append(l, p.Label)
		}
	}

	return l
}

// Defs creates a new slice with data instructions for all new boxes, with their
// values, if initialised.
func (m *MemoryOp) Defs() []*DataInstr {
	var l []*DataInstr
	for _, p := range m.Boxes {
		if p.New {
			l = append(l, NewDataInstr(p.Value, p.Box))
		}
	}

	return l
}

func (m *MemoryOp) AddNewBoxOps(op ...*MemoryOpBoxPair) {
	m.Boxes = append(m.Boxes, op...)
}

func (m *MemoryOp) AddNewLabelOps(op ...*MemoryOpBoxPair) {
	m.Boxes = append(m.Boxes, op...)
}

// ---------- Memory ----------

// Memory handles all mailboxes (including constants), instructions, and labels.
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

func (m *Memory) GetInstructionSet() *InstructionSet {
	return m.instructions
}

func (m *Memory) GetMailboxes() []*Mailbox {
	return m.mailboxes
}

// GetMailboxAddress returns the first mailbox to match the given address. Nil otherwise.
// Only handles addresses >= 0 as -ve addresses are non-user created boxes.
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

// GetMailboxIdentifier returns the first mailbox to match the given identifier.
// Case-sensitive, obviously... Nil otherwise.
func (m *Memory) GetMailboxIdentifier(identifier string) *Mailbox {
	for _, v := range m.mailboxes {
		if v.Identifier() == identifier {
			return v
		}
	}

	return nil
}

// GetLabel returns the first label to match the given identifier.
// Case-sensitive, obviously... Nil otherwise.
func (m *Memory) GetLabel(identifier string) *Label {
	for _, v := range m.labels {
		if v.Identifier() == identifier {
			return v
		}
	}

	return nil
}

// AddMailbox will try to add a given mailbox to the memory; returning an error
// if the address or identifier are already in use.
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

// RemoveMailboxIdentifier will remove all mailboxes with the given identifier.
func (m *Memory) RemoveMailboxIdentifier(identifier string) bool {
	var i, c int

	for _, b := range m.mailboxes {
		if b.Identifier() != identifier {
			m.mailboxes[i] = b
			i++
		} else {
			c++
		}
	}

	if c > 0 {
		for j := i; j < len(m.mailboxes); j++ {
			m.mailboxes[j] = nil
		}

		m.mailboxes = m.mailboxes[:i]
	}

	return c > 0
}

// RemoveMailboxAddress will remove all mailboxes with the given address.
func (m *Memory) RemoveMailboxAddress(address Address) bool {
	var i, c int

	for _, b := range m.mailboxes {
		if b.addr != address {
			m.mailboxes[i] = b
			i++
		} else {
			c++
		}
	}

	if c > 0 {
		for j := i; j < len(m.mailboxes); j++ {
			m.mailboxes[j] = nil
		}
	}

	return c > 0
}

// AddLabel will try to add a given label to the memory; returning an error if
// the identifier is already in use.
func (m *Memory) AddLabel(label *Label) error {
	if m.GetLabel(label.Identifier()) != nil {
		return LabelAlreadyExistsError(label.Identifier())
	}

	m.labels = append(m.labels, label)
	return nil
}

// NewMailbox creates a new mailbox with a given address and identifier. If the
// identifier is the empty string one is generated using the generator function.
//
// This returns a memory operation. See advisory note in overview.
func (m *Memory) NewMailbox(addr Address, identifier string) *MemoryOp {
	if identifier == "" {
		identifier = m.idGen(int(addr))
	}

	box := NewMailbox(addr, identifier)
	return NewMemoryOpBox1(box, true)
}

// NewLabel creates a new label with a given identifier. If the identifier is
// the empty string one is generated using the generator function.
//
// This returns a memory operation. See advisory note in overview.
func (m *Memory) NewLabel(identifier string) *MemoryOp {
	if identifier == "" {
		identifier = "l_" + m.idGen(len(m.labels))
	}

	label := NewLabel(identifier)
	return NewMemoryOpLabel1(label, true)
}

// Constant returns a mailbox with the value given for use in arithmetic etc.
// If one does not exist, it is created. All constants have an auto generated
// identifier, prefixed with 'c_'. The mailboxes have id -1.
//
// This returns a memory operation. See advisory note in overview.
func (m *Memory) Constant(value Value) *MemoryOp {
	if v, ok := m.constants[value]; ok {
		return NewMemoryOpBox1(v, false)
	} else {
		identifier := "c_" + m.idGen(len(m.constants))

		op := m.NewMailbox(-1, identifier)
		box := op.Boxes[0]

		box.Value = value
		m.constants[value] = box.Box

		return op
	}
}

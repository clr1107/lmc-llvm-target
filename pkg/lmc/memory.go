package lmc

import (
	"errors"
	"fmt"
	"strings"
)

var (
	OutOfSpaceError             = errors.New("out of space for more mailboxes")
	MailboxAlreadyExistsAddress = func(addr Address) error {
		return fmt.Errorf("a mailbox with address %d already exists", addr)
	}
	MailboxAlreadyExistsIdentifier = func(identifier string) error {
		return fmt.Errorf("a mailbox with identifier `%s' already exists", identifier)
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

// ---------- Memory ----------

func makeIdentifierGenerator() func(Address) string {
	identifierSymbols := [26]rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R',
		'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	}

	return func(a Address) string {
		buf := &strings.Builder{}

		for {
			s := buf.String()
			buf.Reset()

			buf.WriteRune(identifierSymbols[a%26])
			buf.WriteString(s)

			if a /= 26; a == 0 {
				break
			}
		}

		return buf.String()
	}
}

type Memory struct {
	mailboxes []*Mailbox
	count     int
	idGen     func(Address) string
}

func NewMemory() *Memory {
	return &Memory{
		mailboxes: make([]*Mailbox, 0),
		count:     0,
		idGen:     makeIdentifierGenerator(),
	}
}

func (m *Memory) GetMailboxAddress(addr Address) *Mailbox {
	for _, v := range m.mailboxes {
		if v.Address() == addr {
			return v
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

func (m *Memory) AddMailbox(mailbox *Mailbox) error {
	if m.count == 99 {
		return OutOfSpaceError
	}

	if m.GetMailboxAddress(mailbox.Address()) != nil {
		return MailboxAlreadyExistsAddress(mailbox.Address())
	}

	if m.GetMailboxIdentifier(mailbox.Identifier()) != nil {
		return MailboxAlreadyExistsIdentifier(mailbox.Identifier())
	}

	m.mailboxes = append(m.mailboxes, mailbox)
	m.count += 1

	return nil
}

func (m *Memory) NewMailbox(addr Address) (*Mailbox, error) {
	identifier := m.idGen(addr)
	mailbox := NewMailbox(addr, identifier)

	return mailbox, m.AddMailbox(mailbox)
}

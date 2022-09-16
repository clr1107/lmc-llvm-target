package optimisation

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/lmc"
)

var propStageNames = [...]string{
	"PROP_TREE",
	"PROP_LDA_STA",
}

func propErr(stage int, child error) error {
	return fmt.Errorf("box propogation failed stage %d=%s: %s", stage, propStageNames[stage], child)
}

// ---------- prop_tree ----------

type node struct {
	box      *lmc.Mailbox
	children []*node
	root     bool
}

func (n *node) add(from *lmc.Mailbox, to *lmc.Mailbox) bool {
	if (!n.root && n.box == nil && from == nil) || (n.box != nil && n.box.Address() == from.Address()) {
		n.children = append(n.children, &node{box: to})
		return true
	}

	for _, c := range n.children {
		if c.add(from, to) {
			return true
		}
	}

	if n.root {
		n.children = append(n.children, &node{box: from, children: []*node{{box: to}}})
		return true
	}

	return false
}

func (n *node) propagate(x *lmc.Mailbox, p *lmc.Program) {
	instrs := p.Memory.InstructionsList.Instructions

	if !n.root && n.box != nil && x != nil {
		for _, i := range instrs {
			if len(i.Boxes()) != 0 {
				if i.Boxes()[0].Address() == n.box.Address() {
					*i.Boxes()[0] = *x
				}
			}
		}
	}

	for _, c := range n.children {
		c.propagate(n.box, p)
	}
}

func prop_tree(prog *lmc.Program) error {
	root := node{root: true}

	instrs := prog.Memory.InstructionsList.Instructions
	for k, instr := range instrs {
		var ok bool

		if _, ok = instr.(*lmc.StoreInstr); !ok {
			continue
		}

		for kk := k - 1; kk >= 0; kk-- {
			i := instrs[kk]

			if _, ok := i.(*lmc.LoadInstr); ok {
				root.add(i.Boxes()[0], instr.Boxes()[0])
				break
			} else if _, ok = i.(*lmc.InputInstr); ok {
				root.add(nil, instr.Boxes()[0])
				break
			} else if i.ACC() {
				break
			}
		}
	}

	root.propagate(nil, prog)
	return nil
}

// ---------- prop_lda_sta ----------

func prop_lda_sta(prog *lmc.Program) error {
	instrs := make([]lmc.Instruction, len(prog.Memory.InstructionsList.Instructions))
	copy(instrs, prog.Memory.InstructionsList.Instructions)

	previous := -1

	for i, removed := 0, 0; i < len(instrs); i++ {
		var ok bool

		if _, ok = instrs[i].(*lmc.LoadInstr); ok {
			previous = i
			continue
		}

		ok = previous != -1
		if ok {
			if _, ok2 := instrs[i].(*lmc.StoreInstr); !ok2 {
				ok = false
			}
		}

		if ok {
			if instrs[i].Boxes()[0].Address() != instrs[previous].Boxes()[0].Address() {
				previous = -1
				continue
			}

			remove := true

			for j := previous + 1; j < i; j++ {
				if instrs[j].ACC() {
					remove = false
					break
				}
			}

			if remove {
				if err := prog.Memory.InstructionsList.RemoveInstruction(i - removed); err != nil {
					return err
				} else {
					removed++
				}
			}
		}
	}

	return nil
}

// ---------- OProp ----------

type OProp struct {
	program *lmc.Program
}

func NewOProp(program *lmc.Program) *OProp {
	return &OProp{
		program: program,
	}
}

func (o *OProp) Strategy() OStrategy {
	return BProp
}

func (o *OProp) Optimise() error {
	var err error

	_ = prop_tree(o.program) // currently returns no err

	if err = prop_lda_sta(o.program); err != nil {
		return propErr(1, err)
	}

	return nil
}

func (o *OProp) Program() *lmc.Program {
	return o.program
}

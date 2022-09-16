package optimisation

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/lmc"
)

var cleanStageNames = [...]string{
	"CLEAN_DEAD",
	"CLEAN_MULTI",
}

func cleanErr(stage int, child error) error {
	return fmt.Errorf("cleaning failed stage %d=%s: %s", stage, cleanStageNames[stage], child)
}

func clean_dead_box(prog *lmc.Program) error {
	used := make(map[string]struct{})

	for _, v := range prog.Memory.InstructionsList.Instructions {
		for _, box := range v.Boxes() {
			used[box.Identifier()] = struct{}{}
		}
	}

	var ok bool

	for _, def := range prog.Memory.InstructionsList.DefInstructions {
		if _, ok = used[def.Box.Identifier()]; !ok {
			_ = prog.Memory.InstructionsList.RemoveDef(def.Box.Identifier()) // ignore error
			prog.Memory.RemoveMailboxIdentifier(def.Box.Identifier())
		}
	}

	return nil
}

func clean_multi_dat(prog *lmc.Program) error {
	seen := make(map[lmc.Address]struct{})
	instrs := prog.Memory.InstructionsList.DefInstructions

	var ok bool
	var i int

	for _, ii := range instrs {
		if _, ok = seen[ii.Box.Address()]; !ok {
			instrs[i] = ii
			i++

			seen[ii.Box.Address()] = struct{}{}
		}
	}

	for j := i; j < len(instrs); j++ {
		instrs[j] = nil
	}

	prog.Memory.InstructionsList.DefInstructions = instrs[:i]
	return nil
}

type OClean struct {
	program *lmc.Program
}

func NewOClean(program *lmc.Program) *OClean {
	return &OClean{
		program: program,
	}
}

func (o *OClean) Strategy() OStrategy {
	return Clean
}

func (o *OClean) Optimise() error {
	_ = clean_dead_box(o.program)  // currently no error returned
	_ = clean_multi_dat(o.program) // ^^

	return nil
}

func (o *OClean) Program() *lmc.Program {
	return o.program
}

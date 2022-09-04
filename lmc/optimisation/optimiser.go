package optimisation

import "github.com/clr1107/lmc-llvm-target/lmc"

type OStrategy uint

const (
	Thrashing OStrategy = iota
	Waste
	BProp
	Chaining
	Unroll
	Stacking
)

type Optimiser interface {
	Optimise() error
	Program() *lmc.Program
	Strategy() OStrategy
}

type StackingOptimiser struct {
	program *lmc.Program
	strategies []OStrategy
}

func NewStackingOptimiser(program *lmc.Program, strategies []OStrategy) *StackingOptimiser {
	return &StackingOptimiser{
		program: program,
		strategies: strategies,
	}
}

func (o *StackingOptimiser) createStrategy(s OStrategy) Optimiser {
	switch s {
	case Thrashing:
		return NewOThrashing(o.program)
	case Waste:
		return NewOWaste(o.program)
	default:
		return nil
	}
}

func (o *StackingOptimiser) Optimise() error {
	var optimiser Optimiser
	var err error
	waste := NewOWaste(o.program)

	for _, s := range o.strategies {

		optimiser = o.createStrategy(s)

		if optimiser != nil {
			if err = optimiser.Optimise(); err != nil {
				return err
			} else if s != Waste { // no point running it twice
				if err = waste.Optimise(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (o *StackingOptimiser) Program() *lmc.Program {
	return o.program
}

func (o *StackingOptimiser) Strategy() OStrategy {
	return Stacking
}


package optimisation

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/lmc"
)

type OStrategy uint

const (
	Thrashing OStrategy = iota
	Clean
	BProp
	Chaining
	Unroll
	Stacking
)

var OStrategyNames = map[OStrategy]string{
	Thrashing: "THRASHING",
	Clean:     "CLEAN",
	BProp:     "BOX_PROPAGATION",
	Chaining:  "ADD_CHAIN",
	Unroll:    "UROLL",
	Stacking:  "OSTACK",
}

type Optimiser interface {
	Optimise() error
	Program() *lmc.Program
	Strategy() OStrategy
}

type StackingOptimiser struct {
	program    *lmc.Program
	strategies []OStrategy
}

func NewStackingOptimiser(program *lmc.Program, strategies []OStrategy) *StackingOptimiser {
	return &StackingOptimiser{
		program:    program,
		strategies: strategies,
	}
}

func (o *StackingOptimiser) createStrategy(s OStrategy) Optimiser {
	switch s {
	case Thrashing:
		return NewOThrashing(o.program)
	case Clean:
		return NewOClean(o.program)
	case BProp:
		return NewOProp(o.program)
	default:
		return nil
	}
}

func (o *StackingOptimiser) Optimise() error {
	var optimiser Optimiser
	var err error
	waste := NewOClean(o.program)

	for _, s := range o.strategies {

		optimiser = o.createStrategy(s)

		if optimiser != nil {
			if err = optimiser.Optimise(); err != nil {
				return fmt.Errorf("stacking optimisation, strategy %s: %s", OStrategyNames[optimiser.Strategy()], err)
			} else if s != Clean { // no point running it twice
				if err = waste.Optimise(); err != nil {
					return fmt.Errorf("stacking clean cycle error: %s", err)
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

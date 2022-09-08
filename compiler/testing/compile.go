// An example compiler. It just takes in an ll file, compiles it, and optimises it.
package main

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/compiler"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/clr1107/lmc-llvm-target/lmc/optimisation"
	"github.com/llir/llvm/asm"
	"os"
)

func main() {
	c := compiler.NewCompiler(lmc.NewProgram(lmc.NewBasicMemory()))

	mod, err := asm.ParseFile("ll/test.ll")
	if err != nil {
		fmt.Printf("error whilst parsing file: %s\n", err)
		os.Exit(1)
	}

	f := compiler.GetLLEntry(mod)
	if f == nil {
		println("could not find entry function")
		os.Exit(1)
	}

	for _, block := range f.Blocks {
		for _, instr := range block.Insts {

			if compiled := c.CompileInst(instr); compiled.Err != nil {
				fmt.Printf("could not compile instruction: %s\n\t%s\n", instr.LLString(), compiled.Err)
				os.Exit(1)
			} else {
				if len(compiled.Warnings) > 0 {
					fmt.Printf("Warnings: %s\n", instr.LLString())
					for _, w := range compiled.Warnings {
						fmt.Printf("\t%s\n", w)
					}
				}

				if compiled.Wrapped != nil {
					if err := c.AddCompiledInstruction(compiled.Wrapped); err != nil {
						fmt.Printf("could not add a compiled instruction: %s\n\t%s\n", instr.LLString(), err)
						os.Exit(1)
					}
				}
			}

		}
	}

	fmt.Printf("Unoptimised %d instrs, %d defs:\n%s\n", len(c.Prog.Memory.GetInstructionSet().GetInstructions()), len(c.Prog.Memory.GetInstructionSet().GetDefs()), c.Prog)
	fmt.Printf("\n\n")

	optimiser := optimisation.NewStackingOptimiser(c.Prog, []optimisation.OStrategy{
		optimisation.Thrashing, optimisation.BProp,
	})

	if err := optimiser.Optimise(); err != nil {
		fmt.Printf("could not optimise program: %s\n", err)
		os.Exit(1)
	}

	prog := optimiser.Program()
	fmt.Printf("Optimised %d instrs, %d defs:\n%s\n", len(prog.Memory.GetInstructionSet().GetInstructions()), len(prog.Memory.GetInstructionSet().GetDefs()), prog)
}

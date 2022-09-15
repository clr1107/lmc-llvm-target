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
	comp := compiler.NewCompiler(lmc.NewProgram(lmc.NewBasicMemory()))
	engine := compiler.NewEngine(comp)

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
		matches, err := engine.FindAll(block.Insts)

		if err != nil {
			fmt.Printf("error pattern matching instructions\n\t%s\n", err)
			os.Exit(1)
		}

		for _, m := range matches {
			c := m.Pattern.Compile(m.Instrs)

			if c.Err != nil {
				fmt.Printf("could not compile instructions: ")
				for _, i := range m.Instrs {
					fmt.Printf("%s, ", i.LLString())
				}
				fmt.Printf("\n")

				fmt.Printf("\t%s\n", c.Err)
				os.Exit(1)
			}

			if len(c.Warnings) != 0 {
				fmt.Printf("Warnings: ")
				for _, i := range m.Instrs {
					fmt.Printf("%s, ", i.LLString())
				}
				fmt.Printf("\n")

				fmt.Printf("\t%s\n", c.Err)
			}

			if err := comp.AddCompiledInstruction(c.Wrapped); err != nil {
				fmt.Printf("could not add compiled instruction: ")
				for _, i := range m.Instrs {
					fmt.Printf("%s, ", i.LLString())
				}
				fmt.Printf("\n")

				fmt.Printf("\t%s\n", err)
			}
		}
	}

	fmt.Printf("Unoptimised %d instrs, %d defs:\n%s\n", len(comp.Prog.Memory.GetInstructionSet().Instructions), len(comp.Prog.Memory.GetInstructionSet().DefInstructions), comp.Prog)
	fmt.Printf("\n\n")

	optimiser := optimisation.NewStackingOptimiser(comp.Prog, []optimisation.OStrategy{
		optimisation.BProp, optimisation.Thrashing,
	})

	if err := optimiser.Optimise(); err != nil {
		fmt.Printf("could not optimise program: %s\n", err)
		os.Exit(1)
	}

	prog := optimiser.Program()
	fmt.Printf("Optimised %d instrs, %d defs:\n%s\n", len(prog.Memory.GetInstructionSet().Instructions), len(prog.Memory.GetInstructionSet().DefInstructions), prog)
}

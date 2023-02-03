// An example compiler. It just takes in an ll file, compiles it, and optimises it.
package main

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/compiler"
	"github.com/clr1107/lmc-llvm-target/compiler/errors"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/clr1107/lmc-llvm-target/lmc/optimisation"
	"github.com/llir/llvm/asm"
	"os"
	"strings"
)

func main() {
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

	comp := compiler.NewCompiler(lmc.NewProgram(lmc.NewBasicMemory()))
	comp.Module = mod

	engine := compiler.NewEngine(comp)

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
				var warningBuf strings.Builder

				for _, w := range c.Warnings {
					if w.Level <= errors.WarningLevel(comp.Options.Get("WLEVEL").Value.(int)) {
						warningBuf.WriteString(fmt.Sprintf("\t%s\n", w))
					}
				}

				if warningBuf.Len() > 0 {
					fmt.Printf("Warnings: ")
					for _, i := range m.Instrs {
						fmt.Printf("%s, ", i.LLString())
					}
					fmt.Printf("\n%s", warningBuf.String())
				}
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

	fmt.Printf("\nUnoptimised %d instrs, %d defs:\n%s\n", len(comp.Prog.Memory.InstructionsList.Instructions), len(comp.Prog.Memory.InstructionsList.DefInstructions), comp.Prog)
	fmt.Printf("\n\n")

	var strategies []optimisation.OStrategy

	optValue := comp.Options.Get("OPT").Value.(optimisation.OStrategy)
	for l := 0; l <= 3; l++ {
		mask := (optValue >> l) & 1
		if mask == 1 {
			if l == 0 {
				strategies = append(strategies, optimisation.OStrategy(1))
			} else {
				strategies = append(strategies, optimisation.OStrategy(2<<(l-1)))
			}
		}
	}

	optimiser := optimisation.NewStackingOptimiser(comp.Prog, strategies)

	if err := optimiser.Optimise(); err != nil {
		fmt.Printf("could not optimise program: %s\n", err)
		os.Exit(1)
	}

	prog := optimiser.Program()
	fmt.Printf("Optimised %d instrs, %d defs:\n%s\n", len(prog.Memory.InstructionsList.Instructions), len(prog.Memory.InstructionsList.DefInstructions), prog)

	fmt.Printf("\n\nCompiler options:\n%s\n", comp.Options)
}

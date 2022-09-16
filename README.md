<div align="center"><h1 >LMC-LLVM-Target</h1>
<p><small><i>Massive work in progress.</i></small></p>
<img alt="Logo: 'LMC Simulator' executable icon, by Prof Magnus Bordewich, Uni of Durham, overlayed on LLVM logo" width="20%" src="https://i.clr.is/5HHGpYVds.png" />
<p><small>(<a href="https://mjrbordewich.webspace.durham.ac.uk/lmc/">'LMC Simulator'</a> executable icon by Professor Magnus Bordewich, University of Durham, overlayed on LLVM logo. Temporary, maybe.)</small></p>
</div>

## Description

Little Man Computer (LMC) is an incredibly simple instruction set, used for teaching, that gives a model of a computer  
in the Von Neumann Architecture. It was first introduced by Dr. Stuart Madnick of M.I.T. in 1965. This program provides
a compiler back-end to translate LLVM's IR to LMC. So, in theory, any language compilable to IR can be compiled to LMC.
This project focuses on compiling C to
LMC. [Here is a list of the instructions.](https://en.wikipedia.org/wiki/Little_man_computer#Instructions) (Note: there
are two main conventions for instruction mnemonics I've seen. I am using the ones listed on Wikipedia.)

The purpose of this project is mainly for entertainment. It is quite ridiculous that C could compile to this simple
and (functionally) useless language. Also, this would have made homework at school and my first year university project
far easier.

### Project description

This project includes two standalone packages: `lmc` which allows the construction, analysis, and optimisation of LMC
programs, and `compiler` which will take LLIR and compile it to (unoptimised) LMC instructions. (The `lmc` package can
then be used to optimise it.)

Once the project is at a suitable level of completeness I will include a CLI tool to compile and optimise code from
inputs. For now, there is a program in `compiler/testing/compiler.go` that gives an incredibly simple example of how
this tool may look. It does compilation and optimisation.

The compiler is really simple. I mean, extremely basic. It performs rudimentary pattern matching on IR instructions,
converts them to LMC instructions, producing some of the worst LMC in existence, before optimising it.

LMC typically has a limit of 100 mailboxes, however this is not observed or worried about in this project. Mainly
because I can't be bothered to include that limitation, and my excuse is I am treating it as RAM; a program does not
know ahead of time how large its target device's RAM will be, so neither does this compiler.

## LMC Package

An overview and examples of the `lmc` package.

### Memory

*TODO*

### Utility functions

*TODO*

### Constructing a program

*TODO*

### Optimising a program

*TODO*

## Compiler package

An overview and examples of the `compiler` package. For details of optimisation algorithms included
see [`lmc/optimisation/OPTIMISTAION.md`](lmc/optimisation/OPTIMISATION.md)

### Steps of the compiler

*TODO*

### LMC header file (`lmc.h`)

[`compiler/lmc.h`](compiler/lmc.h) contains some very useful macros and builtin functions. These functions are
marked `extern` and when pattern matched are, essentially, drop-in replaced by some LMC instructions. There are two
types that are _recommended_ to be used when creating a program: `number_t` and `bool_t`.

Also included in this header file is the entry point to any LMC program, `void _lmc(void)`.

### How to use the compiler

*TODO*

## All round examples

*TODO*

### Optimising a program

*TODO*

### Compiling and optimising a program

*TODO*
  


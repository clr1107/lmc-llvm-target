# An example Makefile for testing the compiler
# Compiles all files matching test/*.c to ll/X.ll individually
# Run `make`

CC = clang
CCFLAGS = -emit-llvm -nostdlib -S -I../ -O0

COMPILED = ll/simple_example.ll ll/test.ll

ll/%.ll: c/%.c
	$(CC) $< $(CCFLAGS) -o $@

.PHONY: all clean

all: $(COMPILED)

clean:
	rm -rf ll/*

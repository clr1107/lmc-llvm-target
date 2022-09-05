#!/bin/bash
#
# Compiles a C source to an LLVM output.
#

OUTPUT="out.ll"
FLAGS="-O3"

if [ "$#" -lt 1 ]; then
  echo "Usage: \`cc.sh <C source> [output] [clang flags ...]\`"
  exit 1
fi

if [ "$#" -ge 2 ]; then
  OUTPUT="$2"
fi

if [ "$#" -gt 2 ];  then
  FLAGS="${*:3}"
fi

clang -emit-llvm -nostdlib -S -I. "$FLAGS" -o "$OUTPUT" "$1"
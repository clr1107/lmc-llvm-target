module github.com/clr1107/lmc-llvm-target/compiler

go 1.17

replace github.com/clr1107/lmc-llvm-target/lmc v0.0.0 => ./../lmc

require (
	github.com/clr1107/lmc-llvm-target/lmc v0.0.0
	github.com/llir/llvm v0.3.4
)

require (
	github.com/mewmew/float v0.0.0-20211212214546-4fe539893335 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	golang.org/x/tools v0.1.10 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

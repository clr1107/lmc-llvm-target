package lmc

import "fmt"

type Address int64
type Value int

type LMCType interface {
	fmt.Stringer
	LMCString() string
}

type Addressable interface {
	Address() Address
}

type Identifiable interface {
	Identifier() string
}

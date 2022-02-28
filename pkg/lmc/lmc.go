package lmc

type Address int64

type Addressable interface {
	Address() Address
}

type Identifiable interface {
	Identifier() string
}

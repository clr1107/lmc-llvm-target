package compiler

import (
	"github.com/clr1107/lmc-llvm-target/compiler/errors"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"reflect"
	"sort"
)

func GetLLFunc(funcs []*ir.Func, name string) *ir.Func {
	for _, v := range funcs {
		if v.Name() == name {
			return v
		}
	}

	return nil
}

func GetLLInstrs(blocks []*ir.Block) []ir.Instruction {
	instrs := make([]ir.Instruction, 0)

	for _, block := range blocks {
		instrs = append(instrs, block.Insts...)
	}

	return instrs
}

func GetLLEntry(module *ir.Module) *ir.Func {
	f := GetLLFunc(module.Funcs, "_lmc")

	if f == nil {
		return nil
	}

	if len(f.Params) != 0 || f.Sig.RetType.String() != "void" {
		return nil
	}

	return f
}

func ReflectGetLocalID(x interface{}) (lmc.Address, error) {
	f := reflect.ValueOf(x).MethodByName("ID")
	if f.IsZero() {
		return -1, errors.E_NonexistentProperty("#ID()", nil)
	}

	res := f.Call(nil)
	if len(res) != 1 {
		return -1, errors.E_NonexistentProperty("#ID() []", nil)

	}

	id := res[0]
	if id.Kind() != reflect.Int64 {
		return -1, errors.E_IncorrectType(nil, "#ID()", id.Kind().String(), "int64")
	}

	return lmc.Address(id.Int()), nil
}

func ReflectGetProperty(x interface{}, field string) (interface{}, error) {
	property := reflect.ValueOf(x).FieldByName(field)
	if property.IsZero() {
		return nil, errors.E_NonexistentProperty("#"+field, nil)
	}

	return property.Interface(), nil
}

func ValidLLType(t types.Type) bool {
	if types.IsInt(t) {
		return true
	}

	if c, ok := t.(*types.PointerType); ok {
		return ValidLLType(c.ElemType)
	}

	return false
}

type OrderedSlice struct {
	s        []interface{}
	comparer func(interface{}, interface{}) bool
}

func NewOrderedSlice(l int, comparer func(interface{}, interface{}) bool) *OrderedSlice {
	return &OrderedSlice{
		s:        make([]interface{}, l),
		comparer: comparer,
	}
}

func (o *OrderedSlice) Append(x interface{}) {
	i := sort.Search(len(o.s), func(i int) bool {
		return o.comparer(o.s[i], x)
	})

	if i == len(o.s) {
		o.s = append(o.s, x)
	} else {
		o.s = append(o.s[:i+1], o.s[i:]...)
		o.s[i] = x
	}
}

func (o *OrderedSlice) Get(i int) interface{} {
	return o.s[i]
}

func (o *OrderedSlice) Slice() []interface{} {
	ss := make([]interface{}, len(o.s))
	copy(ss, o.s)

	return ss
}

func (o *OrderedSlice) Len() int {
	return len(o.s)
}

func MultiPredicate(predicate func(interface{}) bool, x ...interface{}) bool {
	for _, v := range x {
		if !predicate(v) {
			return false
		}
	}

	return true
}

func MultiVariablePredicate(predicate func(int, interface{}) bool, x ...interface{}) bool {
	for k, v := range x {
		if !predicate(k, v) {
			return false
		}
	}

	return true
}

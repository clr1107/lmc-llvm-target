package errors

import (
	"fmt"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"reflect"
	"strings"
)

type ErrorCode uint8

const (
	LMCError ErrorCode = iota
	NonexistentPropertyError
	IncorrectTypeError
	InvalidLLTypesError
	UnknownMailboxError
	UnknownLLInstructionError
	BuiltinInvocationError
	UnknownBuiltinError
)

var ErrorNames = map[ErrorCode]string{
	LMCError:                  "LMC_LIB",
	NonexistentPropertyError:  "NONEXISTENT_PROPERTY",
	IncorrectTypeError:        "INCORRECT_TYPE",
	InvalidLLTypesError:       "INVALID_LL_TYPES",
	UnknownMailboxError:       "UNKNOWN_MAILBOX",
	UnknownLLInstructionError: "UNKNOWN_LL_INSTRUCTION",
	BuiltinInvocationError:    "BUILTIN_INVOCATION",
	UnknownBuiltinError:       "UNKNOWN_BUILTIN",
}

type Error struct {
	error
	Code  ErrorCode
	Child error
	msg   string
}

func NewError(code ErrorCode, msg string, child error) *Error {
	return &Error{
		Code:  code,
		msg:   msg,
		Child: child,
	}
}

func (e *Error) Error() string {
	s := fmt.Sprintf("%d=%s; %s", e.Code, ErrorNames[e.Code], e.msg)

	if e.Child != nil {
		s += ": " + e.Child.Error()
	}

	return s
}

// ---------- Errors definitions ----------

func E_LMC(msg string, child error) *Error {
	if msg != "" {
		msg = ", " + msg
	}

	return NewError(LMCError, fmt.Sprintf("LMC error%s", msg), child)
}

func E_NonexistentProperty(property string, child error) *Error {
	return NewError(NonexistentPropertyError, fmt.Sprintf("nonexistent property `%s`", property), child)
}

func E_IncorrectType(child error, value string, got string, expected string) *Error {
	return NewError(IncorrectTypeError, fmt.Sprintf("incorrect type for %s got %s expected %s", value, got, expected), child)
}

func E_InvalidLLTypes(child error, types ...string) *Error {
	return NewError(InvalidLLTypesError, fmt.Sprintf("invalid LL types %s", strings.Join(types, ", ")), child)
}

func E_UnknownMailbox(addr lmc.Address, child error) *Error {
	return NewError(UnknownMailboxError, fmt.Sprintf("unknown mailbox with address %d", addr), child)
}

func E_UnknownLLInstruction(instr ir.Instruction, child error) *Error {
	return NewError(
		UnknownLLInstructionError,
		fmt.Sprintf("unknown LL instruction type %s", reflect.TypeOf(instr)),
		child,
	)
}

func E_BuiltinInvocation(signature string, child error) *Error {
	return NewError(BuiltinInvocationError, fmt.Sprintf("invocation %s", signature), child)
}

func E_UnknownBuiltin(name string, child error) *Error {
	return NewError(UnknownBuiltinError, fmt.Sprintf("unknown builtin function %s", name), child)
}

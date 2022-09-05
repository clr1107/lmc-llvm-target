package compiler

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
	IncorrectTypesError
	InvalidLLTypesError
	UnknownMailboxError
	UnknownLLInstructionError
)

var ErrorNames = map[ErrorCode]string{
	LMCError:                  "LMC_LIB",
	NonexistentPropertyError:  "NONEXISTENT_PROPERTY",
	IncorrectTypesError:       "INCORRECT_TYPES",
	InvalidLLTypesError:       "INVALID_LL_TYPES",
	UnknownMailboxError:       "UNKNOWN_MAILBOX",
	UnknownLLInstructionError: "UNKNOWN_LL_INSTRUCTION",
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

func E_IncorrectTypes(child error, types ...string) *Error {
	return NewError(IncorrectTypesError, fmt.Sprintf("incorrect types %s", strings.Join(types, ", ")), child)
}

func E_InvalidLLTypes(child error, types ...string) *Error {
	return NewError(IncorrectTypesError, fmt.Sprintf("invalid LL types %s", strings.Join(types, ", ")), child)
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

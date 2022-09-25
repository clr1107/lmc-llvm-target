package errors

import (
	"errors"
	"fmt"
	"github.com/clr1107/lmc-llvm-target/lmc"
	"github.com/llir/llvm/ir"
	"reflect"
	"strings"
)

type ErrorCode uint8

const (
	Err ErrorCode = iota
	LMCError
	UnsupportedError
	NonexistentPropertyError
	IncorrectTypeError
	InvalidLLTypesError
	UnknownMailboxError
	UnknownLLInstructionError
	BuiltinInvocationError
	UnknownBuiltinError
	InvalidOptionSyntaxError
)

var errorNames = map[ErrorCode]string{
	Err:                       "ERR",
	LMCError:                  "LMC_LIB",
	UnsupportedError:          "UNSUPPORTED",
	NonexistentPropertyError:  "NONEXISTENT_PROPERTY",
	IncorrectTypeError:        "INCORRECT_TYPE",
	InvalidLLTypesError:       "INVALID_LL_TYPES",
	UnknownMailboxError:       "UNKNOWN_MAILBOX",
	UnknownLLInstructionError: "UNKNOWN_LL_INSTRUCTION",
	BuiltinInvocationError:    "BUILTIN_INVOCATION",
	UnknownBuiltinError:       "UNKNOWN_BUILTIN",
	InvalidOptionSyntaxError:  "INVALID_OPT_SYNTAX",
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
	s := fmt.Sprintf("%d=%s; %s", e.Code, errorNames[e.Code], e.msg)

	if e.Child != nil {
		s += ": " + e.Child.Error()
	}

	return s
}

// ---------- Errors definitions ----------

func E_Err(msg string, child error) *Error {
	if msg == "" {
		msg = "compiler error"
	}

	return NewError(Err, fmt.Sprintf("%s", msg), child)
}

func E_LMC(msg string, child error) *Error {
	if msg != "" {
		msg = ", " + msg
	}

	return NewError(LMCError, fmt.Sprintf("LMC error%s", msg), child)
}

func E_Unsupported(feature string, child error) *Error {
	return NewError(UnsupportedError, feature, child)
}

func E_NonexistentProperty(property string, child error) *Error {
	return NewError(NonexistentPropertyError, fmt.Sprintf("nonexistent property `%s`", property), child)
}

func E_IncorrectType(child error, value string, got string, expected string) *Error {
	return NewError(IncorrectTypeError, fmt.Sprintf("incorrect type for %s got %s expected %s", value, got, expected), child)
}

func E_InvalidLLTypes(child error, types ...string) *Error {
	spacer := ""
	if len(types) > 0 {
		spacer = " "
	}

	return NewError(InvalidLLTypesError, fmt.Sprintf("invalid LL types%s%s", spacer, strings.Join(types, ", ")), child)
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

func E_UnknownBuiltin(f *ir.Func, child error) *Error {
	return NewError(UnknownBuiltinError, fmt.Sprintf("unknown builtin function %s(%d)", f.Name(), len(f.Params)), child)
}

func E_InvalidOptionSyntax(problem string) *Error {
	return NewError(InvalidOptionSyntaxError, "invalid compiler option syntax (__lmc_option__)", errors.New(problem))
}

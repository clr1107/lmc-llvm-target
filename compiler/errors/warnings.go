package errors

import (
	"fmt"
)

type WarningCode uint8
type WarningLevel uint8

const (
	BitcastWarning WarningCode = iota
)

var warningNames = map[WarningCode]string{
	BitcastWarning: "BITCAST",
}

const (
	L_Default WarningLevel = iota
	L_Info
	L_Debug
)

var warningLevelNames = map[WarningLevel]string{
	L_Default: "DEFAULT",
	L_Info:    "INFO",
	L_Debug:   "DEBUG",
}

type Warning struct {
	msg   string
	Code  WarningCode
	Level WarningLevel
}

func (e *Warning) String() string {
	return fmt.Sprintf("%d=%s, %d=%s: %s", e.Code, warningNames[e.Code], e.Level, warningLevelNames[e.Level], e.msg)
}

// ---------- Warnings definitions ----------

func W_Bitcast(from string, to string) *Warning {
	return &Warning{Code: BitcastWarning, Level: L_Debug, msg: fmt.Sprintf("bitcast from %s to %s", from, to)}
}

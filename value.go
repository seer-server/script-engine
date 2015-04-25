package engine

import (
	"github.com/yuin/gopher-lua"

	"fmt"
)

// Provides information about failed Value typecasts.
type ValueError string

func newValueError(exp string, v *Value) ValueError {
	return ValueError(fmt.Sprintf("expected %s, found \"%s\"", exp, v.lval.Type()))
}

// Implements the Error interface for ValueError
func (v ValueError) Error() string {
	return string(v)
}

// Utility wrapper for lua.LValue that provies conveinient methods for casting.
type Value struct {
	lval lua.LValue
}

func newValue(val lua.LValue) *Value {
	return &Value{
		lval: val,
	}
}

// Returns the LValue as a Go string
func (v *Value) AsString() (string, error) {
	if str, ok := v.lval.(lua.LString); ok {
		return string(str), nil
	} else {
		return "", newValueError("string", v)
	}
}

// Returns the LValue as a Go
func (v *Value) AsFloat() (float64, error) {
	if f, ok := v.lval.(lua.LNumber); ok {
		return float64(f), nil
	} else {
		return 0.0, newValueError("number", v)
	}
}

func (v *Value) AsNumber() (float64, error) {
	return v.AsFloat()
}

func (v *Value) AsBool() (bool, error) {
	if b, ok := v.lval.(lua.LBool); ok {
		return b == lua.LTrue, nil
	} else {
		return false, newValueError("bool", v)
	}
}

func (v *Value) IsNil() bool {
	return v.lval.Type() == lua.LTNil
}

package lua

import glua "github.com/yuin/gopher-lua"

// Nil
var Nil = newValue(glua.LNil)

// Number converts a float to a Value representing a number in Lua.
func Number(f float64) *Value {
	return newValue(glua.LNumber(f))
}

// String returns a Value representing a string in Lua.
func String(s string) *Value {
	return newValue(glua.LString(s))
}

// Bool converts a Go bool into a Lua bool type.
func Bool(b bool) *Value {
	if b {
		return newValue(glua.LTrue)
	}

	return newValue(glua.LFalse)
}

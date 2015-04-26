package engine

import "github.com/yuin/gopher-lua"

// LuaNumber converts a float to a Value representing a number in Lua.
func LuaNumber(f float64) *Value {
	return newValue(lua.LNumber(f))
}

// Lua string returns a Value representing a string in Lua.
func LuaString(s string) *Value {
	return newValue(lua.LString(s))
}

// LuaNil returns the Nil value for Lua.
func LuaNil() *Value {
	return newValue(lua.LNil)
}

// LuaBool converts a Go bool into a Lua bool type.
func LuaBool(b bool) *Value {
	if b {
		return newValue(lua.LTrue)
	} else {
		return newValue(lua.LFalse)
	}
}

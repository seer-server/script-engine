package engine

import "github.com/yuin/gopher-lua"

func LuaNumber(f float64) *Value {
	return newValue(lua.LNumber(f))
}

func LuaString(s string) *Value {
	return newValue(lua.LString(s))
}

func LuaNil() *Value {
	return newValue(lua.LNil)
}

func LuaBool(b bool) *Value {
	if b {
		return newValue(lua.LTrue)
	} else {
		return newValue(lua.LFalse)
	}
}

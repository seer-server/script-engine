// Copyright (c) 2015 tree-server contributors

package engine

import "github.com/yuin/gopher-lua"

// Engine struct stores a pointer to a lua.LState providing a simplified API.
type Engine struct {
	state *lua.LState
}

// ScriptFunction is a type alias for a function that receives an Engine and
// returns an int.
type ScriptFunction func(*Engine) int

// ScriptFnMap is a type alias for map[string]ScriptFunction
type ScriptFnMap map[string]ScriptFunction

// Create a new engine containing a new lua.LState.
func NewEngine() *Engine {
	return &Engine{
		state: lua.NewState(),
	}
}

// Close will perform a close on the Lua state.
func (e *Engine) Close() {
	e.state.Close()
}

// LoadFile runs the file through the Lua interpreter.
func (e *Engine) LoadFile(fn string) error {
	return e.state.DoFile(fn)
}

// LoadString runs the given string through the Lua interpreter.
func (e *Engine) LoadString(src string) error {
	return e.state.DoString(src)
}

// SetVal allows for setting global variables in the loaded code.
func (e *Engine) SetGlobal(name string, val *Value) {
	e.state.SetGlobal(name, val.lval)
}

// SetField applies the value to the given table associated with the given
// key.
func (e *Engine) SetField(tbl *Value, key string, v *Value) {
	e.state.SetField(tbl.lval, key, v.lval)
}

// RegisterFunc registers a Go function (ScriptFunction) with the script.
// Using this method makes Go functions accessible through Lua scripts.
func (e *Engine) RegisterFunc(name string, fn ScriptFunction) {
	e.state.SetGlobal(name, e.genScriptFunc(fn))
}

// RegisterModule registers a Go module with the Engine for use within Lua.
func (e *Engine) RegisterModule(name string, loadFn func(*Engine) *Value) {
	loader := func(l *lua.LState) int {
		e := &Engine{l}
		mod := loadFn(e)
		e.PushRet(mod)

		return 1
	}

	e.state.PreloadModule(name, loader)
}

// GenerateModule returns a table that has been loaded with the given script
// function map.
func (e *Engine) GenerateModule(fnMap ScriptFnMap) *Value {
	tbl := e.state.NewTable()
	realFnMap := make(map[string]lua.LGFunction)
	for k, fn := range fnMap {
		realFnMap[k] = e.wrapScriptFunction(fn)
	}

	mod := e.state.SetFuncs(tbl, realFnMap)

	return newValue(mod)
}

// PopArg returns the top value on the Lua stack.
// This method is used to get arguments given to a Go function from a Lua script.
// This method will return a Value pointer that can then be converted into
// an appropriate type.
func (e *Engine) PopArg() *Value {
	lv := e.state.Get(-1)
	e.state.Pop(1)

	return newValue(lv)
}

// PushRet pushes the given Value onto the Lua stack.
// Use this method when 'returning' values from a Go function called from a
// Lua script.
func (e *Engine) PushRet(v *Value) {
	e.state.Push(v.lval)
}

// Call allows for calling a method by name.
// The second parameter is the number of return values the function being
// called should return. These values will be returned in a slice of Value
// pointers.
func (e *Engine) Call(name string, retCount int, params ...*Value) ([]*Value, error) {
	luaParams := make([]lua.LValue, len(params))
	for i, v := range params {
		luaParams[i] = v.lval
	}

	err := e.state.CallByParam(lua.P{
		Fn:      e.state.GetGlobal(name),
		NRet:    retCount,
		Protect: true,
	}, luaParams...)

	if err != nil {
		return nil, err
	}

	retVals := make([]*Value, retCount)
	for i := 0; i < retCount; i++ {
		retVals[i] = newValue(e.state.Get(-1))
	}
	e.state.Pop(retCount)

	return retVals, nil
}

// LuaTable creates and returns a new LuaTable.
func (e *Engine) LuaTable() *Value {
	return newValue(e.state.NewTable())
}

// wrapScriptFunction turns a ScriptFunction into a lua.LGFunction
func (e *Engine) wrapScriptFunction(fn ScriptFunction) lua.LGFunction {
	return func(l *lua.LState) int {
		e := &Engine{state: l}

		return fn(e)
	}
}

// genScriptFunc will wrap a ScriptFunction with a function that gopher-lua
// expects to see when calling method from Lua.
func (e *Engine) genScriptFunc(fn ScriptFunction) *lua.LFunction {
	return e.state.NewFunction(e.wrapScriptFunction(fn))
}

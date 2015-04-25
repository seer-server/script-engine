package engine

import "github.com/yuin/gopher-lua"

type Engine struct {
	state *lua.LState
}

type ScriptFunction func(*Engine) int

func NewEngine() *Engine {
	return &Engine{
		state: lua.NewState(),
	}
}

func (e *Engine) Close() {
	e.state.Close()
}

func (e *Engine) LoadFile(fn string) error {
	return e.state.DoFile(fn)
}

func (e *Engine) LoadString(src string) error {
	return e.state.DoString(src)
}

func (e *Engine) SetVal(name string, val *Value) {
	e.state.SetGlobal(name, val.lval)
}

func (e *Engine) RegisterFunc(name string, fn ScriptFunction) {
	e.state.SetGlobal(name, e.genScriptFunc(fn))
}

func (e *Engine) genScriptFunc(fn ScriptFunction) lua.LValue {
	sfn := func(l *lua.LState) int {
		e := &Engine{state: l}

		return fn(e)
	}

	return e.state.NewFunction(sfn)
}

func (e *Engine) PopArg() *Value {
	lv := e.state.Get(-1)
	e.state.Pop(1)

	return newValue(lv)
}

func (e *Engine) PushRet(v *Value) {
	e.state.Push(v.lval)
}

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

func (e *Engine) LuaTable() *Value {
	return newValue(e.state.NewTable())
}

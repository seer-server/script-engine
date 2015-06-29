// Copyright (c) 2015 tree-server contributors

package lua

import (
	"fmt"

	"github.com/layeh/gopher-luar"
	glua "github.com/yuin/gopher-lua"
)

// Sandbox provides a way to define a script that loads in a secure environment
// (specificed by Script) and setting the variable that stores this secure
// script (EnvName).
type Sandbox struct {
	Script, EnvName string
}

// defaultSandbox is just a constant representative of the default sandbox
// options.
var defaultSandbox = Sandbox{
	Script:  secureSandboxScript,
	EnvName: defualtSandboxEnvName,
}

// Engine struct stores a pointer to a gluaLState providing a simplified API.
type Engine struct {
	state      *glua.LState
	Secure     bool
	sandbox    Sandbox
	securedFns map[string]struct{}
}

// ScriptFunction is a type alias for a function that receives an Engine and
// returns an int.
type ScriptFunction func(*Engine) int

// LuaTableMap interface to speed along the creation of table defining maps
// when creating Go modueles for use in Lua.
type LuaTableMap map[string]interface{}

// NewEngine creates a new engine containing a new lua.LState.
func NewEngine() *Engine {
	return &Engine{
		state:      glua.NewState(),
		Secure:     false,
		sandbox:    defaultSandbox,
		securedFns: make(map[string]struct{}),
	}
}

// NewSecureEngine creates a secure engine that will secure each function before
// it's called.
func NewSecureEngine() (*Engine, error) {
	engine := NewEngine()
	engine.Secure = true
	if err := engine.initiateKnockdown(); err != nil {
		return nil, err
	}

	return engine, nil
}

// NewCustomSecureEngine creates a new secure engine with the custom Sandbox
// provided. This will allow for cusotm security settings.
func NewCustomSecureEngine(sandbox Sandbox) (*Engine, error) {
	engine := NewEngine()
	engine.Secure = true
	engine.sandbox = sandbox
	if err := engine.initiateKnockdown(); err != nil {
		return nil, err
	}

	return engine, nil
}

// initiateKnockdown runs the SecureScript of the engine, this allows for custom
// security settings.
func (e *Engine) initiateKnockdown() error {
	if err := e.LoadString(e.sandbox.Script); err != nil {
		return err
	}

	return nil
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
func (e *Engine) SetGlobal(name string, val interface{}) {
	v := e.ValueFor(val)

	e.state.SetGlobal(name, v.lval)
}

// GetGlobal returns the value associated with the given name, or LuaNil
func (e *Engine) GetGlobal(name string) *Value {
	lv := e.state.GetGlobal(name)

	return newValue(lv)
}

// SetField applies the value to the given table associated with the given
// key.
func (e *Engine) SetField(tbl *Value, key string, val interface{}) {
	v := e.ValueFor(val)
	e.state.SetField(tbl.lval, key, v.lval)
}

// RegisterFunc registers a Go function with the script. Using this method makes
// Go functions accessible through Lua scripts.
func (e *Engine) RegisterFunc(name string, fn interface{}) {
	var lfn glua.LValue
	if sf, ok := fn.(func(*Engine) int); ok {
		lfn = e.genScriptFunc(sf)
	} else {
		v := e.ValueFor(fn)
		lfn = v.lval
	}
	e.state.SetGlobal(name, lfn)
}

// RegisterModule takes the values given, maps them to a LuaTable and then
// preloads the module with the given name to be consumed in Lua code.
func (e *Engine) RegisterModule(name string, fields map[string]interface{}) *Value {
	table := e.NewTable()
	for key, val := range fields {
		if sf, ok := val.(func(*Engine) int); ok {
			table.RawSet(key, e.genScriptFunc(sf))
		} else {
			table.RawSet(key, e.ValueFor(val).lval)
		}
	}

	loader := func(l *glua.LState) int {
		l.Push(table.lval)

		return 1
	}
	e.state.PreloadModule(name, loader)

	return table
}

// PopArg returns the top value on the Lua stack.
// This method is used to get arguments given to a Go function from a Lua script.
// This method will return a Value pointer that can then be converted into
// an appropriate type.
func (e *Engine) PopArg() *Value {
	lv := e.state.Get(-1)
	e.state.Pop(1)
	val := newValue(lv)
	if val.isTable() {
		val.owner = e
	}

	return val
}

// PushRet pushes the given Value onto the Lua stack.
// Use this method when 'returning' values from a Go function called from a
// Lua script.
func (e *Engine) PushRet(val interface{}) {
	v := e.ValueFor(val)
	e.state.Push(v.lval)
}

// PopBool returns the top of the stack as an actual Go bool.
func (e *Engine) PopBool() bool {
	v := e.PopArg()

	return v.AsBool()
}

// PopFunction is an alias for PopArg, provided for readability when specifying
// the desired value from the top of the stack.
func (e *Engine) PopFunction() *Value {
	return e.PopArg()
}

// PopInt returns the top of the stack as an actual Go int.
func (e *Engine) PopInt() int {
	v := e.PopArg()
	i := int(v.AsNumber())

	return i
}

// PopInt64 returns the top of the stack as an actual Go int64.
func (e *Engine) PopInt64() int64 {
	v := e.PopArg()
	i := int64(v.AsNumber())

	return i
}

// PopFloat returns the top of the stack as an actual Go float.
func (e *Engine) PopFloat() float64 {
	v := e.PopArg()

	return v.AsFloat()
}

// PopNumber is an alias for PopArg, provided for readability when specifying
// the desired value from the top of the stack.
func (e *Engine) PopNumber() *Value {
	return e.PopArg()
}

// PopString returns the top of the stack as an actual Go string value.
func (e *Engine) PopString() string {
	v := e.PopArg()

	return v.AsString()
}

// PopTable is an alias for PopArg, provided for readability when specifying
// the desired value from the top of the stack.
func (e *Engine) PopTable() *Value {
	tbl := e.PopArg()
	tbl.owner = e

	return tbl
}

// PopInterface returns the top of the stack as an actual Go interface.
func (e *Engine) PopInterface() interface{} {
	v := e.PopArg()

	return v.Interface()
}

// Call allows for calling a method by name.
// The second parameter is the number of return values the function being
// called should return. These values will be returned in a slice of Value
// pointers.
func (e *Engine) Call(name string, retCount int, params ...interface{}) ([]*Value, error) {
	luaParams := make([]glua.LValue, len(params))
	for i, iface := range params {
		v := e.ValueFor(iface)
		luaParams[i] = v.lval
	}

	if _, ok := e.securedFns[name]; e.Secure && !ok {
		secureScript := fmt.Sprintf("setfenv(%s, %s)", name, e.sandbox.EnvName)
		if err := e.LoadString(secureScript); err != nil {
			return nil, err
		}
		e.securedFns[name] = struct{}{}
	}

	err := e.state.CallByParam(glua.P{
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

// DefineType creates a construtor with the given name that will generate the
// given type.
func (e *Engine) RegisterType(name string, val interface{}) {
	cons := luar.NewType(e.state, val)
	e.state.SetGlobal(name, cons)
}

// RegisterClass assigns a new type, but instead of creating it via "TypeName()"
// it provides a more OO way of creating the object "TypeName.new()" otherwise
// it's functionally equivalent to RegisterType.
func (e *Engine) RegisterClass(name string, val interface{}) {
	cons := luar.NewType(e.state, val)
	table := e.NewTable()
	table.RawSet("new", cons)
	e.state.SetGlobal(name, table.lval)
}

// RegisterClassWithCtor does the same thing as RegisterClass excep the new
// function is mapped to the constructor passed in.
func (e *Engine) RegisterClassWithCtor(name string, typ interface{}, cons interface{}) {
	luar.NewType(e.state, typ)
	lcons := e.ValueFor(cons)
	table := e.NewTable()
	table.RawSet("new", lcons)

	e.state.SetGlobal(name, table.lval)
}

// ValueFor takes a Go type and creates a lua equivalent Value for it.
func (e *Engine) ValueFor(val interface{}) *Value {
	if v, ok := val.(*Value); ok {
		return v
	} else {
		return newValue(luar.New(e.state, val))
	}
}

// NewTable creates and returns a new NewTable.
func (e *Engine) NewTable() *Value {
	tbl := newValue(e.state.NewTable())
	tbl.owner = e

	return tbl
}

// wrapScriptFunction turns a ScriptFunction into a lua.LGFunction
func (e *Engine) wrapScriptFunction(fn ScriptFunction) glua.LGFunction {
	return func(l *glua.LState) int {
		e := &Engine{state: l}

		return fn(e)
	}
}

// genScriptFunc will wrap a ScriptFunction with a function that gopher-lua
// expects to see when calling method from Lua.
func (e *Engine) genScriptFunc(fn ScriptFunction) *glua.LFunction {
	return e.state.NewFunction(e.wrapScriptFunction(fn))
}

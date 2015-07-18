package lua

import (
	"fmt"

	"github.com/layeh/gopher-luar"
	glua "github.com/yuin/gopher-lua"
)

// ValueError provides information about failed Value typecasts.
type ValueError string

// newValueError creates a new error explaining failure from a given type to an
// actual type.
func newValueError(exp string, v *Value) ValueError {
	return ValueError(fmt.Sprintf("expected %s, found \"%s\"", exp, v.lval.Type()))
}

// Implements the Error interface for ValueError
func (v ValueError) Error() string {
	return string(v)
}

// Value is a utility wrapper for lua.LValue that provies conveinient methods
// for casting.
type Value struct {
	lval  glua.LValue
	owner *Engine
}

// newValue constructs a new value from an LValue.
func newValue(val glua.LValue) *Value {
	return &Value{
		lval: val,
	}
}

// String makes Value conform to Stringer
func (v *Value) String() string {
	return v.lval.String()
}

// AsString returns the LValue as a Go string
func (v *Value) AsString() string {
	return glua.LVAsString(v.lval)
}

// AsFloat returns the LValue as a Go float64.
// This method will try to convert the Lua value to a number if possible, if
// not then LuaNumber(0) is returned.
func (v *Value) AsFloat() float64 {
	return float64(glua.LVAsNumber(v.lval))
}

// AsNumber is an alias for AsFloat (Lua calls them "numbers")
func (v *Value) AsNumber() float64 {
	return v.AsFloat()
}

// AsBool returns the Lua boolean representation for an object (this works for
// non bool Values)
func (v *Value) AsBool() bool {
	return glua.LVAsBool(v.lval)
}

// IsNil will only return true if the Value wraps LNil.
func (v *Value) IsNil() bool {
	return v.lval.Type() == glua.LTNil
}

// IsFalse is similar to AsBool except it returns if the Lua value would be
// considered false in Lua.
func (v *Value) IsFalse() bool {
	return glua.LVIsFalse(v.lval)
}

// IsTrue returns whether or not this is a truthy value or not.
func (v *Value) IsTrue() bool {
	return !v.IsFalse()
}

// The following methods allow for type detection

// IsNumber returns true if the stored value is a numeric value.
func (v *Value) IsNumber() bool {
	return v.lval.Type() == glua.LTNumber
}

// IsBool returns true if the stored value is a boolean value.
func (v *Value) IsBool() bool {
	return v.lval.Type() == glua.LTBool
}

// IsFunction returns true if the stored value is a function.
func (v *Value) IsFunction() bool {
	return v.lval.Type() == glua.LTFunction
}

// IsString returns true if the stored value is a string.
func (v *Value) IsString() bool {
	return v.lval.Type() == glua.LTString
}

// IsTable returns true if the stored value is a table.
func (v *Value) IsTable() bool {
	return v.lval.Type() == glua.LTTable
}

// The following methods allow LTable values to be modified through Go.

// isTable returns a bool if the Value is an LTable.
func (v *Value) isTable() bool {
	return v.lval.Type() == glua.LTTable
}

// asTable converts the Value into an LTable.
func (v *Value) asTable() (t *glua.LTable) {
	t, _ = v.lval.(*glua.LTable)

	return
}

// isUserData returns a bool if the Value is an LUserData
func (v *Value) isUserData() bool {
	return v.lval.Type() == glua.LTUserData
}

// asUserData converts the Value into an LUserData
func (v *Value) asUserData() (t *glua.LUserData) {
	t, _ = v.lval.(*glua.LUserData)

	return
}

// Append maps to lua.LTable.Append
func (v *Value) Append(value interface{}) {
	if v.isTable() {
		val := getLValue(v.owner, value)

		t := v.asTable()
		t.Append(val)
	}
}

// ForEach maps to lua.LTable.ForEach
func (v *Value) ForEach(cb func(*Value, *Value)) {
	if v.isTable() {
		actualCb := func(key glua.LValue, val glua.LValue) {
			cb(newValue(key), newValue(val))
		}
		t := v.asTable()
		t.ForEach(actualCb)
	}
}

// Insert maps to lua.LTable.Insert
func (v *Value) Insert(i int, value interface{}) {
	if v.isTable() {
		val := getLValue(v.owner, value)

		t := v.asTable()
		t.Insert(i, val)
	}
}

// Len maps to lua.LTable.Len
func (v *Value) Len() int {
	if v.isTable() {
		t := v.asTable()

		return t.Len()
	}

	return -1
}

// MaxN maps to lua.LTable.MaxN
func (v *Value) MaxN() int {
	if v.isTable() {
		t := v.asTable()

		return t.MaxN()
	}

	return 0
}

// Next maps to lua.LTable.Next
func (v *Value) Next(key interface{}) (*Value, *Value) {
	if v.isTable() {
		val := getLValue(v.owner, key)

		t := v.asTable()
		v1, v2 := t.Next(val)

		return newValue(v1), newValue(v2)
	}

	return Nil, Nil
}

// Remove maps to lua.LTable.Remove
func (v *Value) Remove(pos int) *Value {
	if v.isTable() {
		t := v.asTable()
		ret := t.Remove(pos)

		return newValue(ret)
	}

	return Nil
}

// Helper method for Set and RawSet
func getLValue(e *Engine, item interface{}) glua.LValue {
	switch val := item.(type) {
	case (*Value):
		return val.lval
	case glua.LValue:
		return val
	}

	if e != nil {
		return luar.New(e.state, item)
	}

	return glua.LNil
}

// Set sets the value of a given key on the table, this method checks for
// validity of array keys and handles them accordingly.
func (v *Value) Set(goKey interface{}, val interface{}) {
	if v.isTable() {
		key := getLValue(v.owner, goKey)
		lval := getLValue(v.owner, val)

		v.asTable().RawSet(key, lval)
	}
}

// RawSet bypasses any checks for key existence and sets the value onto the
// table with the given key.
func (v *Value) RawSet(goKey interface{}, val interface{}) {
	if v.isTable() {
		key := getLValue(v.owner, goKey)
		lval := getLValue(v.owner, val)

		v.asTable().RawSetH(key, lval)
	}
}

// The following provde methods for LUserData

// Interface returns the value of the LUserData
func (v *Value) Interface() interface{} {
	if v.isUserData() {
		t := v.asUserData()

		return t.Value
	}

	return nil
}

// The following provide LFunction methods on Value

// FuncLocalName is a function that returns the local name of a LFunction type
// if this Value objects holds an LFunction.
func (v *Value) FuncLocalName(regno, pc int) (string, bool) {
	if f, ok := v.lval.(*glua.LFunction); ok {
		return f.LocalName(regno, pc)
	}

	return "", false
}

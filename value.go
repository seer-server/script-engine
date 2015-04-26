package engine

import (
	"github.com/yuin/gopher-lua"

	"fmt"
)

// Provides information about failed Value typecasts.
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
	lval lua.LValue
}

// newValue constructs a new value from an LValue.
func newValue(val lua.LValue) *Value {
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
	return lua.LVAsString(v.lval)
}

// AsFloat returns the LValue as a Go float64.
// This method will try to convert the Lua value to a number if possible, if
// not then LuaNumber(0) is returned.
func (v *Value) AsFloat() float64 {
	return float64(lua.LVAsNumber(v.lval))
}

// AsNumber is an alias for AsFloat (Lua calls them "numbers")
func (v *Value) AsNumber() (float64, error) {
	return v.AsFloat()
}

// AsBool returns the Lua boolean representation for an object (this works for
// non bool Values)
func (v *Value) AsBool() bool {
	return lua.LVAsBool(v.lval)
}

// IsNil will only return true if the Value wraps LNil.
func (v *Value) IsNil() bool {
	return v.lval.Type() == lua.LTNil
}

// IsFalase is similar to AsBool except it returns if the Lua value would be
// considered false in Lua.
func (v *Value) IsFalse() bool {
	return lua.LVIsFalse(v.lval)
}

// The following methods allow LTable values to be modified through Go.

// isTable returns a bool if the Value is an LTable.
func (v *Value) isTable() bool {
	return v.lval.Type() == lua.LTTable
}

// asTable converts the Value into an LTable.
func (v *Value) asTable() (t *lua.LTable) {
	t, _ = v.lval.(*lua.LTable)

	return
}

// TableAppend maps to lua.LTable.Append
func (v *Value) TableAppend(value *Value) {
	if v.isTable() {
		t := v.asTable()
		t.Append(value.lval)
	}
}

// TableForEach maps to lua.LTable.ForEach
func (v *Value) TableForEach(cb func(*Value, *Value)) {
	if v.isTable() {
		actualCb := func(key lua.LValue, val lua.LValue) {
			cb(newValue(key), newValue(val))
		}
		t := v.asTable()
		t.ForEach(actualCb)
	}
}

// TableInsert maps to lua.LTable.Insert
func (v *Value) TableInsert(i int, value *Value) {
	if v.isTable() {
		t := v.asTable()
		t.Insert(i, value.lval)
	}
}

// TableLen maps to lua.LTable.Len
func (v *Value) TableLen() int {
	if v.isTable() {
		t := v.asTable()

		return t.Len()
	}

	return -1
}

// TableMaxN maps to lua.LTable.MaxN
func (v *Value) TableMaxN() int {
	if v.isTable() {
		t := v.asTable()

		return t.MaxN()
	}

	return 0
}

// TableNext maps to lua.LTable.Next
func (v *Value) TableNext(key *Value) (*Value, *Value) {
	if v.isTable() {
		t := v.asTable()
		v1, v2 := t.Next(key.lval)

		return newValue(v1), newValue(v2)
	}

	return LuaNil(), LuaNil()
}

// TableRawGet maps to lua.LTable.RawGet
func (v *Value) TableRawGet(key *Value) *Value {
	if v.isTable() {
		t := v.asTable()
		ret := t.RawGet(key.lval)

		return newValue(ret)
	}

	return LuaNil()
}

// TableRawGetH maps to lua.LTable.RawGetH
func (v *Value) TableRawGetH(key *Value) *Value {
	if v.isTable() {
		t := v.asTable()
		ret := t.RawGetH(key.lval)

		return newValue(ret)
	}

	return LuaNil()
}

// TableRawGetInt maps to lua.LTable.RawGetInt
func (v *Value) TableRawGetInt(key int) *Value {
	if v.isTable() {
		t := v.asTable()
		ret := t.RawGetInt(key)

		return newValue(ret)
	}

	return LuaNil()
}

// TableRawSet maps to lua.LTable.RawSet
func (v *Value) TableRawSet(key *Value, value *Value) {
	if v.isTable() {
		t := v.asTable()
		t.RawSet(key.lval, value.lval)
	}
}

// TableRawSetH maps to lua.LTable.RawSetH
func (v *Value) TableRawSetH(key *Value, value *Value) {
	if v.isTable() {
		t := v.asTable()
		t.RawSetH(key.lval, value.lval)
	}
}

// TableRawSetInt maps to lua.LTable.RawSetInt
func (v *Value) TableRawSetInt(key int, value *Value) {
	if v.isTable() {
		t := v.asTable()
		t.RawSetInt(key, value.lval)
	}
}

// TableRemove maps to lua.LTable.Remove
func (v *Value) TableRemove(pos int) *Value {
	if v.isTable() {
		t := v.asTable()
		ret := t.Remove(pos)

		return newValue(ret)
	}

	return LuaNil()
}

// The following provide LFunction methods on Value

func (v *Value) FuncLocalName(regno, pc int) (string, bool) {
	if f, ok := v.lval.(lua.LFunction); ok {
		return f.LocalName(regno, pc)
	} else {
		return "", false
	}
}

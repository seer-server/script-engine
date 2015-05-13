package engine

import (
	"testing"
)

func TestLoadStringDoesNotFail(t *testing.T) {
	e := NewEngine()
	defer e.Close()
	if err := e.LoadString("local a = 1"); err != nil {
		t.Error(e)
	}
}

func TestCanCallMethod(t *testing.T) {
	e := NewEngine()
	defer e.Close()

	_ = e.LoadString(`
		function double(x)
			return x * 2
		end`)

	ret, err := e.Call("double", 1, LuaNumber(10))
	if err != nil {
		t.Error(err)

		return
	}

	if len(ret) != 1 {
		t.Errorf("Expected %d return values, but found %d instead", 1, len(ret))
	}

	n := ret[0].AsNumber()

	exp := float64(20)
	if n != exp {
		t.Errorf("Expected return value of %f but found %f", exp, n)
	}
}

func TestCanLoadFromFile(t *testing.T) {
	e := NewEngine()
	defer e.Close()

	err := e.LoadFile("lua_test.lua")
	if err != nil {
		t.Error(err)

		return
	}

	ret, err := e.Call("give_me_one", 1)
	if err != nil {
		t.Error(err)

		return
	}

	n := ret[0].AsNumber()
	exp := float64(1)
	if n != exp {
		t.Errorf("Expected return value of %f but got %f", exp, n)
	}
}

func TestCallingGoFromLua(t *testing.T) {
	e := NewEngine()
	defer e.Close()

	dbl := func(e *Engine) int {
		n := e.PopArg().AsNumber()
		e.PushRet(LuaNumber(n * 2))

		return 1
	}

	e.RegisterFunc("double", dbl)

	e.LoadString(`
		function test(x)
			return double(x)
		end`)
	ret, err := e.Call("test", 1, LuaNumber(10))
	if err != nil {
		t.Error(e)

		return
	}

	n := ret[0].AsNumber()
	exp := float64(20)
	if n != exp {
		t.Errorf("Expected return value of %f but found %f", exp, n)
	}
}

func TestLoadingModules(t *testing.T) {
	e := NewEngine()
	defer e.Close()

	fnMap := ScriptFnMap{
		"double": func(e *Engine) int {
			x := e.PopArg().AsNumber()
			e.PushRet(LuaNumber(x * 2))

			return 1
		},
		"hello": func(e *Engine) int {
			name := e.PopArg().AsString()
			e.PushRet(LuaString("Hello, " + name + "!"))

			return 1
		},
	}

	loader := func(e *Engine) *Value {
		return e.GenerateModule(fnMap)
	}

	e.RegisterModule("test_mod", loader)
	err := e.LoadString(`
		local test = require("test_mod")

		function test_double(x)
			return test.double(x)
		end

		function test_hello(name)
			return test.hello(name)
		end`)
	if err != nil {
		t.Error(err)

		return
	}

	ret, err := e.Call("test_double", 1, LuaNumber(10))
	if err != nil {
		t.Error(err)

		return
	}

	n := ret[0].AsNumber()
	exp := float64(20)
	if n != exp {
		t.Errorf("Expected return value %f got %f", exp, n)

		return
	}

	ret, err = e.Call("test_hello", 1, LuaString("World"))
	if err != nil {
		t.Error(err)

		return
	}

	s := ret[0].AsString()
	expStr := "Hello, World!"
	if s != expStr {
		t.Errorf("Expected return value of %q but got %q", expStr, s)

		return
	}
}

func TestEngineValueFor(t *testing.T) {
	e := NewEngine()
	sexp := "This is a String"
	sval := e.ValueFor(sexp)
	if sexp != sval.AsString() {
		t.Errorf("Expected %q but got %q", sexp, sval.AsString())

		return
	}

	var nexp float64 = 10.0
	nval := e.ValueFor(nexp)
	if nexp != nval.AsNumber() {
		t.Errorf("Expected %f but got %f", nexp, nval.AsNumber())

		return
	}

	bexp := true
	bval := e.ValueFor(bexp)
	if bexp != bval.AsBool() {
		t.Errorf("Expected true, found %b", bval.AsBool())

		return
	}

	empty := e.ValueFor(nil)
	if !empty.IsNil() {
		t.Error("Expected nil, but it wasn't")

		return
	}

	osval := e.ValueFor(sval)
	if osval != sval {
		t.Error("Expected given value pointer to be returned, but it wasn't.")

		return
	}
}

func TestTypeConstructor(t *testing.T) {
	type Song struct {
		Title, Artist string
	}

	e := NewEngine()
	e.DefineType("Song", Song{})
	e.LoadString(`
		local s = Song()
	    s.Title = "Some Song Name"
	    s.Artist = "Some Awesome Artist"

		function test_song()
		  return s.Title .. " - " .. s.Artist
		end

		function get_song()
		  return s
		end`)
	ret, err := e.Call("test_song", 1)
	if err != nil {
		t.Error(err)

		return
	}

	s := ret[0].AsString()
	exp := "Some Song Name - Some Awesome Artist"
	if s != exp {
		t.Errorf("Expected %q but returned %q", exp, s)
	}

	ret, err = e.Call("get_song", 1)
	if err != nil {
		t.Error(err)

		return
	}

	iface := ret[0].Interface()
	if _, ok := ret[0].Interface().(*Song); !ok {
		t.Errorf("Expected a Song, but received %T (%v)", iface, iface)

		return
	}
}

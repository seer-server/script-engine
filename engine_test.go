package engine

import (
	"testing"
)

func TestLoadStringDoesNotFail(t *testing.T) {
	e := NewEngine()
	defer e.Close()
	if err := e.LoadString("print(\"Hello\")"); err != nil {
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
	}

	ret, err := e.Call("test_double", 1, LuaNumber(10))
	if err != nil {
		t.Error(0)
	}

	n := ret[0].AsNumber()
	exp := float64(20)
	if n != exp {
		t.Errorf("Expected return value %f got %f", exp, n)
	}

	ret, err = e.Call("test_hello", 1, LuaString("World"))
	if err != nil {
		t.Error(err)
	}

	s := ret[0].AsString()
	expStr := "Hello, World!"
	if s != expStr {
		t.Errorf("Expected return value of %q but got %q", expStr, s)
	}
}

// Code for testing loading module

func loader(e *Engine) *Value {
	return e.GenerateModule(fnMap)
}

var fnMap = ScriptFnMap{
	"double": double,
	"hello":  hello,
}

func double(e *Engine) int {
	x := e.PopArg().AsNumber()
	e.PushRet(LuaNumber(x * 2))

	return 1
}

func hello(e *Engine) int {
	name := e.PopArg().AsString()
	e.PushRet(LuaString("Hello, " + name + "!"))

	return 1
}

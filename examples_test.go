package engine

import (
	"fmt"
)

func ExampleEngine() {
	e := NewEngine() // create a new engine
	defer e.Close()  // make sure to close the engine

	err := e.LoadFile("some_lua_file.lua")
	if err != nil {
		panic(err)
	}

	ret, err := e.Call("double_this_number", 1, LuaNumber(10))
	if err != nil {
		panic(err)
	}

	n := ret[0].AsNumber()
	fmt.Println(n) // 20.0000000
}

func ExampleEngine_Call_simple() {
	e := NewEngine()
	defer e.Close()

	e.Call("some_method", 0)
}

func ExampleEngine_Call_multiple() {
	e := NewEngine()
	defer e.Close()

	ret, _ := e.Call("swap_these_numbers", 2, LuaNumber(10), LuaNumber(20))

	a, b := ret[0].AsNumber(), ret[1].AsNumber()
	fmt.Println(a) // 20.0000000
	fmt.Println(b) // 10.0000000
}

func ExampleEngine_RegisterModule() {
	fnMap := ScriptFnMap{
		"double": func(e *Engine) int {
			i := e.PopArg().AsNumber()
			e.PushRet(LuaNumber(i * 2))

			return 1
		},
		"swap": func(e *Engine) int {
			a := e.PopArg()
			b := e.PopArg()
			e.PushRet(b)
			e.PushRet(a)

			return 2
		},
	}

	loader := func(e *Engine) *Value {
		return e.GenerateModule(fnMap)
	}

	e := NewEngine()
	defer e.Close()

	e.RegisterModule("example", loader)
	e.LoadString(`
		local e = require("example")

		e.double(10) -- 20
		e.swap(1, 2) -- 2, 1
	`)
}

func ExampleEngine_RegisterFunc() {
	e := NewEngine()
	defer e.Close()

	e.RegisterFunc("double", func(e *Engine) int {
		n := e.PopArg().AsNumber()
		e.PushRet(LuaNumber(e * 2))

		return 1
	})
	e.LoadString("double(10) -- 20")
}

package lua

import (
	"fmt"
	"math"
)

func ExampleEngine() {
	e := NewEngine() // create a new engine
	defer e.Close()  // make sure to close the engine

	err := e.LoadFile("some_lua_file.lua")
	if err != nil {
		panic(err)
	}

	// Raw values can be passed through to the function. They'll be converted
	// appropriately.
	ret, err := e.Call("double_this_number", 1, 10)
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

	ret, _ := e.Call("swap_these_numbers", 2, 10, 20)

	a, b := ret[0].AsNumber(), ret[1].AsNumber()
	fmt.Println(a) // 20.0000000
	fmt.Println(b) // 10.0000000
}

func ExampleEngine_RegisterModule() {
	e := NewEngine()
	defer e.Close()

	e.RegisterModule("example", LuaTableMap{
		// Using a non-standard function definition to provide concise means
		"double": func(x float64) float64 {
			return x * 2
		},
		// Using a standard ScriptFunction definition to provide dynamic
		// capabilities
		"swap": func(e *Engine) int {
			b := e.PopArg()
			a := e.PopArg()
			e.PushRet(b)
			e.PushRet(a)

			return 2
		},
	})
	e.LoadString(`
		local e = require("example")

		e.double(10) -- 20
		e.swap(1, 2) -- 2, 1
	`)
}

func ExampleEngine_RegisterFunc() {
	e := NewEngine()
	defer e.Close()

	// Add script functions
	e.RegisterFunc("double", func(e *Engine) int {
		n := e.PopFloat()
		e.PushRet(n * 2)

		return 1
	})
	e.LoadString("double(10) -- 20")

	// Or add standard Go functions, all stack communication done for you
	e.RegisterFunc("power", func(x, y float64) float64 {
		return math.Pow(x, y)
	})
	e.LoadString("power(2, 2) -- 4")
}

func ExampleEngine_RegisterType() {
	type Song struct {
		Artist, Title string
	}

	e := NewEngine()
	defer e.Close()

	e.RegisterType("Song", Song{})
	e.LoadString(`
		local s = Song() -- create new structs

		s.Artist = "Some artist" -- set values on structs
		s.title = "Some title" -- keys and methods can be "lowercased"
	`)
}

func ExampleEngine_RegisterClass() {
	type Song struct {
		Artist, Title string
	}

	e := NewEngine()
	defer e.Close()

	e.RegisterClass("Song", Song{})
	e.LoadString(`
		local s = Song.new() -- create new structs with a 'new' function

		s.Artist = "Some artist" -- set values on structs
		s.title = "Some title" -- keys and methods can be "lowercased"
	`)
}

func ExampleEngine_RegisterClassWithCtor() {
	type Song struct {
		Artist, Title string
	}

	newSong := func(n, a string) *Song {
		return &Song{a, n}
	}

	e := NewEngine()
	defer e.Close()

	e.RegisterClassWithCtor("Song", Song{}, newSong)
	e.LoadString(`
		-- create new structs with custom constructors!
		local s = Song.new("Some title", "Some artist")
	`)
}

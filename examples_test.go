package engine

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

func ExampleCall_simple() {
	e.Call("some_method", 0)
}

func ExampleCall_multiple_returns() {
	ret, _ := e.Call("swap_these_numbers", 2, LuaNumber(10), LuaNumber(20))

	a, b := ret[0].AsNumber(), ret[1].AsNumber()
	fmt.Println(a) // 20.0000000
	fmt.Println(b) // 10.0000000
}

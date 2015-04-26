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

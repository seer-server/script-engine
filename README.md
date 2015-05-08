# Script Engine [![GoDoc](http://godoc.org/github.com/tree-server/script-engine?status.svg)](http://godoc.org/github.com/tree-server/script-engine)

This provides a smaller, simpler wrapper around the features of [Gopher Lua](http://github.com/yuin/gopher-lua).

Usage of the ScriptEngine is a through a smaller more simple API than dealing with al the raw Lua states. This will make integrating Lua scripts with Go aspects of the softare much quicker.

Simple example, load a string script and call the function. The second argument to `Call` is an int, representing the number of return values you
expect in return from the call.

```go
eng := engine.NewEngine()
defer eng.Close()
eng.LoadScript(`
  function double(x)
    return x * 2
  end
`)
// Call returns a slice of *Value, LuaNumber converts a number to 
// the appropriate type to be treated as a number in the 
// script
ret, _ := eng.Call("double", 1, LuaNumber(10))
// AsNumber() converts a *Value into a float64, this method always succeeds
// which means you should know that a number was returned. This "always 
// succeeds" mentatility stems from Lua's loose type system.
f := ret[0].AsNumber()
fmt.Println(f) // => 20.0000000
```

Calling Go functions from Lua. A Go function that can be registered for the script to access it is of type `engine.ScriptFunction` which is `func(*Engine) int`. The return value from this function is the number of values you've pusehd on the stack.

```go
eng := engine.NewEngine()
defer eng.Close()
eng.RegisterFunction("double", func(e *Engine) int {
        f, _ := e.PopArg().AsNumber()
        e.PushRet(LuaNumber(n * 2))

        return 1
})
eng.LoadScript(`
  function test(x)
    return double(x)
  end
`)
ret, _ := eng.Call("test", 1, LuaNumber(10))
f := ret[0].AsNumber()
fmt.Println(f) // => 20.0000000
```

# Thanks

I have to thank [Yusuke Inuzuka](http://github.com/yuin) for making one of my absolute favority Go -> Lua libraries that are currently avialable. It's easy to understand, pure Go and is generally just a pleasure to work with.

# Contributors

Brandon Buck

# License

MIT

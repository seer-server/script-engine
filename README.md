# Script Engine [![GoDoc](http://godoc.org/github.com/seer-server/script-engine?status.svg)](http://godoc.org/github.com/seer-server/script-engine)

This provides a smaller, simpler wrapper around the features of [gopher-lua](http://github.com/yuin/gopher-lua) and includes features from [gopher-luar](http://github.com/layeh/gopher-luar).

Usage of the ScriptEngine is through a smaller and more simple API than dealing with the row Lua processes. This will make integrating Lua scripts with Go aspects of the softare much quicker and less error prone.

### Calling Lua from Go

Simple example, load a string script and call the function. The second argument to `Call` is an int, representing the number of return values you
expect in return from the call.

```go
eng := lua.NewEngine()
defer eng.Close()
eng.LoadScript(`
  function double(x)
    return x * 2
  end
`)

// Call returns a slice of *Value, LuaNumber converts a number to 
// the appropriate type to be treated as a number in the 
// script
ret, _ := eng.Call("double", 1, lua.LuaNumber(10))

// AsNumber() converts a *Value into a float64, this method always succeeds
// which means you should know that a number was returned. This "always 
// succeeds" mentatility stems from Lua's loose type system.
f := ret[0].AsNumber()
fmt.Println(f) // => 20.0000000
```

### Calling Go from Lua

Calling Go functions from Lua. A Go function that can be registered for the script to access it is of type `lua.ScriptFunction` which is `func(*Engine) int`. The return value from this function is the number of values you've pusehd on the stack.

```go
eng := lua.NewEngine()
defer eng.Close()
eng.RegisterFunc("double", func(e *Engine) int {
        f, _ := e.PopArg().AsNumber()
        e.PushRet(LuaNumber(n * 2))

        return 1
})
eng.LoadScript(`
  function test(x)
    return double(x)
  end
`)
ret, _ := eng.Call("test", 1, lua.LuaNumber(10))
f := ret[0].AsNumber()
fmt.Println(f) // => 20.0000000
```

Thanks to the new [gopher-luar](http://github.com/layeh/gopher-luar) integration this previous example is even easier to implement.

```go
eng := lua.NewEngine()
defer eng.Close()
eng.RegisterFunc("double", func(x float64) float64 {
        return x * 2
})
eng.LoadScript(`
  function test(x)
    return double(x)
  end
`)
ret, _ := eng.Call("test", 1, 10) // no more lua.LuaNumber(10) needed!
f := ret[0].AsNumber()
fmt.Println(f) // => 20.000000
```

### User Data

Again, thanks to the power of gopher-luar we can easily pass in Go types without worry about boilerplate (and a lot of it, at that).

```go
// Accessible fields should be public
type Person struct {
        Name string
        Age  int
}

// Objects created in Lua are pointers
func (p *Person) String() string {
        return fmt.Sprintf("%s (%d)", p.Name, p.Age)
}

eng := lua.NewEngine()
defer eng.Close()
eng.RegisterType("Person", Person{})
eng.LoadString(`
  function testTypes()
    p = Person()
    p.name = "Brandon"
    p.age = 28

    return p.string()
  end
`)
ret, _ := eng.Call("testTypes", 1)
fmt.Println(ret[0].AsString()) // => Brandon (28)
```

# Thanks

I have to thank [Yusuke Inuzuka](http://github.com/yuin) for making one of my absolute favority Go -> Lua libraries that are currently avialable. It's easy to understand, pure Go and is generally just a pleasure to work with.

I also thank [layeh](http://github.com/layeh) (organization) for their development on gopher-luar which has boosted the simplicity and capabilites of this engine significantly.

# Contributors

[Brandon Buck](http://github.com/bbuck)

# License

MIT

package engine

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

var (
	fibCode = `
		local function fib(n)
			if n < 2 then
				return n
			end
			return fib(n - 2) + fib(n - 1)
		end

		function call_fib(value)
			return fib(value)
		end
	`
	addCode = `
		local function add(a, b)
			return a + b
		end

		function call_add(a, b)
			return add(a, b)
		end
	`
)

func Benchmark_EngineFib5(b *testing.B) {
	fibValue := 5
	for i := 0; i < b.N; i++ {
		e := NewEngine()
		defer e.Close()
		e.LoadString(fibCode)
		ret, _ := e.Call("call_fib", 1, fibValue)
		_ = ret[0].AsNumber()
	}
}

func Benchmark_RawGopherLuaFib5(b *testing.B) {
	fibValue := 5
	for i := 0; i < b.N; i++ {
		L := lua.NewState()
		defer L.Close()
		L.DoString(fibCode)
		lnum := lua.LNumber(fibValue)
		L.CallByParam(lua.P{
			Fn:      L.GetGlobal("call_fib"),
			NRet:    1,
			Protect: true,
		}, lnum)
		ret := L.Get(-1)
		L.Pop(1)
		_ = lua.LVAsNumber(ret)
	}
}

func Benchmark_EngineFib30(b *testing.B) {
	fibValue := 30
	for i := 0; i < b.N; i++ {
		e := NewEngine()
		defer e.Close()
		e.LoadString(fibCode)
		ret, _ := e.Call("call_fib", 1, fibValue)
		_ = ret[0].AsNumber()
	}
}

func Benchmark_RawGopherLuaFib30(b *testing.B) {
	fibValue := 30
	for i := 0; i < b.N; i++ {
		L := lua.NewState()
		defer L.Close()
		L.DoString(fibCode)
		lnum := lua.LNumber(fibValue)
		L.CallByParam(lua.P{
			Fn:      L.GetGlobal("call_fib"),
			NRet:    1,
			Protect: true,
		}, lnum)
		ret := L.Get(-1)
		L.Pop(1)
		_ = lua.LVAsNumber(ret)
	}
}

func Benchmark_EngineAdd(b *testing.B) {
	a, c := 1870183.0, 109899.0
	for i := 0; i < b.N; i++ {
		e := NewEngine()
		defer e.Close()
		e.LoadString(addCode)
		ret, _ := e.Call("call_add", 1, a, c)
		_ = ret[0].AsNumber()
	}
}

func Benchmark_RawGopherLuaAdd(b *testing.B) {
	a, c := 1870183.0, 109899.0
	for i := 0; i < b.N; i++ {
		L := lua.NewState()
		defer L.Close()
		L.DoString(addCode)
		la := lua.LNumber(a)
		lc := lua.LNumber(c)
		L.CallByParam(lua.P{
			Fn:      L.GetGlobal("call_add"),
			NRet:    1,
			Protect: true,
		}, la, lc)
		ret := L.Get(-1)
		L.Pop(1)
		_ = lua.LVAsNumber(ret)
	}
}
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

func Benchmark_EngineGoToLua(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := NewEngine()
		defer e.Close()
		e.RegisterFunc("add", func(a, b float64) float64 {
			return a + b
		})
		e.LoadString("num = add(1, 2)")
		_ = e.GetGlobal("num").AsNumber()
	}
}

func Benchmark_GopherLuaGoToLua(b *testing.B) {
	for i := 0; i < b.N; i++ {
		L := lua.NewState()
		defer L.Close()
		L.SetGlobal("add", L.NewFunction(func(l *lua.LState) int {
			la := l.ToNumber(-2)
			lb := l.ToNumber(-1)
			l.Pop(2)

			a, b := float64(la), float64(lb)
			L.Push(lua.LNumber(a + b))

			return 1
		}))
		L.DoString("num = add(1, 2)")
		ln := L.GetGlobal("num")
		_ = float64(lua.LVAsNumber(ln))
	}
}

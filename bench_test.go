package gisp_test

import (
	"encoding/json"
	"testing"

	"github.com/ysmood/gisp"
	lua "github.com/yuin/gopher-lua"
)

func Add(L *lua.LState) int {
	a := L.ToInt(1)            /* get argument */
	b := L.ToInt(2)            /* get argument */
	L.Push(lua.LNumber(a + b)) /* push result */
	return 1                   /* number of results */
}

func BenchmarkLua(b *testing.B) {
	L := lua.NewState()
	defer L.Close()
	L.SetGlobal("add", L.NewFunction(Add))

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		L.DoString("add(1.1,1.2)")
	}
}

func BenchmarkAST(b *testing.B) {
	code := []byte(`["+", 1.1, 1.2]`)
	var ast interface{}
	json.Unmarshal(code, &ast)

	sandbox := gisp.Sandbox{
		"+": func(ctx *gisp.Context) interface{} {
			return ctx.ArgNum(1) + ctx.ArgNum(2)
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gisp.Run(&gisp.Context{
			AST:     ast,
			Sandbox: sandbox,
		})
	}
}

func BenchmarkAST_liftPanic(b *testing.B) {
	code := []byte(`["+", 1, 1]`)
	var ast interface{}
	json.Unmarshal(code, &ast)

	sandbox := gisp.Sandbox{
		"+": func(ctx *gisp.Context) interface{} {
			return ctx.ArgNum(1) + ctx.ArgNum(2)
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gisp.Run(&gisp.Context{
			AST:         ast,
			Sandbox:     sandbox,
			IsLiftPanic: true,
		})
	}
}

func BenchmarkJSON(b *testing.B) {
	sandbox := gisp.Sandbox{
		"+": func(ctx *gisp.Context) interface{} {
			return ctx.ArgNum(1) + ctx.ArgNum(2)
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gisp.RunJSON([]byte(`["+", ["+", 1, 1], ["+", 1, 1]]`), &gisp.Context{
			Sandbox: sandbox,
		})
	}
}

func BenchmarkJSONBase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		code := []byte(`["+", ["+", 1, 1], ["+", 1, 1]]`)
		var ast interface{}
		json.Unmarshal(code, &ast)
	}
}

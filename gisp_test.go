package gisp_test

import (
	"fmt"
	"testing"

	"github.com/a8m/djson"
	"github.com/stretchr/testify/assert"
	"github.com/ysmood/gisp"
	"github.com/ysmood/gisp/lib"
)

func TestReadmeExample(t *testing.T) {
	code := `["+", 1, ["*", 2, 2]]`

	out, _ := gisp.RunJSON(code, &gisp.Context{
		Sandbox: gisp.New(gisp.Box{
			"+": lib.Add,
			"*": lib.Multiply,
		}),
	})

	fmt.Println(out) // 5

	assert.Equal(t, float64(5), out)
}
func TestReturnFn(t *testing.T) {
	sandbox := gisp.New(gisp.Box{
		"foo": func(ctx *gisp.Context) interface{} {
			return func(ctx *gisp.Context) interface{} {
				return ctx.ArgNum(1) + ctx.ArgNum(2)
			}
		},
	})

	out, _ := gisp.RunJSON(`[["foo"], 1, 2]`, &gisp.Context{
		Sandbox: sandbox,
	})

	assert.Equal(t, float64(3), out)
}

func TestEmpty(t *testing.T) {
	out, _ := gisp.RunJSON(`[]`, &gisp.Context{
		IsLiftPanic: true,
		Sandbox:     gisp.New(gisp.Box{}),
	})
	assert.Equal(t, nil, out)
}

func TestStr(t *testing.T) {
	sandbox := gisp.New(gisp.Box{})

	out, _ := gisp.RunJSON(`"foo"`, &gisp.Context{
		Sandbox: sandbox,
	})

	assert.Equal(t, "foo", out)
}

func TestVal(t *testing.T) {
	sandbox := gisp.New(gisp.Box{
		"test": "ok",
	})

	out, _ := gisp.RunJSON(`["test"]`, &gisp.Context{
		Sandbox: sandbox,
	})

	assert.Equal(t, "ok", out)
}

func TestPreRun(t *testing.T) {
	sandbox := gisp.New(gisp.Box{
		"+": lib.Add,
	})

	env := 0

	gisp.RunJSON(`["+", 1, ["+", 1, 1]]`, &gisp.Context{
		Sandbox: sandbox,
		ENV:     &env,
		PreRun: func(ctx *gisp.Context) {
			*ctx.ENV.(*int) = *ctx.ENV.(*int) + 1
		},
	})

	assert.Equal(t, 7, env)
}

func TestPostRun(t *testing.T) {
	sandbox := gisp.New(gisp.Box{
		"+": lib.Add,
	})

	env := 0

	gisp.RunJSON(`["+", 1, ["+", 1, 1]]`, &gisp.Context{
		Sandbox: sandbox,
		ENV:     &env,
		PostRun: func(ctx *gisp.Context) {
			*ctx.ENV.(*int) = *ctx.ENV.(*int) + 1
		},
	})

	assert.Equal(t, 7, env)
}
func TestAST(t *testing.T) {
	code := []byte(`["*", ["*", 2, 5], ["*", 9, 3]]`)
	ast, _ := djson.Decode(code)

	sandbox := gisp.New(gisp.Box{
		"*": func(ctx *gisp.Context) interface{} {
			a := ctx.ArgNum(1)
			b := ctx.ArgNum(2)
			return a * b
		},
	})

	out := gisp.Run(&gisp.Context{
		AST:     ast,
		Sandbox: sandbox,
	})

	assert.Equal(t, float64(270), out)
}
func TestMissName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic")
		} else {
			assert.Equal(
				t,
				"function \"foo\" is undefined",
				r.(error).Error(),
			)
		}
	}()

	gisp.RunJSON(`["foo"]`, &gisp.Context{
		IsLiftPanic: true,
		Sandbox:     gisp.New(gisp.Box{}),
	})
}

func TestRuntimeErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic")
		} else {
			assert.Equal(
				t,
				"[foo 1 @ 2 @ 0]",
				fmt.Sprint(r.(gisp.Error).Stack),
			)
		}
	}()

	gisp.RunJSON(`["@", ["@", 1, 1], ["@", ["foo"], 1]]`, &gisp.Context{
		IsLiftPanic: true,
		Sandbox: gisp.New(gisp.Box{
			"foo": func(ctx *gisp.Context) interface{} {
				a := []int{}
				a[100] = 1
				return nil
			},
			"@": func(ctx *gisp.Context) interface{} {
				ctx.Arg(1)
				ctx.Arg(2)
				return nil
			},
		}),
	})
}

func TestEmptyFn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic")
		} else {
			assert.Equal(
				t,
				"function [\"foo\"] is undefined",
				fmt.Sprint(r.(gisp.Error).Message),
			)
		}
	}()

	gisp.RunJSON(`[["foo"]]`, &gisp.Context{
		IsLiftPanic: true,
		Sandbox: gisp.New(gisp.Box{
			"foo": func(ctx *gisp.Context) interface{} {
				return nil
			},
		}),
	})
}

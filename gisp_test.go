package gisp_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/gisp"
)

func TestReturnFn(t *testing.T) {
	sandbox := gisp.Sandbox{
		"foo": func(ctx *gisp.Context) interface{} {
			return func(ctx *gisp.Context) float64 {
				return ctx.ArgNum(1) + ctx.ArgNum(2)
			}
		},
	}

	out, _ := gisp.RunJSON([]byte(`[["foo"], 1, 2]`), &gisp.Context{
		Sandbox: sandbox,
	})

	assert.Equal(t, float64(3), out)
}

func TestStr(t *testing.T) {
	sandbox := gisp.Sandbox{}

	out, _ := gisp.RunJSON([]byte(`"foo"`), &gisp.Context{
		Sandbox: sandbox,
	})

	assert.Equal(t, "foo", out)
}

func TestVal(t *testing.T) {
	sandbox := gisp.Sandbox{
		"test": "ok",
	}

	out, _ := gisp.RunJSON([]byte(`["test"]`), &gisp.Context{
		Sandbox: sandbox,
	})

	assert.Equal(t, "ok", out)
}

func TestAST(t *testing.T) {
	code := []byte(`["*", ["*", 2, 5], ["*", 9, 3]]`)
	var ast interface{}
	json.Unmarshal(code, &ast)

	sandbox := gisp.Sandbox{
		"*": func(ctx *gisp.Context) interface{} {
			a := ctx.ArgNum(1)
			b := ctx.ArgNum(2)
			return a * b
		},
	}

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
				"\"foo\" is undefined",
				r.(error).Error(),
			)
		}
	}()

	gisp.RunJSON([]byte(`["foo"]`), &gisp.Context{
		IsLiftPanic: true,
		Sandbox:     gisp.Sandbox{},
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

	gisp.RunJSON([]byte(`["@", ["@", 1, 1], ["@", ["foo"], 1]]`), &gisp.Context{
		IsLiftPanic: true,
		Sandbox: gisp.Sandbox{
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
		},
	})
}

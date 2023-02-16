package test

import (
	"github.com/go-steven/atengo/pkg"
	_ "github.com/go-steven/atengo/pkg/plugins" // 引用所有内部的plugins
	"testing"
)

func TestDefaultEngine_Run(t *testing.T) {
	script := `export func() {
	println("hello world")
	return "xxx"
}
`
	engine := pkg.NewEngine()
	if a, b, err := engine.Eval(script); err != nil {
		panic(err)
	} else {
		println("a:", a, ", b:", b)
	}
}

func TestDefaultEngine_Run2(t *testing.T) {
	script := `
	x := 3
	return x
`
	engine := pkg.NewEngine()
	ret, err := engine.Run(script)
	if err != nil {
		panic(err)
	}
	if ret != "3" {
		t.Errorf("ret:%v, not matched expected 3", ret)
		return
	}
}

func TestDefaultEngine_Eval(t *testing.T) {
	script := `
	is_eval := false
	xxx := func(...args) {
		if is_eval {
			println("hello world")
		} else {
			println(args)
		}
	}

	xxx("aaa")
`
	engine := pkg.NewEngine()
	if _, _, err := engine.Eval(script); err != nil {
		panic(err)
	}
}

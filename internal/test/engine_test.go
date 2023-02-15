package test

import (
	"github.com/go-steven/atengo/internal"
	_ "github.com/go-steven/atengo/internal/plugins" // 引用所有内部的plugins
	"testing"
)

func TestDefaultEngine_Run(t *testing.T) {
	script := `export func() {
	fmt.println("hello world")
}
`
	engine := internal.NewEngine()
	if _, _, err := engine.Eval(script); err != nil {
		panic(err)
	}
}

func TestDefaultEngine_Run2(t *testing.T) {
	script := `
	x := 3
	return x
`
	engine := internal.NewEngine()
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
			fmt.println("hello world")
		} else {
			fmt.println(args)
		}
	}

	xxx("aaa")
`
	engine := internal.NewEngine()
	if _, _, err := engine.Eval(script); err != nil {
		panic(err)
	}
}

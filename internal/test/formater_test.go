package test

import (
	"fmt"
	"github.com/go-steven/atengo/internal"
	"testing"
	"time"
)

func TestScriptFormater_ToSourceModel(t *testing.T) {
	startT := time.Now()
	script := `
export func(a, b) {
	return a + b
}
`
	s := internal.NewScriptFormater()
	ret, err := s.ToSourceModel(script)
	if err != nil {
		panic(err)
	}
	fmt.Println("Params: ", ret.Params)
	fmt.Println("FullCode: ", ret.FullCode)
	println("done, duration:", time.Now().Sub(startT))
}

func TestScriptFormater_ToSourceModel2(t *testing.T) {
	startT := time.Now()
	script := `
	return 3
`
	s := internal.NewScriptFormater()
	ret, err := s.ToSourceModel(script)
	if err != nil {
		panic(err)
	}
	fmt.Println("Params: ", ret.Params)
	fmt.Println("FullCode: ", ret.FullCode)
	println("done, duration:", time.Now().Sub(startT))
}

func TestScriptFormater_ToScript(t *testing.T) {
	startT := time.Now()
	script := `
export func(a, b) {
	return a + b
}
`
	s := internal.NewScriptFormater()
	ret, err := s.ToScript(script)
	if err != nil {
		panic(err)
	}
	fmt.Println("Params: ", ret.Params)
	fmt.Println("FullCode: ", ret.FullCode)
	println("done, duration:", time.Now().Sub(startT))
}

func TestScriptFormater_ToScript2(t *testing.T) {
	startT := time.Now()
	script := `
	return 3
`
	s := internal.NewScriptFormater()
	ret, err := s.ToScript(script)
	if err != nil {
		panic(err)
	}
	fmt.Println("Params: ", ret.Params)
	fmt.Println("FullCode: ", ret.FullCode)
	println("done, duration:", time.Now().Sub(startT))
}

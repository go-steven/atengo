package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/go-steven/atengo/pkg"
	_ "github.com/go-steven/atengo/pkg/plugins" // 引用所有内部的plugins
	"github.com/go-steven/atengo/pkg/util"
	"github.com/mkideal/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// 命令行参数
var (
	versionFlag = flag.Bool("version", false, "show version")
	actFlag     = flag.String("act", "eval", "script act")
)

const (
	ACT_COMPILE = "compile"
	ACT_EVAL    = "eval"
	ACT_RUN     = "run"
)

func main() {
	flag.Parse()
	// 打印版本信息
	Version()
	if *versionFlag {
		return
	}

	log.Info("Enter")
	startT := time.Now()

	modules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	for _, m := range pkg.Plugins() {
		modules.AddBuiltinModule(m.Name(), m.Module())
	}

	inputFile := flag.Arg(0)
	if inputFile == "" || filepath.Ext(inputFile) != ".tengo" {
		fmt.Fprintln(os.Stderr, "need .tengo file")
		os.Exit(1)
	}

	inputData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %s", err.Error())
		os.Exit(1)
	}

	if len(inputData) > 1 && string(inputData[:2]) == "#!" {
		copy(inputData, "//")
	}

	var dataMap map[string]interface{}
	dataFile := flag.Arg(1)
	if dataFile != "" {
		if filepath.Ext(dataFile) != ".json" {
			fmt.Fprintln(os.Stderr, "need .json file")
			os.Exit(1)
		}
		data, err := ioutil.ReadFile(dataFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		//读取的数据为json格式，进行解码
		var _dataMap map[string]interface{}
		if err = json.Unmarshal(data, &_dataMap); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		if len(_dataMap) > 0 {
			dataMap = util.FixJsonMapFloat(_dataMap)
		}
	}

	engine := pkg.NewEngine()
	opts := []pkg.RunOption{}
	if len(dataMap) > 0 {
		opts = append(opts, pkg.FuncData(dataMap))
	}

	var ret, outputLog string
	switch *actFlag {
	case ACT_COMPILE:
		err = engine.Compile(string(inputData), opts...)
	case ACT_RUN:
		ret, err = engine.Run(string(inputData), opts...)
	default:
		ret, outputLog, err = engine.Eval(string(inputData), opts...)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		println("返回日志: ", outputLog)
		os.Exit(1)
	}
	if ret != "" {
		println("返回结果:")
		fmt.Fprintln(os.Stdout, util.Json(ret))
	}
	println("返回日志: ", outputLog)
	fmt.Printf("执行结束, 执行时间：%v\n", time.Now().Sub(startT))
}

var (
	BuildTime   string // 编译时间
	GitRevision string // Git版本
	GitBranch   string // Git分支
	GoVersion   string // Golang信息
)

// Version 版本信息
func Version() {
	fmt.Printf("Build time:\t%s\n", BuildTime)
	fmt.Printf("Git revision:\t%s\n", GitRevision)
	fmt.Printf("Git branch:\t%s\n", GitBranch)
	fmt.Printf("Golang Version:\t%s\n", GoVersion)
}

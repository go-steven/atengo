package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/d5/tengo/v2"
	"github.com/go-steven/atengo/pkg/script_log"
	"github.com/go-steven/atengo/pkg/show_source_line"
	"github.com/mkideal/log"
)

// tengo脚本执行引擎
type Engine interface {
	// 检查tengo脚本是否存在编译错误
	Compile(script string, opts ...RunOption) error

	// 执行tengo脚本，并从给定的返回变量名中获得返回值
	Eval(script string, opts ...RunOption) (string, string, error)

	// 执行tengo脚本，并从给定的返回变量名中获得返回值
	Run(script string, opts ...RunOption) (string, error)
}

func NewEngine() Engine {
	e := &defaultEngine{
		formater: NewScriptFormater(),
	}
	return e
}

type defaultEngine struct {
	formater ScriptFormater
}

func (e *defaultEngine) Compile(script string, opts ...RunOption) error {
	_opts := NewRunOptions(opts...)
	_opts.ShowSourceLine = true
	_opts.Eval = true
	_, _, err := e.CompileAndEval(script, true, opts...)
	return err
}

func (e *defaultEngine) Eval(script string, opts ...RunOption) (string, string, error) {
	_opts := NewRunOptions(opts...)
	_opts.OutputLog = true
	_opts.ShowSourceLine = true
	_opts.Eval = true
	return e.CompileAndEval(script, false, _opts.OptionList()...)
}

func (e *defaultEngine) Run(script string, opts ...RunOption) (string, error) {
	_opts := NewRunOptions(opts...)

	scriptInfo, err := e.formater.ToSourceModel(script, FormaterFuncParams(_opts.FuncParams))
	if err != nil {
		return "", err
	}
	moduleMap := scriptInfo.ModuleMap

	scriptName := V_SCRIPT_NAME
	moduleMap.AddSourceModule(scriptName, []byte(scriptInfo.FullCode))
	//println(buff.String())

	new_buff := bytes.NewBufferString("")
	new_buff.WriteString(fmt.Sprintf("%s := import(\"%s\")\n%s := %s(", scriptName, scriptName, V_RETURN, scriptName))
	for idx := range scriptInfo.Params {
		if idx > 0 {
			new_buff.WriteString(", ")
		}
		new_buff.WriteString(fmt.Sprintf("%s_%d", V_PARAM, idx))
	}
	new_buff.WriteString(")\n")
	//println(new_buff.String())

	inst := tengo.NewScript(new_buff.Bytes())
	inst.SetImports(moduleMap)
	for idx, param := range scriptInfo.Params {
		val, ok := _opts.FuncData[param]
		if !ok {
			return "", fmt.Errorf("函数调用缺少参数值: %s", param)
		}
		inst.Add(fmt.Sprintf("%s_%d", V_PARAM, idx), val)
	}

	compiled, err := inst.Compile()
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	if err := compiled.RunContext(context.Background()); err != nil {
		log.Error(err.Error())
		return "", err
	}
	v := compiled.Get(V_RETURN)
	if err := v.Error(); err != nil {
		log.Error(err.Error())
		return "", err
	}

	if !v.IsUndefined() {
		json_val, err := json.Marshal(v.Value())
		if err != nil {
			log.Error(err.Error())
			return "", err
		}
		return string(json_val), nil
	}
	return "", nil
}

func (e *defaultEngine) CompileAndEval(script string, onlyCompile bool, opts ...RunOption) (string, string, error) {
	_opts := NewRunOptions(opts...)

	scriptLog := script_log.LogFactoryInst()
	logID := scriptLog.NewLogger(_opts.OutputLog)

	scriptInfo, err := e.formater.ToScript(script, FormaterFuncParams(_opts.FuncParams), FormaterEval(_opts.Eval), FormaterLogID(logID))
	if err != nil {
		return "", "", err
	}
	moduleMap := scriptInfo.ModuleMap

	//moduleMap.AddSourceModule(scriptName, []byte(script))
	script = scriptInfo.FullCode
	if _opts.ShowSourceLine {
		println(show_source_line.GetSourceWithLineNum(script))
	}

	inst := tengo.NewScript([]byte(script))
	inst.SetImports(moduleMap)

	for idx, param := range scriptInfo.Params {
		val, ok := _opts.FuncData[param]
		if !ok {
			return "", "", fmt.Errorf("函数调用缺少参数值: %s", param)
		}
		inst.Add(fmt.Sprintf("%s_%d", V_PARAM, idx), val)
	}

	compiled, err := inst.Compile()
	if err != nil {
		log.Error(err.Error())
		if _opts.ShowSourceLine {
			return "", "", fmt.Errorf("%s\n%v", show_source_line.GetSourceWithLineNum(script), err.Error())
		}
		return "", scriptLog.Output(logID), err
	}
	if err := compiled.RunContext(context.Background()); err != nil {
		log.Error(err.Error())
		if _opts.ShowSourceLine {
			return "", "", fmt.Errorf("%s\n%v", show_source_line.GetSourceWithLineNum(script), err.Error())
		}
		return "", scriptLog.Output(logID), err
	}
	v := compiled.Get(V_RETURN)
	if err := v.Error(); err != nil {
		log.Error(err.Error())
		if _opts.ShowSourceLine {
			return "", "", fmt.Errorf("%s\n%v", show_source_line.GetSourceWithLineNum(script), err.Error())
		}
		return "", scriptLog.Output(logID), err
	}

	if !v.IsUndefined() {
		json_val, err := json.Marshal(v.Value())
		if err != nil {
			log.Error(err.Error())
			if _opts.ShowSourceLine {
				return "", "", fmt.Errorf("%s\n%v", show_source_line.GetSourceWithLineNum(script), err.Error())
			}
			return "", scriptLog.Output(logID), err
		}
		return string(json_val), scriptLog.Output(logID), nil
	}
	return "", scriptLog.Output(logID), nil
}

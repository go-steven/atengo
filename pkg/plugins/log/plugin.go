package internal_plugin_log

import (
	"fmt"
	"github.com/d5/tengo/v2"
	"github.com/go-steven/atengo/pkg"
	"github.com/go-steven/atengo/pkg/script_log"
)

func init() {
	pkg.RegisterPlugin(plugin_name, &LogPlugin{})
}

const plugin_name = "LOG"

const packageName = "atengo_LOG"

type LogPlugin struct{}

func (m *LogPlugin) Name() string {
	return plugin_name
}

func (m *LogPlugin) AliasMethod() bool {
	return true
}

func (m *LogPlugin) IsInternal() bool {
	return true
}

func (m *LogPlugin) Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"print":   &tengo.UserFunction{Name: "print", Value: log_print},
		"println": &tengo.UserFunction{Name: "println", Value: log_println},
		"printf":  &tengo.UserFunction{Name: "printf", Value: log_printf},
	}
}

var method_docs = map[string]string{}

func (m *LogPlugin) Doc() map[string]string {
	ret := make(map[string]string)
	for method := range m.Module() {
		ret[method] = method_docs[method]
	}
	return ret
}

func log_print(args ...tengo.Object) (ret tengo.Object, err error) {
	funcName := "log_print"

	if len(args) != 2 {
		return nil, fmt.Errorf("%s: invalid args", funcName)
	}

	logID, ok := args[0].(*tengo.String)
	if !ok {
		return nil, fmt.Errorf("%s: invalid args[0]", funcName)
	}
	scriptLog := script_log.LogFactoryInst().GetLogger(logID.Value)
	if scriptLog == nil {
		return nil, fmt.Errorf("%s: invalid args[0]", funcName)
	}
	args_v, ok := args[1].(*tengo.Array)
	if !ok {
		return nil, fmt.Errorf("%s: invalid args[2]", funcName)
	}

	var log_args []interface{}
	for _, v := range args_v.Value {
		log_args = append(log_args, convert_tengo_value(v))
	}
	scriptLog.Print(log_args...)
	return nil, nil
}

func log_println(args ...tengo.Object) (ret tengo.Object, err error) {
	funcName := "log_println"

	if len(args) != 2 {
		return nil, fmt.Errorf("%s: invalid args", funcName)
	}

	logID, ok := args[0].(*tengo.String)
	if !ok {
		return nil, fmt.Errorf("%s: invalid args[0]", funcName)
	}
	scriptLog := script_log.LogFactoryInst().GetLogger(logID.Value)
	if scriptLog == nil {
		return nil, fmt.Errorf("%s: invalid args[0]", funcName)
	}
	args_v, ok := args[1].(*tengo.Array)
	if !ok {
		return nil, fmt.Errorf("%s: invalid args[2]", funcName)
	}

	var log_args []interface{}
	for _, v := range args_v.Value {
		log_args = append(log_args, convert_tengo_value(v))
	}
	scriptLog.Println(log_args...)

	return nil, nil
}

func log_printf(args ...tengo.Object) (ret tengo.Object, err error) {
	funcName := "log_printf"

	if len(args) < 2 {
		return nil, fmt.Errorf("%s: invalid args", funcName)
	}

	logID, ok := args[0].(*tengo.String)
	if !ok {
		return nil, fmt.Errorf("%s: invalid args[0]", funcName)
	}
	scriptLog := script_log.LogFactoryInst().GetLogger(logID.Value)
	if scriptLog == nil {
		return nil, fmt.Errorf("%s: invalid args[0]", funcName)
	}
	format, ok := args[1].(*tengo.String)
	if !ok {
		return nil, fmt.Errorf("%s: invalid args[1]", funcName)
	}
	args_v, ok := args[2].(*tengo.Array)
	if !ok {
		return nil, fmt.Errorf("%s: invalid args[2]", funcName)
	}

	var log_args []interface{}
	for _, v := range args_v.Value {
		log_args = append(log_args, convert_tengo_value(v))
	}
	scriptLog.Printf(format.Value, log_args...)
	return nil, nil
}

func convert_tengo_value(v tengo.Object) (ret interface{}) {
	switch v.(type) {
	case *tengo.Int:
		val, _ := v.(*tengo.Int)
		ret = val.Value
	case *tengo.String:
		val, _ := v.(*tengo.String)
		ret = val.Value
	case *tengo.Float:
		val, _ := v.(*tengo.Float)
		ret = val.Value
	case *tengo.Bool:
		val, _ := v.(*tengo.Bool)
		ret = !val.IsFalsy()
	case *tengo.Char:
		val, _ := v.(*tengo.Char)
		ret = val.Value
	case *tengo.Time:
		val, _ := v.(*tengo.Time)
		ret = val.Value
	case *tengo.Error:
		val, _ := v.(*tengo.Error)
		ret = val.Value
	case *tengo.Array:
		val, _ := v.(*tengo.Array)
		_ret := []interface{}{}
		for _, _val := range val.Value {
			_ret = append(_ret, convert_tengo_value(_val))
		}
		ret = _ret
	case *tengo.ImmutableArray:
		val, _ := v.(*tengo.ImmutableArray)
		_ret := []interface{}{}
		for _, _val := range val.Value {
			_ret = append(_ret, convert_tengo_value(_val))
		}
		ret = _ret
	case *tengo.Map:
		val, _ := v.(*tengo.Map)
		_ret := make(map[string]interface{})
		for _k, _val := range val.Value {
			_ret[_k] = convert_tengo_value(_val)
		}
		ret = _ret
	case *tengo.ImmutableMap:
		val, _ := v.(*tengo.ImmutableMap)
		_ret := make(map[string]interface{})
		for _k, _val := range val.Value {
			_ret[_k] = convert_tengo_value(_val)
		}
		ret = _ret
	case *tengo.Bytes:
		val, _ := v.(*tengo.Bytes)
		ret = val.Value
	default:
		ret = v.String()
	}
	return
}

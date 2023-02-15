package internal

import (
	"bytes"
	"fmt"
	"github.com/d5/tengo/v2"
	tengo_stdlib "github.com/d5/tengo/v2/stdlib"
	"regexp"
	"strings"
)

type FormaterOptions struct {
	FuncParams    []string
	LogID         string
	Eval          bool
	ImportModules map[string]struct{}
	ModuleMap     *tengo.ModuleMap
}

type FormaterOption func(*FormaterOptions)

// constructor
func NewFormaterOptions(opts ...FormaterOption) FormaterOptions {
	opt := FormaterOptions{}
	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func (o *FormaterOptions) OptionList() []FormaterOption {
	opts := []FormaterOption{
		FormaterEval(o.Eval),
		FormaterLogID(o.LogID),
	}
	if len(o.FuncParams) > 0 {
		opts = append(opts, FormaterFuncParams(o.FuncParams))
	}
	if len(o.ImportModules) > 0 {
		opts = append(opts, FormaterImportModules(o.ImportModules))
	}
	if o.ModuleMap != nil {
		opts = append(opts, FormaterModuleMap(o.ModuleMap))
	}
	return opts
}

func FormaterFuncParams(params []string) FormaterOption {
	return func(o *FormaterOptions) {
		o.FuncParams = params
	}
}

func FormaterEval(enabled bool) FormaterOption {
	return func(o *FormaterOptions) {
		o.Eval = enabled
	}
}

func FormaterLogID(logID string) FormaterOption {
	return func(o *FormaterOptions) {
		o.LogID = logID
	}
}

func FormaterImportModules(importModules map[string]struct{}) FormaterOption {
	return func(o *FormaterOptions) {
		o.ImportModules = importModules
	}
}

func FormaterModuleMap(moduleMap *tengo.ModuleMap) FormaterOption {
	return func(o *FormaterOptions) {
		o.ModuleMap = moduleMap
	}
}

type FullScriptInfo struct {
	ModuleMap     *tengo.ModuleMap
	ImportModules map[string]struct{}
	Params        []string
	FullCode      string
}

type ScriptFormater interface {
	ToSourceModel(script string, opts ...FormaterOption) (*FullScriptInfo, error)
	ToScript(script string, opts ...FormaterOption) (*FullScriptInfo, error)
}

type scriptFormaterImpl struct {
	packageName string

	std_modules map[string]struct{}
	plugins     map[string]IPlugin
}

func NewScriptFormater() ScriptFormater {
	e := &scriptFormaterImpl{
		packageName: "scriptFormaterImpl",
		plugins:     Plugins(),
	}
	e.std_modules = make(map[string]struct{})
	for _, v := range tengo_stdlib.AllModuleNames() {
		e.std_modules[v] = struct{}{}
	}
	return e
}

func (e *scriptFormaterImpl) ToSourceModel(script string, opts ...FormaterOption) (*FullScriptInfo, error) {
	_opts := NewFormaterOptions(opts...)

	rtssFuncMap := e.getFuncMap(script)
	if len(rtssFuncMap) > 0 {
		script = strings.Replace(script, "RTSS_FUNC.", "RTSS_FUNC_", -1)
	}

	params, hasExport := e.getExportParams(script)
	if !hasExport {
		params = []string{}
		if len(_opts.FuncParams) > 0 {
			for _, v := range _opts.FuncParams {
				params = append(params, v)
			}
		}
	}

	importModules, moduleMap, err := e.GetModuleMapByNoImportScript(script, _opts.ImportModules, _opts.ModuleMap)
	if err != nil {
		return nil, err
	}

	//if len(rtssFuncMap) > 0 {
	//	funcGeter := script_lib_proxy.ScriptGeter()
	//	for name, funcModelName := range rtssFuncMap {
	//		if _, ok := importModules[funcModelName]; ok {
	//			return nil, errors.New("函数存在循环引用")
	//		}
	//
	//		f, err := funcGeter.GetFunc(name)
	//		if err != nil {
	//			return nil, err
	//		}
	//		fullFunc, err := e.ToSourceModel(f.Content, FormaterImportModules(importModules), FormaterModuleMap(moduleMap))
	//		if err != nil {
	//			return nil, err
	//		}
	//		importModules = fullFunc.ImportModules
	//		moduleMap = fullFunc.ModuleMap
	//
	//		moduleMap.AddSourceModule(funcModelName, []byte(fullFunc.FullCode))
	//		importModules[funcModelName] = struct{}{}
	//	}
	//	script = strings.Replace(script, "RTSS_FUNC.", "RTSS_FUNC_", -1)
	//}

	buff := e.export_script_with_import_enhance(&ScriptEnhanceInputModel{
		Script:        script,
		HasExport:     hasExport,
		ExportParams:  params,
		ImportModules: importModules,
		LogID:         _opts.LogID,
		IsEval:        _opts.Eval,
	})

	return &FullScriptInfo{
		ModuleMap:     moduleMap,
		ImportModules: importModules,
		FullCode:      buff.String(),
		Params:        params,
	}, nil
}

func (e *scriptFormaterImpl) ToScript(script string, opts ...FormaterOption) (*FullScriptInfo, error) {
	_opts := NewFormaterOptions(opts...)

	funcMap := e.getFuncMap(script)
	if len(funcMap) > 0 {
		script = strings.Replace(script, "FUNC.", "FUNC_", -1)
	}

	params, hasExport := e.getExportParams(script)
	if !hasExport {
		params = []string{}
		if len(_opts.FuncParams) > 0 {
			for _, v := range _opts.FuncParams {
				params = append(params, v)
			}
		}
	}

	importModules, moduleMap, err := e.GetModuleMapByNoImportScript(script, _opts.ImportModules, _opts.ModuleMap)
	if err != nil {
		return nil, err
	}
	//if len(funcMap) > 0 {
	//	funcGeter := script_lib_proxy.ScriptGeter()
	//	for name, funcModelName := range funcMap {
	//		f, err := funcGeter.GetFunc(name)
	//		if err != nil {
	//			return nil, err
	//		}
	//		fullFunc, err := e.ToSourceModel(f.Content, FormaterImportModules(importModules), FormaterModuleMap(moduleMap))
	//		if err != nil {
	//			return nil, err
	//		}
	//		importModules = fullFunc.ImportModules
	//		moduleMap = fullFunc.ModuleMap
	//
	//		moduleMap.AddSourceModule(funcModelName, []byte(fullFunc.FullCode))
	//		importModules[funcModelName] = struct{}{}
	//	}
	//	script = strings.Replace(script, "FUNC.", "FUNC_", -1)
	//}

	buff := e.script_with_import_enhance(&ScriptEnhanceInputModel{
		Script:        script,
		HasExport:     hasExport,
		ExportParams:  params,
		ImportModules: importModules,
		LogID:         _opts.LogID,
		IsEval:        _opts.Eval,
	})

	return &FullScriptInfo{
		ModuleMap:     moduleMap,
		ImportModules: importModules,
		FullCode:      buff.String(),
		Params:        params,
	}, nil
}

func (s *scriptFormaterImpl) getFuncMap(script string) (ret map[string]string) {
	ret = make(map[string]string)
	re := regexp.MustCompile(`FUNC\.(\w+)\(`)
	matchedList := re.FindAllStringSubmatch(script, -1)
	if len(matchedList) == 0 {
		return
	}
	for _, vals := range matchedList {
		if len(vals) != 2 {
			continue
		}
		ret["FUNC."+vals[1]] = "FUNC_" + vals[1]
	}
	return
}

func (s *scriptFormaterImpl) getExportParams(script string) ([]string, bool) {
	re := regexp.MustCompile(`export func\((.*?)\)`)
	matchedList := re.FindAllStringSubmatch(script, 1)
	if len(matchedList) == 0 {
		return nil, false
	}
	var ret []string
	for _, vals := range matchedList {
		if len(vals) != 2 {
			continue
		}
		args := strings.Split(vals[1], ",")
		for _, v := range args {
			v = strings.TrimSpace(v)
			if v != "" {
				ret = append(ret, strings.TrimSpace(v))
			}
		}
	}
	return ret, true
}

func (s *scriptFormaterImpl) GetModuleMapByNoImportScript(script string, importModules map[string]struct{}, moduleMap *tengo.ModuleMap) (map[string]struct{}, *tengo.ModuleMap, error) {
	if moduleMap == nil {
		moduleMap = tengo_stdlib.GetModuleMap(tengo_stdlib.AllModuleNames()...)
	}
	if importModules == nil {
		importModules = make(map[string]struct{})
	}

	re := regexp.MustCompile(`(\w+)\.`)
	matchedList := re.FindAllStringSubmatch(script, -1)
	for _, vals := range matchedList {
		if len(vals) != 2 {
			continue
		}
		moduleName := strings.TrimSpace(vals[1])
		if _, ok := s.std_modules[moduleName]; ok {
			importModules[moduleName] = struct{}{}
			continue
		}
		if m, ok := s.plugins[moduleName]; ok {
			importModules[moduleName] = struct{}{}
			moduleMap.AddBuiltinModule(m.Name(), m.Module())
			continue
		}
	}

	// 使用方法别名的插件，默认自动加载
	for _, v := range Plugins() {
		if v.AliasMethod() {
			importModules[v.Name()] = struct{}{}
			moduleMap.AddBuiltinModule(v.Name(), v.Module())
		}
		if v.Name() == "LOG" || v.Name() == "RTSS" {
			importModules[v.Name()] = struct{}{}
			moduleMap.AddBuiltinModule(v.Name(), v.Module())
		}
	}
	return importModules, moduleMap, nil
}

const V_RETURN = "V_RETURN"
const V_INTERNAL_F = "V_INTERNAL_F"
const V_PARAM = "V_PARAM"
const V_SCRIPT_NAME = "V_SCRIPT_NAME"

type ScriptEnhanceInputModel struct {
	Script        string
	HasExport     bool
	ExportParams  []string
	ImportModules map[string]struct{}
	LogID         string
	IsEval        bool
}

func (e *scriptFormaterImpl) script_with_import_enhance(input *ScriptEnhanceInputModel) *bytes.Buffer {
	buff := bytes.NewBufferString("\n")
	buff.WriteString(fmt.Sprintf("is_eval := %v\n", input.IsEval))
	for moduleName := range input.ImportModules {
		buff.WriteString(fmt.Sprintf(`%s := import("%s")`, moduleName, moduleName) + "\n")
	}
	buff.WriteString("\n")

	alias_methods := make(map[string]struct{})
	if input.LogID == "" {
		buff.WriteString(`
// alias
print := func(...args) {}
println := func(...args) {}
printf := func(format, ...args) {}
`)
	} else {
		buff.WriteString(`
// alias
print := func(...args) {
	if is_eval {
		LOG.print("` + input.LogID + `", args)
	}
}
println := func(...args) {
	if is_eval {
		LOG.println("` + input.LogID + `", args)
	}
}
printf := func(format, ...args) {
	if is_eval {
		LOG.printf("` + input.LogID + `", format, args)
	}
}
`)
	}
	buff.WriteString("\n")

	alias_methods["print"] = struct{}{}
	alias_methods["println"] = struct{}{}
	alias_methods["printf"] = struct{}{}

	// 使用方法别名的插件，默认自动加载
	for _, v := range Plugins() {
		if v.AliasMethod() {
			buff.WriteString("\n")
			for method := range v.Module() {
				if _, ok := alias_methods[method]; ok {
					continue
				}
				buff.WriteString(fmt.Sprintf("%s := %s.%s\n", method, v.Name(), method))
			}
			buff.WriteString("\n")
		}
	}

	script_func_name := V_INTERNAL_F
	if input.HasExport {
		buff.WriteString(strings.Replace(input.Script, `export func(`, fmt.Sprintf("%s := func(", script_func_name), 1))
	} else {
		buff.WriteString(fmt.Sprintf("%s := func(", script_func_name))

		for idx, param := range input.ExportParams {
			if idx > 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(param)
		}
		buff.WriteString(") {")
		buff.WriteString(input.Script)
		buff.WriteString("\n}")
		buff.WriteString("")
	}
	buff.WriteString("\n")

	buff.WriteString(fmt.Sprintf("V_RETURN := %s(", script_func_name))
	for idx := range input.ExportParams {
		if idx > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(fmt.Sprintf("%s_%d", V_PARAM, idx))
	}
	buff.WriteString(")")

	return buff
}

func (e *scriptFormaterImpl) export_script_with_import_enhance(input *ScriptEnhanceInputModel) *bytes.Buffer {
	buff := bytes.NewBufferString("\n")
	buff.WriteString(fmt.Sprintf("is_eval := %v\n", input.IsEval))
	for moduleName := range input.ImportModules {
		buff.WriteString(fmt.Sprintf(`%s := import("%s")`, moduleName, moduleName) + "\n")
	}
	buff.WriteString("\n")

	alias_methods := make(map[string]struct{})
	if input.LogID == "" {
		buff.WriteString(`
// alias
print := func(...args) {}
println := func(...args) {}
printf := func(format, ...args) {}
`)
	} else {
		buff.WriteString(`
// alias
print := func(...args) {
	if is_eval {
		LOG.print("` + input.LogID + `", args)
	}
}
println := func(...args) {
	if is_eval {
		LOG.println("` + input.LogID + `", args)
	}
}
printf := func(format, ...args) {
	if is_eval {
		LOG.printf("` + input.LogID + `", format, args)
	}
}
`)
	}
	buff.WriteString("\n")

	alias_methods["print"] = struct{}{}
	alias_methods["println"] = struct{}{}
	alias_methods["printf"] = struct{}{}

	// 使用方法别名的插件，默认自动加载
	for _, v := range Plugins() {
		if v.AliasMethod() {
			buff.WriteString("\n")
			for method := range v.Module() {
				if _, ok := alias_methods[method]; ok {
					continue
				}
				buff.WriteString(fmt.Sprintf("%s := %s.%s\n", method, v.Name(), method))
			}
			buff.WriteString("\n")
		}
	}

	if !input.HasExport {
		buff.WriteString(`export func(`)

		for idx, param := range input.ExportParams {
			if idx > 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(param)
		}
		buff.WriteString(") {")
		buff.WriteString(input.Script)
		buff.WriteString("\n}")
	} else {
		buff.WriteString(input.Script)
	}
	buff.WriteString("\n")

	return buff
}

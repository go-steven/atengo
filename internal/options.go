package internal

type RunOptions struct {
	FuncParams     []string
	FuncData       map[string]interface{}
	OutputLog      bool
	ShowSourceLine bool
	Eval           bool
}

type RunOption func(*RunOptions)

// constructor
func NewRunOptions(opts ...RunOption) RunOptions {
	opt := RunOptions{}
	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func (o *RunOptions) OptionList() []RunOption {
	opts := []RunOption{
		OutputLog(o.OutputLog),
		ShowSourceLine(o.ShowSourceLine),
		Eval(o.Eval),
	}
	if len(o.FuncParams) > 0 {
		opts = append(opts, FuncParams(o.FuncParams))
	}
	if o.FuncData != nil {
		opts = append(opts, FuncData(o.FuncData))
	}
	return opts
}

func FuncParams(params []string) RunOption {
	return func(o *RunOptions) {
		o.FuncParams = params
	}
}

func FuncData(data map[string]interface{}) RunOption {
	return func(o *RunOptions) {
		o.FuncData = data
	}
}

func OutputLog(enabled bool) RunOption {
	return func(o *RunOptions) {
		o.OutputLog = enabled
	}
}

func ShowSourceLine(enabled bool) RunOption {
	return func(o *RunOptions) {
		o.ShowSourceLine = enabled
	}
}

func Eval(enabled bool) RunOption {
	return func(o *RunOptions) {
		o.Eval = enabled
	}
}

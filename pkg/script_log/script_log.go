package script_log

import (
	"bytes"
	"fmt"
)

type ScriptLog interface {
	Print(a ...interface{})
	Printf(format string, a ...interface{})
	Println(a ...interface{})
	Output() string
}

type scriptLogImpl struct {
	buff    *bytes.Buffer
	enabled bool
}

func NewScriptLog(enabled bool) ScriptLog {
	return &scriptLogImpl{
		buff:    bytes.NewBufferString(""),
		enabled: enabled,
	}
}

func (s *scriptLogImpl) Print(a ...interface{}) {
	if s.enabled {
		for _, v := range a {
			s.buff.WriteString(fmt.Sprintf("%v", v))
		}
	}
}

func (s *scriptLogImpl) Println(a ...interface{}) {
	if s.enabled {
		for _, v := range a {
			s.buff.WriteString(fmt.Sprintf("%v", v))
		}
		s.buff.WriteString("\n")
	}
}

func (s *scriptLogImpl) Printf(format string, a ...interface{}) {
	if s.enabled {
		s.buff.WriteString(fmt.Sprintf(format, a...))
	}
}

func (s *scriptLogImpl) Output() (ret string) {
	if s.enabled {
		ret = s.buff.String()
		s.buff.Reset()
	}
	return
}

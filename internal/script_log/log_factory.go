package script_log

import (
	"github.com/go-steven/atengo/internal/util"
	"sync"
)

type LoggerFactory interface {
	NewLogger(enabled bool) string
	GetLogger(logID string) ScriptLog
	Print(logID string, a ...interface{})
	Printf(logID string, format string, a ...interface{})
	Println(logID string, a ...interface{})
	Output(logID string) string
}

type logFactoryImpl struct {
	data map[string]ScriptLog
	m    *sync.RWMutex
}

// 单例模式
var _inst_LoggerFactory LoggerFactory
var _once_LoggerFactory sync.Once

func LogFactoryInst() LoggerFactory {
	_once_LoggerFactory.Do(func() {
		_inst_LoggerFactory = &logFactoryImpl{
			data: make(map[string]ScriptLog),
			m:    new(sync.RWMutex),
		}
	})
	return _inst_LoggerFactory
}

func (f *logFactoryImpl) NewLogger(enabled bool) string {
	logID := util.Token()
	f.m.Lock()
	f.data[logID] = NewScriptLog(enabled)
	f.m.Unlock()
	return logID
}

func (f *logFactoryImpl) GetLogger(logID string) ScriptLog {
	f.m.RLock()
	v, ok := f.data[logID]
	f.m.RUnlock()
	if !ok {
		return nil
	}
	return v
}

func (f *logFactoryImpl) Print(logID string, a ...interface{}) {
	logger := f.GetLogger(logID)
	if logger != nil {
		logger.Print(a...)
	}
}

func (f *logFactoryImpl) Println(logID string, a ...interface{}) {
	logger := f.GetLogger(logID)
	if logger != nil {
		logger.Println(a...)
	}
}

func (f *logFactoryImpl) Printf(logID string, format string, a ...interface{}) {
	logger := f.GetLogger(logID)
	if logger != nil {
		logger.Printf(format, a...)
	}
}

func (f *logFactoryImpl) Output(logID string) (ret string) {
	logger := f.GetLogger(logID)
	if logger != nil {
		ret = logger.Output()
		f.m.Lock()
		delete(f.data, logID)
		f.m.Unlock()
	}
	return
}

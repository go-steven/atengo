package internal

import (
	"github.com/d5/tengo/v2"
	"sync"
)

type IPlugin interface {
	Doc() map[string]string
	Name() string
	Module() map[string]tengo.Object
	AliasMethod() bool
	IsInternal() bool
}

type ISubPlugin interface {
	Doc() map[string]string
	Module() map[string]tengo.Object
}

type IPluginFactory interface {
	Register(name string, m IPlugin)
	UnRegister(name string)
	Get(name string) IPlugin
	Plugins() map[string]IPlugin
}

func RegisterPlugin(name string, m IPlugin) {
	iPluginFactory.Register(name, m)
}

func UnRegisterPlugin(name string) {
	iPluginFactory.UnRegister(name)
}

func GetPlugin(name string) IPlugin {
	return iPluginFactory.Get(name)
}

func Plugins() map[string]IPlugin {
	return iPluginFactory.Plugins()
}

func init() {
	once.Do(func() {
		iPluginFactory = newPluginFactory()
	})
}

// 单例
var iPluginFactory IPluginFactory
var once sync.Once

type defaultPluginFactory struct {
	data map[string]IPlugin
	m    *sync.RWMutex
}

func newPluginFactory() IPluginFactory {
	s := &defaultPluginFactory{
		data: make(map[string]IPlugin),
		m:    new(sync.RWMutex),
	}
	return s
}

func (s *defaultPluginFactory) Register(name string, m IPlugin) {
	s.m.Lock()
	defer s.m.Unlock()

	s.data[name] = m
}

func (s *defaultPluginFactory) UnRegister(name string) {
	s.m.Lock()
	defer s.m.Unlock()

	delete(s.data, name)
}

func (s *defaultPluginFactory) Get(name string) IPlugin {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.data[name]
}

func (s *defaultPluginFactory) Plugins() map[string]IPlugin {
	s.m.RLock()
	defer s.m.RUnlock()

	ret := make(map[string]IPlugin)
	for k, v := range s.data {
		ret[k] = v
	}
	return ret
}

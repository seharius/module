package module

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type OnErrorHandler func(error)
type OnStartHandler func()
type OnEndHandler func()
type Module struct {
	key, name      string
	parrent        *Module
	onErrorHandler OnErrorHandler
	onStartHandler OnStartHandler
	onEndHandler   OnEndHandler
	status         Status
	onShutingDown  func()
}

type ModuleRegisteredHandler func(*Module)

var (
	onModuleRegisteredHandlers []ModuleRegisteredHandler
	modules                    map[string]*Module = make(map[string]*Module)
	onceStart                  sync.Once
	onceEnd                    sync.Once
)

func RegisterOnModuleRegisteredHandler(handler ModuleRegisteredHandler) {
	onModuleRegisteredHandlers = append(onModuleRegisteredHandlers, handler)
}

func dispatchModuleRegisteredEvent(module *Module) {
	for _, handler := range onModuleRegisteredHandlers {
		handler(module)
	}
}

func Start() {
	defer End()
	go listenSigTerm()
	onceStart.Do(func() {
		for _, module := range modules {
			module.start()
		}
	})
}

func listenSigTerm() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	<-sigs
	ShutDown()
}

func ShutDown() {
	for _, module := range modules {
		module.onShutingDown()
	}
}

func End() {
	onceEnd.Do(func() {
		for _, module := range modules {
			module.End()
		}
	})
}

func Register(key, name string, handler OnStartHandler) *Module {
	var module = &Module{
		key:            key,
		name:           name,
		onStartHandler: handler,
		status:         starting,
	}
	modules[key] = module
	dispatchModuleRegisteredEvent(module)

	return module
}

func Find(key string) *Module {
	if module, ok := modules[key]; ok {
		return module
	}

	return nil
}

func (m *Module) RegisterSubmodule(key, name string, handler OnStartHandler) *Module {
	var submodule = Register(m.key+"."+key, name, handler)
	submodule.parrent = m

	return submodule
}

func (m *Module) FindSubmodule(key string) *Module {
	return Find(m.key + "." + key)
}

func (m *Module) Key() string {
	return m.key
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Parrent() *Module {
	return m.parrent
}

func (m *Module) RegisterOnErrorHandler(handler OnErrorHandler) {
	m.onErrorHandler = handler
}

func (m *Module) RegisterOnEndHandler(handler OnEndHandler) {
	m.onEndHandler = handler
}

func (m *Module) RaiseError(err error) {
	m.onErrorHandler(err)
}

func (m *Module) start() {
	m.onStartHandler()
	m.status = running
}

func (m *Module) End() {
	if m.onEndHandler != nil {
		m.onEndHandler()
	}
}

func (m *Module) Status() Status {
	return m.status
}

func (m *Module) OnShutingDown() {
	m.onShutingDown()
}

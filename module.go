package module

type OnErrorHandler func(error)

type Module struct {
	key, name       string
	parrent         *Module
	onErrorHandlers []OnErrorHandler
}

type ModuleRegisteredHandler func(*Module)

var (
	onModuleRegisteredHandlers []ModuleRegisteredHandler
	modules                    map[string]*Module
)

func RegisterOnModuleRegisteredHandler(handler ModuleRegisteredHandler) {
	onModuleRegisteredHandlers = append(onModuleRegisteredHandlers, handler)
}

func dispatchModuleRegisteredEvent(module *Module) {
	for _, handler := range onModuleRegisteredHandlers {
		handler(module)
	}
}

func Register(key, name string) *Module {
	var module = &Module{
		key:  key,
		name: name,
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

func (m *Module) RegisterSubmodule(key, name string) *Module {
	var submodule = Register(m.key+"."+key, name)
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
	m.onErrorHandlers = append(m.onErrorHandlers, handler)
}

func (m *Module) RaiseError(err error) {
	for _, handler := range m.onErrorHandlers {
		handler(err)
	}
}

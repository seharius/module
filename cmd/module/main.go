package main

import (
	"fmt"

	"github.com/Kankeran/module"
)

var m *module.Module

func main() {
	module.RegisterOnModuleRegisteredHandler(onModuleRegistered)
	onInit()

	m1 := module.Register("api", "api", func() {})
	m2 := m1.RegisterSubmodule("v1", "api v1", func() {})
	m2.RegisterSubmodule("example", "exampler api", func() { fmt.Println("Elo") })
	m3 := m1.RegisterSubmodule("v2", "api v2", func() {})
	m3.RegisterSubmodule("example", "exampler api", func() { fmt.Println("Elo2") })
	fmt.Printf("%v\n", m)
	module.Start()
	module.Start()
	module.Start()
	module.End()
	fmt.Println(m.Status())

}

func onModuleRegistered(module *module.Module) {
	fmt.Println("Module registered", module.Key())
}

func onInit() {
	m = module.Register("test", "test", onStart)
	m.RegisterOnErrorHandler(onError)
	m.RegisterOnEndHandler(onEnd)
}

func onStart() {
	fmt.Println("Module started")
	m.RaiseError(fmt.Errorf("error"))
}

func onError(err error) {
	fmt.Println("Error", err)
}

func onEnd() {
	fmt.Println("Module ended")
}

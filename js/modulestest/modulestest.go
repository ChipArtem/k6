package modulestest

import (
	"context"

	"github.com/ChipArtem/k6/js/common"
	"github.com/ChipArtem/k6/js/modules"
	"github.com/ChipArtem/k6/lib"
	"github.com/dop251/goja"
)

var _ modules.VU = &VU{}

// VU is a modules.VU implementation meant to be used within tests
type VU struct {
	CtxField              context.Context
	InitEnvField          *common.InitEnvironment
	EventsField           common.Events
	StateField            *lib.State
	RuntimeField          *goja.Runtime
	RegisterCallbackField func() func(f func() error)
}

// Context returns internally set field to conform to modules.VU interface
func (m *VU) Context() context.Context {
	return m.CtxField
}

// Events returns internally set field to conform to modules.VU interface
func (m *VU) Events() common.Events {
	return m.EventsField
}

// InitEnv returns internally set field to conform to modules.VU interface
func (m *VU) InitEnv() *common.InitEnvironment {
	m.checkIntegrity()
	return m.InitEnvField
}

// State returns internally set field to conform to modules.VU interface
func (m *VU) State() *lib.State {
	m.checkIntegrity()
	return m.StateField
}

// Runtime returns internally set field to conform to modules.VU interface
func (m *VU) Runtime() *goja.Runtime {
	return m.RuntimeField
}

// RegisterCallback is not really implemented
func (m *VU) RegisterCallback() func(f func() error) {
	return m.RegisterCallbackField()
}

func (m *VU) checkIntegrity() {
	if m.InitEnvField != nil && m.StateField != nil {
		panic("there is a bug in the test: InitEnvField and StateField are not allowed at the same time")
	}
}

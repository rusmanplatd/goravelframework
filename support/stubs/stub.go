package stubs

import "strings"

func ConsoleKernel() string {
	return `package console

import (
	"github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/schedule"
)

type Kernel struct {
}

func (kernel Kernel) Commands() []console.Command {
	return []console.Command{}
}

func (kernel Kernel) Schedule() []schedule.Event {
	return []schedule.Event{}
}
`
}

func DatabaseConfig(module string) string {
	content := `package config

import (
	"DummyModule/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("database", map[string]any{})
}
`

	return strings.ReplaceAll(content, "DummyModule", module)
}

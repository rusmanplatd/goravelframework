package main

type Stubs struct{}

func (s Stubs) ViewFacade() string {
	return `package facades

import (
	"github.com/rusmanplatd/goravelframework/contracts/view"
)

func View() view.View {
	return App().MakeView()
}
`
}

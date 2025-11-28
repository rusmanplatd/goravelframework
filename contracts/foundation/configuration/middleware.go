package configuration

import "github.com/rusmanplatd/goravelframework/contracts/http"

type Middleware interface {
	Append(middleware ...http.Middleware) Middleware
	GetGlobalMiddleware() []http.Middleware
	GetRecover() func(ctx http.Context, err any)
	Prepend(middleware ...http.Middleware) Middleware
	Recover(fn func(ctx http.Context, err any)) Middleware
	Use(middleware ...http.Middleware) Middleware
}

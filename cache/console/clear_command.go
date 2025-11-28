package console

import (
	"github.com/rusmanplatd/goravelframework/contracts/cache"
	"github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/console/command"
)

type ClearCommand struct {
	cache cache.Cache
}

func NewClearCommand(cache cache.Cache) *ClearCommand {
	return &ClearCommand{cache: cache}
}

// Signature The name and signature of the console command.
func (r *ClearCommand) Signature() string {
	return "cache:clear"
}

// Description The console command description.
func (r *ClearCommand) Description() string {
	return "Flush the application cache"
}

// Extend The console command extend.
func (r *ClearCommand) Extend() command.Extend {
	return command.Extend{
		Category: "cache",
	}
}

// Handle Execute the console command.
func (r *ClearCommand) Handle(ctx console.Context) error {
	if r.cache.Flush() {
		ctx.Success("Application cache cleared")
	} else {
		ctx.Error("Clear Application cache Failed")
	}

	return nil
}

package facades

import (
	"context"

	"github.com/rusmanplatd/goravelframework/contracts/auth"
	"github.com/rusmanplatd/goravelframework/contracts/auth/access"
	"github.com/rusmanplatd/goravelframework/contracts/broadcast"
	"github.com/rusmanplatd/goravelframework/contracts/cache"
	"github.com/rusmanplatd/goravelframework/contracts/config"
	"github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/crypt"
	"github.com/rusmanplatd/goravelframework/contracts/database/db"
	"github.com/rusmanplatd/goravelframework/contracts/database/orm"
	"github.com/rusmanplatd/goravelframework/contracts/database/schema"
	"github.com/rusmanplatd/goravelframework/contracts/database/seeder"
	"github.com/rusmanplatd/goravelframework/contracts/event"
	"github.com/rusmanplatd/goravelframework/contracts/filesystem"
	foundationcontract "github.com/rusmanplatd/goravelframework/contracts/foundation"
	"github.com/rusmanplatd/goravelframework/contracts/grpc"
	"github.com/rusmanplatd/goravelframework/contracts/hash"
	"github.com/rusmanplatd/goravelframework/contracts/http"
	"github.com/rusmanplatd/goravelframework/contracts/http/client"
	"github.com/rusmanplatd/goravelframework/contracts/log"
	"github.com/rusmanplatd/goravelframework/contracts/mail"
	"github.com/rusmanplatd/goravelframework/contracts/notification"
	"github.com/rusmanplatd/goravelframework/contracts/process"
	"github.com/rusmanplatd/goravelframework/contracts/queue"
	"github.com/rusmanplatd/goravelframework/contracts/route"
	"github.com/rusmanplatd/goravelframework/contracts/schedule"
	"github.com/rusmanplatd/goravelframework/contracts/session"
	"github.com/rusmanplatd/goravelframework/contracts/testing"
	"github.com/rusmanplatd/goravelframework/contracts/translation"
	"github.com/rusmanplatd/goravelframework/contracts/validation"
	"github.com/rusmanplatd/goravelframework/contracts/view"
	"github.com/rusmanplatd/goravelframework/errors"
	"github.com/rusmanplatd/goravelframework/foundation"
)

func App() foundationcontract.Application {
	if foundation.App == nil {
		panic(errors.ApplicationNotSet.SetModule(errors.ModuleFacade))
	} else {
		return foundation.App
	}
}

func Artisan() console.Artisan {
	return App().MakeArtisan()
}

func Auth(ctx ...http.Context) auth.Auth {
	return App().MakeAuth(ctx...)
}

func Broadcast() broadcast.Manager {
	return App().MakeBroadcast()
}

func Cache() cache.Cache {
	return App().MakeCache()
}

func Config() config.Config {
	return App().MakeConfig()
}

func Crypt() crypt.Crypt {
	return App().MakeCrypt()
}

func DB() db.DB {
	return App().MakeDB()
}

func Event() event.Instance {
	return App().MakeEvent()
}

func Gate() access.Gate {
	return App().MakeGate()
}

func Grpc() grpc.Grpc {
	return App().MakeGrpc()
}

func Hash() hash.Hash {
	return App().MakeHash()
}

func Http() client.Request {
	return App().MakeHttp()
}

func Lang(ctx context.Context) translation.Translator {
	return App().MakeLang(ctx)
}

func Log() log.Log {
	return App().MakeLog()
}

func Mail() mail.Mail {
	return App().MakeMail()
}

func Notification() notification.Factory {
	return App().MakeNotification()
}

func Orm() orm.Orm {
	return App().MakeOrm()
}

func Process() process.Process {
	return App().MakeProcess()
}

func Queue() queue.Queue {
	return App().MakeQueue()
}

func RateLimiter() http.RateLimiter {
	return App().MakeRateLimiter()
}

func Route() route.Route {
	return App().MakeRoute()
}

func Schedule() schedule.Schedule {
	return App().MakeSchedule()
}

func Schema() schema.Schema {
	return App().MakeSchema()
}

func Seeder() seeder.Facade {
	return App().MakeSeeder()
}

func Session() session.Manager {
	return App().MakeSession()
}

func Storage() filesystem.Storage {
	return App().MakeStorage()
}

func Testing() testing.Testing {
	return App().MakeTesting()
}

func Validation() validation.Validation {
	return App().MakeValidation()
}

func View() view.View {
	return App().MakeView()
}

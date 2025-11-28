package foundation

import (
	"context"
	"fmt"
	"sync"

	contractsauth "github.com/rusmanplatd/goravelframework/contracts/auth"
	contractsaccess "github.com/rusmanplatd/goravelframework/contracts/auth/access"
	"github.com/rusmanplatd/goravelframework/contracts/binding"
	contractsbroadcast "github.com/rusmanplatd/goravelframework/contracts/broadcast"
	contractscache "github.com/rusmanplatd/goravelframework/contracts/cache"
	contractsconfig "github.com/rusmanplatd/goravelframework/contracts/config"
	contractsconsole "github.com/rusmanplatd/goravelframework/contracts/console"
	contractscrypt "github.com/rusmanplatd/goravelframework/contracts/crypt"
	contractsdb "github.com/rusmanplatd/goravelframework/contracts/database/db"
	contractsorm "github.com/rusmanplatd/goravelframework/contracts/database/orm"
	contractsmigration "github.com/rusmanplatd/goravelframework/contracts/database/schema"
	contractsseerder "github.com/rusmanplatd/goravelframework/contracts/database/seeder"
	contractsevent "github.com/rusmanplatd/goravelframework/contracts/event"
	"github.com/rusmanplatd/goravelframework/contracts/facades"
	contractsfilesystem "github.com/rusmanplatd/goravelframework/contracts/filesystem"
	contractsfoundation "github.com/rusmanplatd/goravelframework/contracts/foundation"
	contractsgrpc "github.com/rusmanplatd/goravelframework/contracts/grpc"
	contractshash "github.com/rusmanplatd/goravelframework/contracts/hash"
	contractshttp "github.com/rusmanplatd/goravelframework/contracts/http"
	contractshttpclient "github.com/rusmanplatd/goravelframework/contracts/http/client"
	contractslog "github.com/rusmanplatd/goravelframework/contracts/log"
	contractsmail "github.com/rusmanplatd/goravelframework/contracts/mail"
	contractsnotification "github.com/rusmanplatd/goravelframework/contracts/notification"
	contractspipeline "github.com/rusmanplatd/goravelframework/contracts/pipeline"
	contractsprocess "github.com/rusmanplatd/goravelframework/contracts/process"
	contractsqueue "github.com/rusmanplatd/goravelframework/contracts/queue"
	contractsroute "github.com/rusmanplatd/goravelframework/contracts/route"
	contractsschedule "github.com/rusmanplatd/goravelframework/contracts/schedule"
	contractsession "github.com/rusmanplatd/goravelframework/contracts/session"
	contractstesting "github.com/rusmanplatd/goravelframework/contracts/testing"
	contractstranslation "github.com/rusmanplatd/goravelframework/contracts/translation"
	contractsvalidation "github.com/rusmanplatd/goravelframework/contracts/validation"
	contractsview "github.com/rusmanplatd/goravelframework/contracts/view"
	"github.com/rusmanplatd/goravelframework/support/color"
)

type instance struct {
	concrete any
	shared   bool
}

type Container struct {
	bindings  sync.Map
	instances sync.Map
}

func NewContainer() *Container {
	return &Container{}
}

func (r *Container) Bind(key any, callback func(app contractsfoundation.Application) (any, error)) {
	r.bindings.Store(key, instance{concrete: callback, shared: false})
}

func (r *Container) Bindings() []any {
	var bindings []any
	r.bindings.Range(func(key, value any) bool {
		bindings = append(bindings, key)
		return true
	})
	return bindings
}

func (r *Container) BindWith(key any, callback func(app contractsfoundation.Application, parameters map[string]any) (any, error)) {
	r.bindings.Store(key, instance{concrete: callback, shared: false})
}

func (r *Container) Fresh(bindings ...any) {
	if len(bindings) == 0 {
		r.instances.Range(func(key, value any) bool {
			if key != binding.Config {
				r.instances.Delete(key)
			}

			return true
		})
	} else {
		for _, binding := range bindings {
			r.instances.Delete(binding)
		}
	}
}

func (r *Container) Instance(key any, ins any) {
	r.bindings.Store(key, instance{concrete: ins, shared: true})
}

func (r *Container) Make(key any) (any, error) {
	return r.make(key, nil)
}

func (r *Container) MakeArtisan() contractsconsole.Artisan {
	instance, err := r.Make(facades.FacadeToBinding[facades.Artisan])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsconsole.Artisan)
}

func (r *Container) MakeAuth(ctx ...contractshttp.Context) contractsauth.Auth {
	parameters := map[string]any{}
	if len(ctx) > 0 {
		parameters["ctx"] = ctx[0]
	}

	instance, err := r.MakeWith(facades.FacadeToBinding[facades.Auth], parameters)
	if err != nil {
		color.Errorln(err)
		return nil
	}
	if instance == nil {
		return nil
	}

	return instance.(contractsauth.Auth)
}

func (r *Container) MakeBroadcast() contractsbroadcast.Manager {
	instance, err := r.Make(facades.FacadeToBinding[facades.Broadcast])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsbroadcast.Manager)
}

func (r *Container) MakeCache() contractscache.Cache {
	instance, err := r.Make(facades.FacadeToBinding[facades.Cache])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractscache.Cache)
}

func (r *Container) MakeConfig() contractsconfig.Config {
	instance, err := r.Make(facades.FacadeToBinding[facades.Config])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsconfig.Config)
}

func (r *Container) MakeCrypt() contractscrypt.Crypt {
	instance, err := r.Make(facades.FacadeToBinding[facades.Crypt])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractscrypt.Crypt)
}

func (r *Container) MakeDB() contractsdb.DB {
	instance, err := r.Make(facades.FacadeToBinding[facades.DB])
	if err != nil {
		color.Errorln(err)
		return nil
	}
	if instance == nil {
		return nil
	}

	return instance.(contractsdb.DB)
}

func (r *Container) MakeEvent() contractsevent.Instance {
	instance, err := r.Make(facades.FacadeToBinding[facades.Event])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsevent.Instance)
}

func (r *Container) MakeGate() contractsaccess.Gate {
	instance, err := r.Make(facades.FacadeToBinding[facades.Gate])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsaccess.Gate)
}

func (r *Container) MakeGrpc() contractsgrpc.Grpc {
	instance, err := r.Make(facades.FacadeToBinding[facades.Grpc])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsgrpc.Grpc)
}

func (r *Container) MakeHash() contractshash.Hash {
	instance, err := r.Make(facades.FacadeToBinding[facades.Hash])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractshash.Hash)
}

func (r *Container) MakeHttp() contractshttpclient.Request {
	instance, err := r.Make(facades.FacadeToBinding[facades.Http])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractshttpclient.Request)
}

func (r *Container) MakeLang(ctx context.Context) contractstranslation.Translator {
	instance, err := r.MakeWith(facades.FacadeToBinding[facades.Lang], map[string]any{
		"ctx": ctx,
	})
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractstranslation.Translator)
}

func (r *Container) MakeLog() contractslog.Log {
	instance, err := r.Make(facades.FacadeToBinding[facades.Log])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractslog.Log)
}

func (r *Container) MakeMail() contractsmail.Mail {
	instance, err := r.Make(facades.FacadeToBinding[facades.Mail])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsmail.Mail)
}

func (r *Container) MakeNotification() contractsnotification.Factory {
	instance, err := r.Make(facades.FacadeToBinding[facades.Notification])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsnotification.Factory)
}

func (r *Container) MakeOrm() contractsorm.Orm {
	instance, err := r.Make(facades.FacadeToBinding[facades.Orm])
	if err != nil {
		color.Errorln(err)
		return nil
	}
	if instance == nil {
		return nil
	}

	return instance.(contractsorm.Orm)
}

func (r *Container) MakePipeline() contractspipeline.Pipeline {
	instance, err := r.Make(binding.Pipeline)
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractspipeline.Pipeline)
}

func (r *Container) MakeProcess() contractsprocess.Process {
	instance, err := r.Make(facades.FacadeToBinding[facades.Process])
	if err != nil {
		color.Errorln(err)
		return nil
	}
	if instance == nil {
		return nil
	}

	return instance.(contractsprocess.Process)
}

func (r *Container) MakeQueue() contractsqueue.Queue {
	instance, err := r.Make(facades.FacadeToBinding[facades.Queue])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsqueue.Queue)
}

func (r *Container) MakeRateLimiter() contractshttp.RateLimiter {
	instance, err := r.Make(facades.FacadeToBinding[facades.RateLimiter])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractshttp.RateLimiter)
}

func (r *Container) MakeRoute() contractsroute.Route {
	instance, err := r.Make(facades.FacadeToBinding[facades.Route])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsroute.Route)
}

func (r *Container) MakeSchedule() contractsschedule.Schedule {
	instance, err := r.Make(facades.FacadeToBinding[facades.Schedule])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsschedule.Schedule)
}

func (r *Container) MakeSchema() contractsmigration.Schema {
	instance, err := r.Make(facades.FacadeToBinding[facades.Schema])
	if err != nil {
		color.Errorln(err)
		return nil
	}
	if instance == nil {
		return nil
	}

	return instance.(contractsmigration.Schema)
}

func (r *Container) MakeSeeder() contractsseerder.Facade {
	instance, err := r.Make(facades.FacadeToBinding[facades.Seeder])

	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsseerder.Facade)
}

func (r *Container) MakeSession() contractsession.Manager {
	instance, err := r.Make(facades.FacadeToBinding[facades.Session])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsession.Manager)
}

func (r *Container) MakeStorage() contractsfilesystem.Storage {
	instance, err := r.Make(facades.FacadeToBinding[facades.Storage])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsfilesystem.Storage)
}

func (r *Container) MakeTesting() contractstesting.Testing {
	instance, err := r.Make(facades.FacadeToBinding[facades.Testing])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractstesting.Testing)
}

func (r *Container) MakeValidation() contractsvalidation.Validation {
	instance, err := r.Make(facades.FacadeToBinding[facades.Validation])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsvalidation.Validation)
}

func (r *Container) MakeView() contractsview.View {
	instance, err := r.Make(facades.FacadeToBinding[facades.View])
	if err != nil {
		color.Errorln(err)
		return nil
	}

	return instance.(contractsview.View)
}

func (r *Container) MakeWith(key any, parameters map[string]any) (any, error) {
	return r.make(key, parameters)
}

func (r *Container) Singleton(key any, callback func(app contractsfoundation.Application) (any, error)) {
	r.bindings.Store(key, instance{concrete: callback, shared: true})
}

func (r *Container) make(key any, parameters map[string]any) (any, error) {
	binding, ok := r.bindings.Load(key)
	if !ok {
		return nil, fmt.Errorf("binding not found: %+v", key)
	}

	if parameters == nil {
		instance, ok := r.instances.Load(key)
		if ok {
			return instance, nil
		}
	}

	bindingImpl := binding.(instance)
	switch concrete := bindingImpl.concrete.(type) {
	case func(app contractsfoundation.Application) (any, error):
		concreteImpl, err := concrete(App)
		if err != nil {
			return nil, err
		}
		if bindingImpl.shared {
			r.instances.Store(key, concreteImpl)
		}

		return concreteImpl, nil
	case func(app contractsfoundation.Application, parameters map[string]any) (any, error):
		concreteImpl, err := concrete(App, parameters)
		if err != nil {
			return nil, err
		}

		return concreteImpl, nil
	default:
		r.instances.Store(key, concrete)

		return concrete, nil
	}
}

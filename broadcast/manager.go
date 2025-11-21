package broadcast

import (
	"fmt"
	"sync"

	contractsbroadcast "github.com/goravel/framework/contracts/broadcast"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	contractslog "github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/errors"
)

// Ensure interface implementation
var _ contractsbroadcast.Manager = (*Manager)(nil)

// Manager manages broadcast connections and drivers.
type Manager struct {
	config config.Config
	json   foundation.Json
	log    contractslog.Log
	redis  any // Optional Redis instance

	drivers       map[string]contractsbroadcast.Broadcaster
	defaultDriver string

	mu sync.RWMutex
}

// NewManager creates a new broadcast manager instance.
func NewManager(config config.Config, json foundation.Json, log contractslog.Log, redis any) *Manager {
	defaultDriver := config.GetString("broadcasting.default", "null")

	return &Manager{
		config:        config,
		json:          json,
		log:           log,
		redis:         redis,
		drivers:       make(map[string]contractsbroadcast.Broadcaster),
		defaultDriver: defaultDriver,
	}
}

// Connection gets a broadcaster instance by name.
func (m *Manager) Connection(name ...string) (contractsbroadcast.Broadcaster, error) {
	driverName := m.defaultDriver
	if len(name) > 0 && name[0] != "" {
		driverName = name[0]
	}

	m.mu.RLock()
	driverInstance, exists := m.drivers[driverName]
	m.mu.RUnlock()
	if exists {
		return driverInstance, nil
	}

	err := m.registerDriver(driverName)
	if err != nil {
		return nil, err
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.drivers[driverName], nil
}

// Driver is an alias for Connection.
func (m *Manager) Driver(name ...string) (contractsbroadcast.Broadcaster, error) {
	return m.Connection(name...)
}

// registerDriver registers a broadcast driver based on configuration.
func (m *Manager) registerDriver(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring lock
	if _, exists := m.drivers[name]; exists {
		return nil
	}

	driver := m.config.GetString(fmt.Sprintf("broadcasting.connections.%s.driver", name))
	if driver == "" {
		return errors.BroadcastDriverNotSupported.Args(name, "driver not configured")
	}

	var broadcaster contractsbroadcast.Broadcaster
	var err error

	switch driver {
	case "pusher":
		broadcaster, err = m.createPusherDriver(name)
	case "ably":
		broadcaster, err = m.createAblyDriver(name)
	case "redis":
		broadcaster, err = m.createRedisDriver(name)
	case "log":
		broadcaster, err = m.createLogDriver(name)
	case "null":
		broadcaster = m.createNullDriver()
	case "custom":
		broadcaster, err = m.createCustomDriver(name)
	default:
		return errors.BroadcastDriverNotSupported.Args(driver)
	}

	if err != nil {
		return errors.BroadcastDriverRegisterFailed.Args(name, err)
	}

	m.drivers[name] = broadcaster
	return nil
}

// createPusherDriver creates a Pusher broadcaster instance.
func (m *Manager) createPusherDriver(name string) (contractsbroadcast.Broadcaster, error) {
	// Pusher is commented out as it's an optional dependency
	// Uncomment pusher_broadcaster.go and add pusher-http-go/v5 to go.mod to use
	return nil, fmt.Errorf("pusher driver not available - add pusher-http-go/v5 package to use")
}

// createAblyDriver creates an Ably broadcaster instance.
func (m *Manager) createAblyDriver(name string) (contractsbroadcast.Broadcaster, error) {
	key := m.config.GetString(fmt.Sprintf("broadcasting.connections.%s.key", name))
	if key == "" {
		return nil, fmt.Errorf("ably key not configured for connection: %s", name)
	}

	return NewAbly(key), nil
}

// createRedisDriver creates a Redis broadcaster instance.
func (m *Manager) createRedisDriver(name string) (contractsbroadcast.Broadcaster, error) {
	if m.redis == nil {
		return nil, fmt.Errorf("redis facade not available")
	}

	connection := m.config.GetString(fmt.Sprintf("broadcasting.connections.%s.connection", name), "default")
	prefix := m.config.GetString(fmt.Sprintf("broadcasting.connections.%s.prefix", name), "")

	return NewRedis(m.redis, connection, prefix), nil
}

// createLogDriver creates a Log broadcaster instance.
func (m *Manager) createLogDriver(name string) (contractsbroadcast.Broadcaster, error) {
	if m.log == nil {
		return nil, fmt.Errorf("log facade not available")
	}

	return NewLog(m.log), nil
}

// createNullDriver creates a Null broadcaster instance.
func (m *Manager) createNullDriver() contractsbroadcast.Broadcaster {
	return NewNull()
}

// createCustomDriver creates a custom broadcaster instance.
func (m *Manager) createCustomDriver(name string) (contractsbroadcast.Broadcaster, error) {
	via := m.config.Get(fmt.Sprintf("broadcasting.connections.%s.via", name))

	if custom, ok := via.(contractsbroadcast.Broadcaster); ok {
		return custom, nil
	}

	if factory, ok := via.(func() (contractsbroadcast.Broadcaster, error)); ok {
		return factory()
	}

	return nil, fmt.Errorf("custom driver must provide a Broadcaster instance or factory function")
}

// getDriverConfig gets the configuration for a specific driver.
func (m *Manager) getDriverConfig(name string) map[string]any {
	configKey := fmt.Sprintf("broadcasting.connections.%s", name)
	config := m.config.Get(configKey)

	if configMap, ok := config.(map[string]any); ok {
		return configMap
	}

	return make(map[string]any)
}

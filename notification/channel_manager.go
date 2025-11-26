package notification

import (
	"fmt"

	"github.com/goravel/framework/contracts/config"
	contractsorm "github.com/goravel/framework/contracts/database/orm"
	contractsevent "github.com/goravel/framework/contracts/event"
	contractslog "github.com/goravel/framework/contracts/log"
	contractsmail "github.com/goravel/framework/contracts/mail"
	contractsnotification "github.com/goravel/framework/contracts/notification"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/notification/channels"
)

// ChannelManager manages notification channels.
type ChannelManager struct {
	config         config.Config
	event          contractsevent.Instance
	log            contractslog.Log
	mail           contractsmail.Mail
	orm            contractsorm.Orm
	queue          contractsqueue.Queue
	channels       map[string]contractsnotification.Channel
	defaultChannel string
}

// NewChannelManager creates a new channel manager instance.
func NewChannelManager(
	config config.Config,
	event contractsevent.Instance,
	log contractslog.Log,
	mail contractsmail.Mail,
	orm contractsorm.Orm,
	queue contractsqueue.Queue,
) *ChannelManager {
	manager := &ChannelManager{
		config:         config,
		event:          event,
		log:            log,
		mail:           mail,
		orm:            orm,
		queue:          queue,
		channels:       make(map[string]contractsnotification.Channel),
		defaultChannel: "database",
	}

	// Register default channels
	manager.registerDefaultChannels()

	return manager
}

// registerDefaultChannels registers the built-in notification channels.
func (m *ChannelManager) registerDefaultChannels() {
	// Register database channel
	if m.orm != nil {
		m.channels["database"] = channels.NewDatabaseChannel(m.orm)
	}

	// Register mail channel
	if m.mail != nil {
		m.channels["mail"] = channels.NewMailChannel(m.mail)
	}
}

// Channel gets a channel instance by name.
func (m *ChannelManager) Channel(name string) (contractsnotification.Channel, error) {
	if channel, exists := m.channels[name]; exists {
		return channel, nil
	}
	return nil, fmt.Errorf("notification channel [%s] not found", name)
}

// Send sends the given notification to the given notifiable entities.
func (m *ChannelManager) Send(notifiables any, notification contractsnotification.Notification) error {
	sender := NewNotificationSender(m, m.queue, m.log)
	return sender.Send(notifiables, notification)
}

// SendNow sends the given notification immediately to the given notifiable entities.
func (m *ChannelManager) SendNow(notifiables any, notification contractsnotification.Notification, channels ...string) error {
	sender := NewNotificationSender(m, m.queue, m.log)
	return sender.SendNow(notifiables, notification, channels...)
}

// Extend registers a custom notification channel.
func (m *ChannelManager) Extend(name string, channel contractsnotification.Channel) {
	m.channels[name] = channel
}

// GetDefaultChannel returns the default channel name.
func (m *ChannelManager) GetDefaultChannel() string {
	return m.defaultChannel
}

// SetDefaultChannel sets the default channel name.
func (m *ChannelManager) SetDefaultChannel(channel string) {
	m.defaultChannel = channel
}

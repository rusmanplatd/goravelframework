package notification

import (
	"fmt"

	contractsorm "github.com/goravel/framework/contracts/database/orm"
	contractsnotification "github.com/goravel/framework/contracts/notification"
	databaseorm "github.com/goravel/framework/database/orm"
)

// Notifiable provides notification functionality to models.
type Notifiable struct {
	model any
	orm   contractsorm.Orm
}

// NewNotifiable creates a new notifiable instance.
func NewNotifiable(model any, orm contractsorm.Orm) *Notifiable {
	return &Notifiable{
		model: model,
		orm:   orm,
	}
}

// Notify sends a notification to the entity.
func (n *Notifiable) Notify(notification contractsnotification.Notification, manager contractsnotification.Factory) error {
	return manager.Send(n.model, notification)
}

// NotifyNow sends a notification immediately to the entity.
func (n *Notifiable) NotifyNow(notification contractsnotification.Notification, manager contractsnotification.Factory, channels ...string) error {
	return manager.SendNow(n.model, notification, channels...)
}

// Notifications returns a query builder for the entity's notifications.
func (n *Notifiable) Notifications() ([]databaseorm.DatabaseNotification, error) {
	var notifications []databaseorm.DatabaseNotification

	// Get notifiable type and ID
	notifiableType, notifiableID, err := n.getNotifiableInfo()
	if err != nil {
		return nil, err
	}

	// Query notifications
	err = n.orm.Query().
		Where("notifiable_type = ?", notifiableType).
		Where("notifiable_id = ?", notifiableID).
		Order("created_at DESC").
		Find(&notifications)

	return notifications, err
}

// ReadNotifications returns the entity's read notifications.
func (n *Notifiable) ReadNotifications() ([]databaseorm.DatabaseNotification, error) {
	var notifications []databaseorm.DatabaseNotification

	notifiableType, notifiableID, err := n.getNotifiableInfo()
	if err != nil {
		return nil, err
	}

	err = n.orm.Query().
		Where("notifiable_type = ?", notifiableType).
		Where("notifiable_id = ?", notifiableID).
		Where("read_at IS NOT NULL").
		Order("created_at DESC").
		Find(&notifications)

	return notifications, err
}

// UnreadNotifications returns the entity's unread notifications.
func (n *Notifiable) UnreadNotifications() ([]databaseorm.DatabaseNotification, error) {
	var notifications []databaseorm.DatabaseNotification

	notifiableType, notifiableID, err := n.getNotifiableInfo()
	if err != nil {
		return nil, err
	}

	err = n.orm.Query().
		Where("notifiable_type = ?", notifiableType).
		Where("notifiable_id = ?", notifiableID).
		Where("read_at IS NULL").
		Order("created_at DESC").
		Find(&notifications)

	return notifications, err
}

// RouteNotificationFor returns the notification routing information for the given channel.
func (n *Notifiable) RouteNotificationFor(channel string) any {
	// Check if model implements custom routing
	if routable, ok := n.model.(interface {
		RouteNotificationFor(channel string) any
	}); ok {
		return routable.RouteNotificationFor(channel)
	}

	// Default routing based on channel
	switch channel {
	case "database":
		return n.model
	case "mail":
		// Try to get email from model
		if emailer, ok := n.model.(interface{ GetEmail() string }); ok {
			return emailer.GetEmail()
		}
	}

	return nil
}

// getNotifiableInfo extracts type and ID information from the model.
func (n *Notifiable) getNotifiableInfo() (string, uint, error) {
	// This is a helper to get the model's type and ID
	// Implementation would depend on the model structure

	// Try to get ID from model
	if idGetter, ok := n.model.(interface{ GetID() uint }); ok {
		typeName := fmt.Sprintf("%T", n.model)
		return typeName, idGetter.GetID(), nil
	}

	return "", 0, fmt.Errorf("model must implement GetID() method or have ID field")
}

// MarkNotificationAsRead marks a specific notification as read.
func (n *Notifiable) MarkNotificationAsRead(notificationID string) error {
	var notification databaseorm.DatabaseNotification
	if err := n.orm.Query().Where("id = ?", notificationID).First(&notification); err != nil {
		return err
	}

	// Use ORM Update directly without type assertion
	_, err := n.orm.Query().Model(&notification).Where("id = ?", notificationID).Update("read_at", "NOW()")
	return err
}

// MarkAllNotificationsAsRead marks all notifications as read.
func (n *Notifiable) MarkAllNotificationsAsRead() error {
	notifiableType, notifiableID, err := n.getNotifiableInfo()
	if err != nil {
		return err
	}

	_, err = n.orm.Query().
		Model(&databaseorm.DatabaseNotification{}).
		Where("notifiable_type = ?", notifiableType).
		Where("notifiable_id = ?", notifiableID).
		Where("read_at IS NULL").
		Update("read_at", "NOW()")
	return err
}

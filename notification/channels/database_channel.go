package channels

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"

	"github.com/rusmanplatd/goravelframework/contracts/database/orm"
	contractsnotification "github.com/rusmanplatd/goravelframework/contracts/notification"
	databaseorm "github.com/rusmanplatd/goravelframework/database/orm"
)

// DatabaseChannel sends notifications to the database.
type DatabaseChannel struct {
	orm orm.Orm
}

// NewDatabaseChannel creates a new database channel instance.
func NewDatabaseChannel(orm orm.Orm) *DatabaseChannel {
	return &DatabaseChannel{
		orm: orm,
	}
}

// Send sends the given notification to the given notifiable entity.
func (c *DatabaseChannel) Send(notifiable any, notification contractsnotification.Notification) error {
	// Get the notification data
	var data map[string]any

	// Try ToDatabase first
	if dbMsg := notification.ToDatabase(notifiable); dbMsg != nil && dbMsg.Data != nil {
		data = dbMsg.Data
	} else {
		// Fallback to ToArray
		data = notification.ToArray(notifiable)
	}

	// Get notifiable type and ID
	notifiableType, notifiableID, err := c.getNotifiableInfo(notifiable)
	if err != nil {
		return err
	}

	// Get notification type
	notificationType := c.getNotificationType(notification)

	// Get notification ID
	var notificationID string
	if baseNotif, ok := notification.(interface{ ID() string }); ok {
		notificationID = baseNotif.ID()
	} else {
		notificationID = uuid.New().String()
	}

	// Create the database notification
	dbNotification := &databaseorm.DatabaseNotification{
		ID:             notificationID,
		Type:           notificationType,
		NotifiableType: notifiableType,
		NotifiableID:   notifiableID,
		Data:           data,
	}

	// Save to database
	return c.orm.Query().Create(dbNotification)
}

// getNotifiableInfo extracts the type and ID from the notifiable entity.
func (c *DatabaseChannel) getNotifiableInfo(notifiable any) (string, uint, error) {
	val := reflect.ValueOf(notifiable)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Get the type name
	typeName := val.Type().String()

	// Try to get ID field
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		return "", 0, fmt.Errorf("notifiable entity must have an ID field")
	}

	var id uint
	switch idField.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		id = uint(idField.Uint())
	default:
		return "", 0, fmt.Errorf("notifiable ID must be a uint type")
	}

	return typeName, id, nil
}

// getNotificationType returns the notification type name.
func (c *DatabaseChannel) getNotificationType(notification contractsnotification.Notification) string {
	t := reflect.TypeOf(notification)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.String()
}

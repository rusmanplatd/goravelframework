package orm

import (
	"gorm.io/gorm"

	"github.com/goravel/framework/support/carbon"
)

// DatabaseNotification represents a notification stored in the database.
type DatabaseNotification struct {
	ID             string           `gorm:"primaryKey;type:char(36)" json:"id"`
	Type           string           `gorm:"type:varchar(255);not null" json:"type"`
	NotifiableType string           `gorm:"type:varchar(255);not null;index:notifications_notifiable_index" json:"notifiable_type"`
	NotifiableID   uint             `gorm:"not null;index:notifications_notifiable_index" json:"notifiable_id"`
	Data           map[string]any   `gorm:"type:json;serializer:json" json:"data"`
	ReadAt         *carbon.DateTime `gorm:"type:timestamp" json:"read_at"`
	CreatedAt      *carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt      *carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

// TableName returns the table name for the DatabaseNotification model.
func (DatabaseNotification) TableName() string {
	return "notifications"
}

// Notifiable returns the notifiable entity that the notification belongs to.
func (n *DatabaseNotification) Notifiable() *gorm.DB {
	// This would need to be implemented based on the polymorphic relationship
	// For now, returning nil as it requires more context about the notifiable model
	return nil
}

// MarkAsRead marks the notification as read.
func (n *DatabaseNotification) MarkAsRead(db *gorm.DB) error {
	if n.ReadAt == nil {
		now := carbon.Now()
		readAt := carbon.DateTime{Carbon: now}
		n.ReadAt = &readAt
		return db.Model(n).Update("read_at", readAt).Error
	}
	return nil
}

// MarkAsUnread marks the notification as unread.
func (n *DatabaseNotification) MarkAsUnread(db *gorm.DB) error {
	if n.ReadAt != nil {
		n.ReadAt = nil
		return db.Model(n).Update("read_at", nil).Error
	}
	return nil
}

// Read returns true if the notification has been read.
func (n *DatabaseNotification) Read() bool {
	return n.ReadAt != nil
}

// Unread returns true if the notification has not been read.
func (n *DatabaseNotification) Unread() bool {
	return n.ReadAt == nil
}

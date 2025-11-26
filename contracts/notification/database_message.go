package notification

// DatabaseMessage represents the data structure for database notifications.
type DatabaseMessage struct {
	// Data is the notification data to be stored in the database.
	Data map[string]any
}

// NewDatabaseMessage creates a new database message.
func NewDatabaseMessage() *DatabaseMessage {
	return &DatabaseMessage{
		Data: make(map[string]any),
	}
}

// WithData sets the data for the database message.
func (m *DatabaseMessage) WithData(data map[string]any) *DatabaseMessage {
	m.Data = data
	return m
}

package console

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rusmanplatd/goravelframework/contracts/console"
	"github.com/rusmanplatd/goravelframework/contracts/console/command"
)

// NotificationMakeCommand creates a new notification class.
type NotificationMakeCommand struct {
}

// NewNotificationMakeCommand creates a new notification make command instance.
func NewNotificationMakeCommand() *NotificationMakeCommand {
	return &NotificationMakeCommand{}
}

// Signature returns the command signature.
func (c *NotificationMakeCommand) Signature() string {
	return "make:notification"
}

// Description returns the command description.
func (c *NotificationMakeCommand) Description() string {
	return "Create a new notification class"
}

// Extend extends the command.
func (c *NotificationMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Overwrite the notification if it already exists",
			},
		},
	}
}

// Handle handles the command.
func (c *NotificationMakeCommand) Handle(ctx console.Context) error {
	name := ctx.Argument(0)
	if name == "" {
		return fmt.Errorf("notification name is required")
	}

	// Ensure the name ends with "Notification"
	if !strings.HasSuffix(name, "Notification") {
		name = name + "Notification"
	}

	// Convert to snake case for file name (simple implementation)
	fileName := toSnakeCase(name) + ".go"

	// Default path
	path := filepath.Join("app", "notifications", fileName)

	// Check if file exists
	if _, err := os.Stat(path); err == nil && !ctx.OptionBool("force") {
		return fmt.Errorf("notification already exists: %s", path)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate the notification content
	content := c.generateNotificationContent(name)

	// Write the file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write notification file: %w", err)
	}

	ctx.Success(fmt.Sprintf("Notification created successfully: %s", path))
	return nil
}

// generateNotificationContent generates the notification class content.
func (c *NotificationMakeCommand) generateNotificationContent(name string) string {
	return fmt.Sprintf(`package notifications

import (
	"github.com/rusmanplatd/goravelframework/contracts/mail"
	contractsnotification "github.com/rusmanplatd/goravelframework/contracts/notification"
	"github.com/rusmanplatd/goravelframework/notification"
)

// %s represents a notification.
type %s struct {
	*notification.BaseNotification
}

// New%s creates a new notification instance.
func New%s() *%s {
	return &%s{
		BaseNotification: notification.NewBaseNotification(),
	}
}

// Via returns the channels the notification should be sent through.
func (n *%s) Via(notifiable any) []string {
	return []string{"database"}
}

// ToDatabase returns the database representation of the notification.
func (n *%s) ToDatabase(notifiable any) *contractsnotification.DatabaseMessage {
	return contractsnotification.NewDatabaseMessage().WithData(map[string]any{
		"message": "This is a notification",
	})
}

// ToMail returns the mail representation of the notification.
func (n *%s) ToMail(notifiable any) mail.Mailable {
	// Implement mail notification here
	return nil
}

// ToArray returns the array representation of the notification.
func (n *%s) ToArray(notifiable any) map[string]any {
	return map[string]any{
		"message": "This is a notification",
	}
}
`, name, name, name, name, name, name, name, name, name, name)
}

// toSnakeCase converts a string to snake_case.
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

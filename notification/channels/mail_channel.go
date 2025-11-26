package channels

import (
	"fmt"

	contractsmail "github.com/goravel/framework/contracts/mail"
	contractsnotification "github.com/goravel/framework/contracts/notification"
)

// MailChannel sends notifications via email.
type MailChannel struct {
	mail contractsmail.Mail
}

// NewMailChannel creates a new mail channel instance.
func NewMailChannel(mail contractsmail.Mail) *MailChannel {
	return &MailChannel{
		mail: mail,
	}
}

// Send sends the given notification to the given notifiable entity.
func (c *MailChannel) Send(notifiable any, notification contractsnotification.Notification) error {
	// Get the mailable from the notification
	mailable := notification.ToMail(notifiable)
	if mailable == nil {
		return fmt.Errorf("notification does not implement ToMail method")
	}

	// Convert to mail.Mailable if needed
	mailMailable, ok := mailable.(contractsmail.Mailable)
	if !ok {
		return fmt.Errorf("ToMail must return a mail.Mailable instance")
	}

	// Get the recipient email address
	recipient := c.getRecipient(notifiable)
	if recipient == "" {
		return fmt.Errorf("notifiable entity does not have an email address")
	}

	// Send the email
	return c.mail.To([]string{recipient}).Send(mailMailable)
}

// getRecipient extracts the email address from the notifiable entity.
func (c *MailChannel) getRecipient(notifiable any) string {
	// Try to get email from RouteNotificationForMail method
	if routable, ok := notifiable.(interface {
		RouteNotificationForMail() string
	}); ok {
		return routable.RouteNotificationForMail()
	}

	// Try to get email from Email field using reflection
	if emailer, ok := notifiable.(interface{ GetEmail() string }); ok {
		return emailer.GetEmail()
	}

	return ""
}

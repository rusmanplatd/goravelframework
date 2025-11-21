package broadcast

import (
	"encoding/json"
	"fmt"

	contractsbroadcast "github.com/goravel/framework/contracts/broadcast"
)

// Ensure interface implementations
var (
	_ contractsbroadcast.Channel = (*channel)(nil)
	_ contractsbroadcast.Channel = (*privateChannel)(nil)
	_ contractsbroadcast.Channel = (*presenceChannel)(nil)
)

// channel represents a public broadcast channel.
type channel struct {
	name string
}

// NewChannel creates a new public channel.
func NewChannel(name string) contractsbroadcast.Channel {
	return &channel{name: name}
}

func (c *channel) GetName() string {
	return c.name
}

func (c *channel) IsPrivate() bool {
	return false
}

func (c *channel) IsPresence() bool {
	return false
}

// privateChannel represents a private broadcast channel.
type privateChannel struct {
	name string
}

// NewPrivateChannel creates a new private channel.
func NewPrivateChannel(name string) contractsbroadcast.Channel {
	return &privateChannel{name: fmt.Sprintf("private-%s", name)}
}

func (p *privateChannel) GetName() string {
	return p.name
}

func (p *privateChannel) IsPrivate() bool {
	return true
}

func (p *privateChannel) IsPresence() bool {
	return false
}

// presenceChannel represents a presence broadcast channel.
type presenceChannel struct {
	name string
}

// NewPresenceChannel creates a new presence channel.
func NewPresenceChannel(name string) contractsbroadcast.Channel {
	return &presenceChannel{name: fmt.Sprintf("presence-%s", name)}
}

func (p *presenceChannel) GetName() string {
	return p.name
}

func (p *presenceChannel) IsPrivate() bool {
	return true
}

func (p *presenceChannel) IsPresence() bool {
	return true
}

// PresenceChannelMember represents a member in a presence channel.
type PresenceChannelMember struct {
	ID   string                 `json:"id"`
	Info map[string]interface{} `json:"info"`
}

// PresenceChannelData represents the data structure for presence channel authentication.
type PresenceChannelData struct {
	UserID   string                 `json:"user_id"`
	UserInfo map[string]interface{} `json:"user_info"`
}

// NewPresenceChannelMember creates a new presence channel member.
func NewPresenceChannelMember(id string, info map[string]interface{}) *PresenceChannelMember {
	return &PresenceChannelMember{
		ID:   id,
		Info: info,
	}
}

// NewPresenceChannelData creates new presence channel data for authentication.
func NewPresenceChannelData(userID string, userInfo map[string]interface{}) *PresenceChannelData {
	return &PresenceChannelData{
		UserID:   userID,
		UserInfo: userInfo,
	}
}

// ToJSON converts the presence channel member to JSON.
func (m *PresenceChannelMember) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// ToJSON converts the presence channel data to JSON.
func (d *PresenceChannelData) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

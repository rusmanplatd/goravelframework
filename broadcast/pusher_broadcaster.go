package broadcast

// NOTE: Pusher broadcaster is commented out because pusher-http-go/v5 is an optional dependency.
// To use Pusher, add the package to go.mod and uncomment this file.

/*
import (
	"fmt"

	contractsbroadcast "github.com/goravel/framework/contracts/broadcast"
	"github.com/pusher/pusher-http-go/v5"
)

// Ensure interface implementation
var _ contractsbroadcast.Broadcaster = (*Pusher)(nil)

// Pusher is a broadcaster that uses Pusher service.
type Pusher struct {
	*BaseBroadcaster
	client *pusher.Client
}

// NewPusher creates a new Pusher broadcaster.
func NewPusher(config map[string]any) (*Pusher, error) {
	// Extract configuration
	key, _ := config["key"].(string)
	secret, _ := config["secret"].(string)
	appID, _ := config["app_id"].(string)

	if key == "" || secret == "" || appID == "" {
		return nil, fmt.Errorf("pusher configuration incomplete: key, secret, and app_id are required")
	}

	// Extract options
	options := make(map[string]any)
	if opts, ok := config["options"].(map[string]any); ok {
		options = opts
	}

	cluster, _ := options["cluster"].(string)
	if cluster == "" {
		cluster = "mt1"
	}

	host, _ := options["host"].(string)
	scheme, _ := options["scheme"].(string)
	if scheme == "" {
		scheme = "https"
	}

	encrypted := true
	if enc, ok := options["encrypted"].(bool); ok {
		encrypted = enc
	}

	// Create Pusher client
	client := &pusher.Client{
		AppID:   appID,
		Key:     key,
		Secret:  secret,
		Cluster: cluster,
		Secure:  encrypted,
	}

	if host != "" {
		client.Host = host
	}

	if scheme == "http" {
		client.Secure = false
	}

	return &Pusher{
		BaseBroadcaster: NewBaseBroadcaster(),
		client:          client,
	}, nil
}

// Broadcast sends the event to Pusher channels.
func (p *Pusher) Broadcast(channels []contractsbroadcast.Channel, event string, payload map[string]any) error {
	if len(channels) == 0 {
		return nil
	}

	channelNames := p.formatChannels(channels)

	// Extract socket ID for exclusion
	var socketID string
	if socket, ok := payload["socket"].(string); ok {
		socketID = socket
		delete(payload, "socket")
	}

	// Pusher has a limit of 100 channels per request, so we chunk them
	chunkSize := 100
	for i := 0; i < len(channelNames); i += chunkSize {
		end := i + chunkSize
		if end > len(channelNames) {
			end = len(channelNames)
		}

		chunk := channelNames[i:end]

		// Trigger event on Pusher
		params := pusher.TriggerParams{}
		if socketID != "" {
			params.SocketID = &socketID
		}

		_, err := p.client.TriggerMulti(chunk, event, payload, params)
		if err != nil {
			return fmt.Errorf("failed to broadcast to pusher: %w", err)
		}
	}

	return nil
}

// GetClient returns the underlying Pusher client.
func (p *Pusher) GetClient() *pusher.Client {
	return p.client
}
*/

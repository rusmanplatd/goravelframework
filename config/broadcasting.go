package config

// Broadcasting returns the default broadcasting configuration.
func Broadcasting() map[string]any {
	return map[string]any{
		/*
			|--------------------------------------------------------------------------
			| Default Broadcaster
			|--------------------------------------------------------------------------
			|
			| This option controls the default broadcaster that will be used by the
			| framework when an event needs to be broadcast. You may set this to
			| any of the connections defined in the "connections" section below.
			|
			| Supported: "reverb", "pusher", "ably", "redis", "log", "null"
			|
		*/
		"default": "BROADCAST_DRIVER",

		/*
			|--------------------------------------------------------------------------
			| Broadcast Connections
			|--------------------------------------------------------------------------
			|
			| Here you may define all of the broadcast connections that will be used
			| to broadcast events to other systems or over websockets. Samples of
			| each available type of connection are provided inside this map.
			|
		*/
		"connections": map[string]any{
			"reverb": map[string]any{
				"driver": "reverb",
				"key":    "REVERB_APP_KEY",
				"secret": "REVERB_APP_SECRET",
				"app_id": "REVERB_APP_ID",
				"options": map[string]any{
					"host":   "REVERB_HOST",
					"port":   "REVERB_PORT",
					"scheme": "REVERB_SCHEME",
					"useTLS": true,
				},
			},

			"pusher": map[string]any{
				"driver": "pusher",
				"key":    "PUSHER_APP_KEY",
				"secret": "PUSHER_APP_SECRET",
				"app_id": "PUSHER_APP_ID",
				"options": map[string]any{
					"cluster":   "PUSHER_APP_CLUSTER",
					"host":      "PUSHER_HOST",
					"port":      "PUSHER_PORT",
					"scheme":    "PUSHER_SCHEME",
					"encrypted": true,
				},
			},

			"ably": map[string]any{
				"driver": "ably",
				"key":    "ABLY_KEY",
			},

			"redis": map[string]any{
				"driver":     "redis",
				"connection": "BROADCAST_REDIS_CONNECTION",
				"prefix":     "BROADCAST_REDIS_PREFIX",
			},

			"log": map[string]any{
				"driver": "log",
			},

			"null": map[string]any{
				"driver": "null",
			},
		},
	}
}

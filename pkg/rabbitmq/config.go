package rabbitmq

import "time"

type Config struct {
	AppName        string
	AmqpURL        string
	Timeout        time.Duration
	ConnectionName string
	Heartbeat      time.Duration
	Exchange       string
}

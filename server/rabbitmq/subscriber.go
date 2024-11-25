package rabbitmq

import (
	"github.com/skyapps-id/go-rabbitmq-exmaple/handler/rabbitmq/logs"
	"github.com/skyapps-id/go-rabbitmq-exmaple/pkg/rabbitmq"
)

func SetupSubscriber() {
	// Handler
	logshandler := logs.NewHandler()

	rabbitmq.Subscriber("logs.info", logshandler.Info)
	rabbitmq.Subscriber("logs.error", logshandler.Error)
	rabbitmq.Subscriber("logs.debug", logshandler.Debug)
}

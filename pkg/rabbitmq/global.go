package rabbitmq

var svcRabbitMQ SvcRabbitMQ

func Init(cfg Config) {
	svcRabbitMQ = New(cfg).InitializeConnection()
}

func Publish(topic string, payload string) {
	svcRabbitMQ.SendMessage(topic, payload)
}

func Subscriber(topic string, handler HandlerFunc) {
	svcRabbitMQ.CreateChannel(topic, handler)
}

func Close() {
	svcRabbitMQ.Close()
}

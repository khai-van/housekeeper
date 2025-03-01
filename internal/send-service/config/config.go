package config

type Config struct {
	RabbitMQURL   string `yaml:"rabbitmq_url"`
	RabbitMQQueue string `yaml:"rabbitmq_queue"`
}

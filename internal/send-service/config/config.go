package config

type Config struct {
	RabbitMQURL   string `mapstructure:"rabbitmqUrl"`
	RabbitMQQueue string `mapstructure:"rabbitmqQueue"`
}

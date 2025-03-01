package main

import "housekeeper/internal/send-service/config"

type Config struct {
	config.Config `mapstructure:",squash"`
	Port          int32 `mapstructure:"port"`
}

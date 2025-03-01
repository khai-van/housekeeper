package main

import "housekeeper/internal/send-service/config"

type Config struct {
	config.Config `yaml:",inline"`
	Port          int32 `yaml:"port"`
}

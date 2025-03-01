package utils

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Load config by file yaml and also parse env variable to that config file
func LoadConfig[T any](filename string) (*T, error) {
	confPath, err := filepath.Abs("./config")
	if err != nil {
		return nil, err
	}

	var config T

	viper.SetConfigName(filename)                          // set file name config
	viper.SetConfigType("yaml")                            // set file config type
	viper.AddConfigPath(confPath)                          // set path to config
	viper.SetEnvPrefix("")                                 // no need prefix in env
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // child and parent field in the config parse to env separate by _
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

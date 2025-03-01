package config

type Config struct {
	Port                  int    `mapstructure:"port"`
	MongoDBURI            string `mapstructure:"mongodbUri"`
	MongoDBName           string `mapstructure:"mongodbName"`
	PricingServiceAddress string `mapstructure:"pricingServiceAddress"`
	SendServiceAddress    string `mapstructure:"sendServiceAddress"`
}

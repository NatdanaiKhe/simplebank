package util

import "github.com/spf13/viper"

type Config struct {
	DB_Driver      string `mapstructure:"DB_DRIVER"`
	DB_URL         string `mapstructure:"DB_URL"`
	SERVER_ADDRESS string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // json, xml
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&config)
	return config, err
}

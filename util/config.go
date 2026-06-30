package util

import "github.com/spf13/viper"

type Config struct {
	DB_Driver      string `mapstructure:"DB_DRIVER"`
	DB_URL         string `mapstructure:"DB_URL"`
	SERVER_ADDRESS string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env") // json, xml
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, err
		}
	}

	err = viper.Unmarshal(&config)
	return config, err
}

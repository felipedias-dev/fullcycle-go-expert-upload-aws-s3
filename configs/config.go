package configs

import (
	"github.com/spf13/viper"
)

type config struct {
	AwsKey    string `mapstructure:"AWS_KEY"`
	AwsSecret string `mapstructure:"AWS_SECRET"`
	S3Bucket  string `mapstructure:"S3_BUCKET"`
	AwsRegion string `mapstructure:"AWS_REGION"`
}

func LoadConfig(path string) (*config, error) {
	var conf *config

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}

	return conf, nil
}

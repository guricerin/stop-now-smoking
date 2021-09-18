package util

import "github.com/spf13/viper"

type Config struct {
	ServerHost   string `mapstructure:"server_host"`
	ServerPort   string `mapstructure:"server_port"`
	DbUser       string `mapstructure:"db_user"`
	DbPassword   string `mapstructure:"db_password"`
	DbName       string `mapstructure:"db_name"`
	DbProtocol   string `mapstructure:"db_protocol"`
	DbConnOption string `mapstructure:"db_conn_option"`
}

func LoadConfig(filepath string) (cfg Config, err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return
	}
	return
}

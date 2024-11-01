package config

import "github.com/spf13/viper"

type Config struct {
	SECRETKEY    string `mapstructure:"JWTSECRET"`
	DBHost       string `mapstructure:"DBHOST"`
	DBUser       string `mapstructure:"DBUSER"`
	DBPassword   string `mapstructure:"DBPASSWORD"`
	DBDatabase   string `mapstructure:"DBNAME"`
	DBPort       string `mapstructure:"DBPORT"`
	DBSslmode    string `mapstructure:"DBSSL"`
	REDISHOST    string `mapstructure:"REDISHOST"`
	GrpcUserPort string `mapstructure:"GRPCUSERPORT"`
}

func LoadConfig() *Config {
	var cnfg Config
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.Unmarshal(&cnfg)
	return &cnfg
}

package options

import (
	"fmt"
	configs "github.com/iwinder/geekGoWork/internal/pkg/options"
	"github.com/spf13/viper"
)

func InitConfig() *configs.Option {
	viper.SetConfigType("yaml")
	viper.SetConfigName("week04")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	var config *configs.Option
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("config.server: %#v\n", config.ServerOption)
	fmt.Printf("config.mysql: %#v\n", config.MysqlOption)
	return config
}

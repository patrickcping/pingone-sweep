package config

import "github.com/spf13/viper"

var config *viper.Viper

func Init() {
	config = viper.New()

}

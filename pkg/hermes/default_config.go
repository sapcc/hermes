package hermes

import "github.com/spf13/viper"

func SetDefaultConfig() {
	viper.SetDefault("hermes.keystone_driver", "real")
	viper.SetDefault("API.ListenAddress", "0.0.0.0:8788")
}

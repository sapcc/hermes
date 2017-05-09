package hermes

import "github.com/spf13/viper"

func SetDefaultConfig() {
	viper.SetDefault("hermes.keystone_driver", "keystone")
	viper.SetDefault("hermes.storage_driver", "mock")
	viper.SetDefault("API.ListenAddress", "0.0.0.0:8788")
}

package hermes

import "github.com/spf13/viper"

func SetDefaultConfig() {
	viper.SetDefault("API.ListenAddress", "0.0.0.0:8788")
}

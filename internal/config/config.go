package config

import (
	"strings"

	"github.com/NataliaZabelina/monitoring/internal/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Host      string        `json:"host" mapstructure:"host"`
	Port      string        `json:"port" mapstructure:"port"`
	Log       logger.Params `json:"log" mapstructure:"log"`
	Collector Collector     `json:"collector" mapstructure:"collector"`
}

type Collector struct {
	Timeout    int        `json:"timeout" mapstructure:"timeout"`
	Statistics Statistics `json:"statistics" mapstructure:"statistics"`
}

type Statistics struct {
	LoadSystem  bool `json:"load_system" mapstructure:"load_system"`
	LoadCPU     bool `json:"load_cpu" mapstructure:"load_cpu"`
	LoadDisk    bool `json:"load_disk" mapstructure:"load_disk"`
	TopTalkers  bool `json:"top_talkers" mapstructure:"top_talkers"`
	NetworkStat bool `json:"stat_network" mapstructure:"network_stat"`
}

func LoadConfig(filePath string) (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("host", "localhost")
	viper.SetDefault("collector.timeout", 5)

	if filePath == "" {
		logger.Logger.Info("Config file path was not provided.")
	} else {
		viper.SetConfigFile(filePath)
		viper.SetConfigType("json")

		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

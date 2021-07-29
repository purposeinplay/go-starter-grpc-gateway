package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config the config.json file should be set at the root level
type Config struct {
	SERVER struct {
		Port    int    `mapstructure:"port"`
		Address string `mapstructure:"address"`
	}
	DB struct {
		Driver      string `mapstructure:"driver"`
		HOST        string `mapstructure:"host"`
		USER        string `mapstructure:"user"`
		PASSWORD    string `mapstructure:"password"`
		NAME        string `mapstructure:"name"`
		Automigrate bool   `mapstructure:"automigrate"`
	}

	JWT struct {
		Secret          string `mapstructure:"secret"`
		RefreshTokenExp int    `mapstructure:"refresh_token_exp"`
		AccessTokenExp  int    `mapstructure:"access_token_exp"`
	}
}

var ConfigFile string

func LoadTestConfig(path string) (*Config, error) {
	cmd := &cobra.Command{}
	cmd.Flags().StringVar(&ConfigFile, "config", path, "Config file")

	return LoadConfig(cmd)
}

// LoadConfig should load and unmarshal the config file
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	err := viper.BindPFlags(cmd.Flags())

	if err != nil {
		return nil, err
	}

	viper.SetEnvPrefix("GOSTARTER")

	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath("./") // adding home directory as first search path
		viper.SetConfigName("config")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	config := new(Config)

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}

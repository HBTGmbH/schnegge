package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"schnegge/internal/base"
)

var homePath string

func init() {
	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case "windows":
		homePath = os.Getenv("HOMEPATH")
	default:
		homePath = os.Getenv("HOME")
	}
}

func ReadConfig() (Config, error) {
	initViper()
	var cfg Config = NewSimpleConfig()
	err := readConfigFile(cfg)
	return cfg, err
}

func WriteConfigFile(cfg Config) error {
	if value, ok := cfg.GetValue(TokenID); ok {
		viper.Set("TokenID", value)
	}
	if value, ok := cfg.GetValue(TokenSecret); ok {
		viper.Set("TokenSecret", value)
	}
	if value, ok := cfg.GetValue(Server); ok {
		viper.Set("Server", value)
	}
	if value, ok := cfg.GetValue(NoSplash); ok {
		viper.Set("NoSplash", value)
	}
	if value, ok := cfg.GetValue(Order); ok {
		viper.Set("Order", value)
	}
	return viper.WriteConfigAs(homePath + "/.schnegge")
}

func PrintVerboseConfig() {
	base.Log.Println()
	base.Log.Println("=== CONFIG ===")
	base.Log.Println("TokenID: ", viper.GetString("TokenID"))
	base.Log.Println("TokenSecret: ", viper.GetString("TokenSecret"))
	base.Log.Println("Server: ", viper.GetString("Server"))
	base.Log.Println("NoSplash: ", viper.GetString("NoSplash"))
	base.Log.Println("Auftrag: ", viper.GetString("Order"))
}

func initViper() {
	viper.SetConfigName(".schnegge")
	viper.AddConfigPath(homePath)
	viper.SetConfigType("yaml")
	viper.SetConfigPermissions(os.FileMode(0600))
}

func readConfigFile(cfg Config) error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("No .schnegge.*extension* File found")
			return nil
		} else {
			// Config file was found but another error was produced
			return err
		}
	}
	if viper.GetString("TokenID") != "" {
		cfg.AddValue(TokenID, viper.GetString("TokenID"))
	}
	if viper.GetString("TokenSecret") != "" {
		cfg.AddValue(TokenSecret, viper.GetString("TokenSecret"))
	}
	if viper.GetString("Server") != "" {
		cfg.AddValue(Server, viper.GetString("Server"))
	}
	if viper.GetString("NoSplash") != "" {
		cfg.AddValue(NoSplash, viper.GetString("NoSplash"))
	}
	if viper.GetString("Order") != "" {
		cfg.AddValue(Order, viper.GetString("Order"))
	}
	return nil
}

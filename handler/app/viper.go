package app

import (
	"chera_khube/internal/helper"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func SetupViper(path string) (*helper.ServiceConfig, error) {
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	f, err := os.Open(path)
	if err != nil {
		msg := fmt.Sprintf("cannot read config file: %s", err.Error())
		return nil, errors.New(msg)
	}
	err = viper.ReadConfig(f)
	if err != nil {
		msg := fmt.Sprintf("viper read config error: %s", err.Error())
		return nil, errors.New(msg)
	}
	var c helper.ServiceConfig
	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

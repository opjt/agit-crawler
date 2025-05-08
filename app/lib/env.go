package lib

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Env struct {
	Server ServerConfig
	Log    LogConfig
	Agit   Agit
}

type ServerConfig struct {
	Port        string `mapstructure:"PORT"`
	Environment string `mapstructure:"ENV"`
	Url         string `mapstructure:"URL"`
}

type LogConfig struct {
	Output string `mapstructure:"OUTPUT"`
	Level  string `mapstructure:"LEVEL"`
}

type Agit struct {
	Name     string `mapstructure:"NAME"`
	UserId   string `mapstructure:"USERID"`
	Password string `mapstructure:"PASSWORD"`
}

var (
	once sync.Once
	env  Env
)

func LoadEnv() (Env, error) {
	viper.SetConfigFile(".env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("☠️ cannot read configuration")
		return env, err
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("☠️ environment can't be loaded: ", err)
		return env, err
	}

	if env.Server.Port != "" && env.Server.Url != "" {
		env.Server.Url = fmt.Sprintf("%s:%s", env.Server.Url, env.Server.Port)
	}

	return env, nil
}

func NewEnv() Env {

	once.Do(func() {
		env, _ = LoadEnv()

	})
	return env
}

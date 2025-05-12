package lib

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Environment string

const (
	EnvDev  Environment = "dev"
	EnvProd Environment = "prod"
)

type Env struct {
	Server ServerConfig
	Log    LogConfig
	Agit   Agit
}

type ServerConfig struct {
	Port        string      `mapstructure:"PORT"`
	Environment Environment `mapstructure:"ENV"`
	Url         string      `mapstructure:"URL"`
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

func LoadEnv() Env {
	envFile := ".env"
	envSuffix := os.Getenv("APP_ENV")
	if envSuffix != "" {
		envFile = fmt.Sprintf(".env.%s", envSuffix) // 예: ".env.dev"
	}

	viper.SetConfigFile(envFile)
	viper.SetConfigType("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("☠️ cannot read configuration")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("☠️ environment can't be loaded: ", err)
	}
	// 환경 검증: dev 또는 prod 만 허용
	if env.Server.Environment != EnvDev && env.Server.Environment != EnvProd {
		log.Fatalf("☠️ Invalid environment value: %s. Only 'dev' and 'prod' are allowed.", env.Server.Environment)
	}

	if env.Server.Port != "" && env.Server.Url != "" {
		env.Server.Url = fmt.Sprintf("%s:%s", env.Server.Url, env.Server.Port)
	}

	return env
}

func NewEnv() Env {

	once.Do(func() {
		env = LoadEnv()

	})
	return env
}

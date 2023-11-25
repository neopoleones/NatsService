package config

import (
	"L0_task/pkg/utils"
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const (
	flagConfigPath    = "config"
	envConfigPath     = "CONFIG_PATH"
	defaultConfigPath = "./etc/default.yml"
)

type Config struct {
	Env string `yaml:"env"`

	Subscriber Subscriber `yaml:"nats"`
	Postgres   Postgres   `yaml:"postgres"`
	Service    Service    `yaml:"service"`
}

type Subscriber struct {
	Stream        string        `yaml:"stream"`
	Address       string        `yaml:"address"`
	ReconnectWait time.Duration `yaml:"reconnect_wait"`
	MaxReconnect  int           `yaml:"max_reconnect"`
}

type Postgres struct {
	Address  string `yaml:"address"`
	User     string `yaml:"user"`
	Password string `env:"PG_PASSWORD" env-required:"true"`
	Database string `yaml:"database"`
}

func (cfg Postgres) AsSchema() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Address,
		cfg.Database,
	)
}

type Service struct {
	Port int `yaml:"port"`
}

func MustLoad() *Config {
	var config Config
	path := getConfigPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Errorf("failed to load configuration file: %v", err))
	}

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		desc, _ := cleanenv.GetDescription(&config, nil)
		fmt.Println(desc)

		panic(fmt.Errorf("failed to load configuration file: %v", err))
	}

	return &config
}

// getConfigPath tries to get config path from several sources in priority:
// flag > env > default
func getConfigPath() string {
	var path string

	flag.StringVar(&path, flagConfigPath, "", "Specifies the path of configuration file")
	flag.Parse()

	if path == "" {
		path = utils.GetEnvDefault(envConfigPath, defaultConfigPath)
	}

	return path
}

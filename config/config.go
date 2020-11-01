package config

import (
	"io/ioutil"

	"github.com/jbuchbinder/swarm-monitor/logging"
	"gopkg.in/yaml.v2"
)

var (
	Config           *SwarmConfig
	cachedConfigPath string
)

type SwarmConfig struct {
	Debug                bool   `yaml:"debug"`
	LogLevel             string `yaml:"log-level"`
	Port                 int    `yaml:"port"`
	HostID               int    `yaml:"host-id"`
	PoolSize             int    `yaml:"pool-size"`
	DefaultCheckInterval uint64 `yaml:"check-interval"`
	Database             struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		Db   int    `yaml:"db"`
	} `yaml:"database"`
}

func (c *SwarmConfig) SetDefaults() {
	c.Debug = true
	c.LogLevel = "DEBUG"
	c.Port = 48666
	c.HostID = 1
	c.PoolSize = 10
	c.Database.Host = "localhost"
	c.Database.Port = 6379
	c.Database.Db = 13
	c.DefaultCheckInterval = 300
}

func LoadConfigWithDefaults(configPath string) (*SwarmConfig, error) {
	cachedConfigPath = configPath
	c := &SwarmConfig{}
	c.SetDefaults()
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal([]byte(data), c)
	if err == nil {
		logging.LogLevel = logging.StringToLevel(c.LogLevel)
	}
	return c, err
}

func ConfigReload() error {
	c := &SwarmConfig{}
	c.SetDefaults()
	data, err := ioutil.ReadFile(cachedConfigPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(data), c)
	if err == nil {
		Config = c
	}
	if err == nil {
		logging.LogLevel = logging.StringToLevel(c.LogLevel)
	}
	return err
}

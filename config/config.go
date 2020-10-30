package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	Config           *SwarmConfig
	cachedConfigPath string
)

type SwarmConfig struct {
	Debug    bool `yaml:"debug"`
	Port     int  `yaml:"port"`
	HostID   int  `yaml:"host-id"`
	PoolSize int    `yaml:"pool-size"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Db       int    `yaml:"db"`
	} `yaml:"database"`
}

func (c *SwarmConfig) SetDefaults() {
	c.Debug = true
	c.Port = 48666
	c.HostID = 1
	c.PoolSize = 10
	c.Database.Host = "localhost"
	c.Database.Port = 6379
	c.Database.Db = 5
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
	return err
}

package core

import (
	"github.com/BurntSushi/toml"
	"github.com/Unknwon/com"
	"os"
	"time"
)

var ConfigFile string = "config.toml"

type coreConfig struct {
	InstallTime int64  `toml:"install_time"`
	LogFile     string `toml:"log_file"`

	App    coreVars     `toml:"app"`
	Server serverConfig `toml:"server"`
	Db     dbConfig     `toml:"db"`
}

type serverConfig struct {
	Addr     string
	Port     string
	Protocol string
}

type dbConfig struct {
	Directory string
	Index     int
}

func (c *coreConfig) WriteFile() error {
	c.InstallTime = time.Now().Unix()
	f, err := os.OpenFile(ConfigFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	encoder := toml.NewEncoder(f)
	return encoder.Encode(c)
}

func (c *coreConfig) ReadFile() error {
	_, err := toml.DecodeFile(ConfigFile, c)
	return err
}

func (c *coreConfig) HasFile() bool {
	return com.IsFile(ConfigFile)
}

func newCoreConfig() *coreConfig {
	cfg := &coreConfig{
		App: *Vars,
		Server: serverConfig{
			Addr:     "0.0.0.0",
			Port:     "3030",
			Protocol: "http",
		},
		Db: dbConfig{
			Directory: "data",
			Index:     14,
		},
	}

	// if config file exist, load it
	if cfg.HasFile() {
		cfg.ReadFile()
	}
	return cfg
}

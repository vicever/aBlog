package core

import (
	"github.com/BurntSushi/toml"
	"github.com/Unknwon/com"
	"os"
)

var configFile string = "config.toml"

type CoreConfig struct {
	InstallTime int64

	TcpAddress string

	DbDriver     string
	DbSQLiteFile string
}

func NewConfig() *CoreConfig {
	return &CoreConfig{
		TcpAddress:   "0.0.0.0:3999",
		DbDriver:     "sqlite3",
		DbSQLiteFile: "sqlite.db",
	}
}

// check config in file or not
func (cfg *CoreConfig) IsFiled() bool {
	return com.IsFile(configFile)
}

// write config to file
func (cfg *CoreConfig) WriteFile() error {
	// open file handler
	f, err := os.OpenFile(configFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	// use toml encoder
	enc := toml.NewEncoder(f)
	return enc.Encode(cfg)
}

// read config from file
func (cfg *CoreConfig) ReadFile() error {
	_, err := toml.DecodeFile(configFile, cfg)
	return err
}

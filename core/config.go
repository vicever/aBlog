package core

import (
	"github.com/BurntSushi/toml"
	"github.com/Unknwon/com"
	"os"
	"time"
)

var Config *config = &config{
	App: appConfig{
		Name:        BLOG_NAME,
		Desc:        BLOG_DESC,
		Version:     BLOG_VERSION,
		VersionDate: BLOG_VERSION_DATE,
		Author:      BLOG_AUTHOR,
		Github:      BLOG_GITHUB,
		Official:    BLOG_OFFICIAL,
	},
	Server: serverConfig{
		Addr:     "0.0.0.0",
		Port:     3003,
		Protocol: "http",
	},
}

type config struct {
	App       appConfig    `toml:"app"`
	Server    serverConfig `toml:"server"`
	Installed time.Time
}

type appConfig struct {
	Name        string `toml:"name"`
	Desc        string `toml:"desc"`
	Version     string `toml:"version"`
	VersionDate string `toml:"version_date"`
	Author      string `toml:"author"`
	Github      string `toml:"github"`
	Official    string `toml:"official"`
}

type serverConfig struct {
	Addr     string `toml:"addr"`
	Port     int    `toml:"port"`
	Protocol string `toml:"protocol"`
}

func (c *config) WriteFile() error {
	c.Installed = time.Now()
	f, err := os.OpenFile(CONFIG_FILE, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	encoder := toml.NewEncoder(f)
	return encoder.Encode(c)
}

func (c *config) ReadFile() error {
	_, err := toml.DecodeFile(CONFIG_FILE, c)
	return err
}

func (c *config) HasFile() bool {
	return com.IsFile(CONFIG_FILE)
}

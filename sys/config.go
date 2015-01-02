package sys

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/Unknwon/com"
	"io/ioutil"
	"os"
	"time"
)

var Config *sysConfig

type sysConfig struct {
	fileName string
	IsLoad   bool `toml:"-"`
	InitTime time.Time
	NoDb     sysNoDbConfig  `toml:"nodb"`
	Tango    sysTangoConfig `toml:"tango"`
}

type sysNoDbConfig struct {
	Directory string
	Index     int
}

type sysTangoConfig struct {
	AddrString string `toml:"addr_string"`
}

func (s *sysConfig) setDefault() {
	s.NoDb = sysNoDbConfig{
		Directory: "storage",
		Index:     1,
	}
	s.Tango = sysTangoConfig{
		AddrString: "0.0.0.0:3003",
	}
	s.IsLoad = false
}

func (s *sysConfig) ToFile() {

	var buff bytes.Buffer
	encoder := toml.NewEncoder(&buff)
	if err := encoder.Encode(s); err != nil {
		Fatal("[Config] encode error : %v", err)
		return
	}

	err := ioutil.WriteFile(s.fileName, buff.Bytes(), os.ModePerm)
	if err != nil {
		Fatal("[Config] config to file error : %v")
		return
	}
}

func newSysConfig() *sysConfig {
	s := new(sysConfig)
	s.fileName = "ablog.toml"
	s.setDefault()

	// load from config file
	if com.IsFile(s.fileName) {
		_, err := toml.DecodeFile(s.fileName, s)
		if err != nil {
			Fatal("[Config] load config file error : %v", err)
		}
		s.IsLoad = true
	}
	return s
}

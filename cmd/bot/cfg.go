package main

import (
	"flag"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type (
	Duration struct {
		time.Duration
	}

	Config struct {
		TgAPI struct {
			Token string `toml:"token"`
		} `toml:"tg_api"`
		NotionAPI struct {
			Token          string    `toml:"token"`
			Timeout        *Duration `toml:"timeout"`
			Version        string    `toml:"version"`
			DbId           string    `toml:"db_id"`
			SearchUrl      string    `toml:"search_url"`
			UpdateInterval *Duration `toml:"update_interval"`
		} `toml:"notion_api"`
		Path string
	}
)

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func ConfigInit() *Config {
	fCfgPath := flag.String("c", "conf.toml", "path to conf file")
	flag.Parse()

	conf := new(Config)
	file, err := os.Open(*fCfgPath)
	eh(err)

	defer func() {
		if file == nil {
			return
		}
		eh(file.Close())
	}()

	_, err = toml.DecodeFile(*fCfgPath, &conf)
	eh(err)
	conf.Path = *fCfgPath
	return conf
}

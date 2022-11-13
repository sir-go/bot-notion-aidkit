package main

import (
	"time"

	"github.com/BurntSushi/toml"
)

type (
	Duration struct {
		time.Duration
	}

	CfgWarehouse struct {
		UpdateInterval *Duration `toml:"update_interval"`
	}

	CfgTelegram struct {
		Token string `toml:"token"`
	}

	CfgNotionAPI struct {
		Token     string    `toml:"token"`
		Timeout   *Duration `toml:"timeout"`
		Version   string    `toml:"version"`
		DbId      string    `toml:"db_id"`
		SearchUrl string    `toml:"search_url"`
	}

	CfgApp struct {
		TgAPI     CfgTelegram  `toml:"tg_api"`
		Wh        CfgWarehouse `toml:"warehouse"`
		NotionAPI CfgNotionAPI `toml:"notion_api"`
		Path      string
	}
)

func (d *Duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return
}

func LoadConfig(fileName string) (conf *CfgApp, err error) {
	if _, err = toml.DecodeFile(fileName, &conf); err != nil {
		return nil, err
	}
	conf.Path = fileName
	return
}

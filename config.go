package main

var GlobCfg = Config{}

type Config struct {
	PORT         int64    `toml:"port"`
	ALLOW_ORIGIN []string `toml:"allow_origin"`
	TG_KEY       string   `toml:"tg_key"`
}

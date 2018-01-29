package config

var GlobCfg = Config{}

type Config struct {
	PORT             int64    `toml:"port"`
	ALLOW_ORIGIN     []string `toml:"allow_origin"`
	API_KEY          string   `toml:"api_key"`
	TG_KEY           string   `toml:"tg_key"`
	HEARTBEAT_ALARM  int64    `toml:"heartbeat_alarm_val"`
	MONITOR_INTERVAL int64    `toml:"monitor_interval"`
	MANAGER_NAME     string   `toml:"manager_name"`
	DB_NAME          string   `toml:"db_name"`
	DB_USER          string   `toml:"db_user"`
	DB_PASS          string   `toml:"db_pass"`
	DB_CHARSET       string   `toml:"db_charset"`
	ADMIN            []Admin  `toml:"admin"`
}

type Admin struct {
	Username string
	Password string
}

func ParseDSN(config Config) string {
	return config.DB_USER + ":" + config.DB_PASS + "@/" + config.DB_NAME + "?charset=" + config.DB_CHARSET + "&parseTime=true"
}

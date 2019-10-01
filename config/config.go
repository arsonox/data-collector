package config

import (
	"sync"

	"github.com/JeremyLoy/config"
)

type Config struct {
	InfluxIP   string `config:"INFLUX_IP"`
	InfluxDB   string `config:"UNFLUX_DB"`
	UseSunspec bool   `config:"USE_SUNSPEC"`
	SunspecIP  string `config:"SUNSPEC_IP"`
	UseYouless bool   `config:"USE_YOULESS"`
	YoulessIP  string `config:"YOULESS_IP"`
	UseSolcast bool   `config:"USE_SOLCAST"`
	SolcastURL string `config:"SOLCAST_URL"`
}

var cfg Config
var mux sync.RWMutex

func ReloadConfig() {
	mux.Lock()
	config.From("config.env").FromEnv().To(&cfg)
	mux.Unlock()
}

func init() {
	ReloadConfig()
}

func GetConfig() Config {
	mux.RLock()
	defer mux.RUnlock()
	return cfg
}

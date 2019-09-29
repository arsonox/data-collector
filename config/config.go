package config

import (
	"sync"

	"github.com/JeremyLoy/config"
)

type Config struct {
	UseSunspec bool   `config:"USE_SUNSPEC"`
	SunspecIp  string `config:"SUNSPEC_IP"`
	UseYouless bool   `config:"USE_YOULESS"`
	YoulessIp  string `config:"YOULESS_IP"`
	UseSolcast bool   `config:"USE_SOLCAST"`
	SolcastUrl string `config:"SOLCAST_URL"`
}

var cfg Config
var mux sync.RWMutex

func ReloadConfig() {
	mux.Lock()
	config.FromEnv().To(&cfg)
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

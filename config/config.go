package config

import (
	"sync"

	"github.com/JeremyLoy/config"
)

//Config contains the configuration for the application
type Config struct {
	InfluxIP    string `config:"INFLUX_IP"`
	InfluxDB    string `config:"UNFLUX_DB"`
	InfluxDummy bool   `config:"INFLUX_DUMMY"`
	UseSunspec  bool   `config:"USE_SUNSPEC"`
	SunspecIP   string `config:"SUNSPEC_IP"`
	UseYouless  bool   `config:"USE_YOULESS"`
	YoulessIP   string `config:"YOULESS_IP"`
	UseSolcast  bool   `config:"USE_SOLCAST"`
	SolcastURL  string `config:"SOLCAST_URL"`
	SolcastKey  string `config:"SOLCAST_API_KEY"`
}

var cfg Config
var mux sync.RWMutex

//ReloadConfig reloads the configuration
func ReloadConfig() {
	mux.Lock()
	config.From("config.env").FromEnv().To(&cfg)
	mux.Unlock()
}

func init() {
	ReloadConfig()
}

//GetConfig gets a copy of the current configuration variables
func GetConfig() Config {
	mux.RLock()
	defer mux.RUnlock()
	return cfg
}

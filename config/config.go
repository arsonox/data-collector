package config

var Config struct {
	UseSunspec bool   `config:"USE_SUNSPEC"`
	SunspecIp  string `config:"SUNSPEC_IP"`
	UseYouless bool   `config:"USE_YOULESS"`
	YoulessIp  string `config:"YOULESS_IP"`
	UseSolcast bool   `config:"USE_SOLCAST"`
	SolcastUrl string `config:"SOLCAST_URL"`
}

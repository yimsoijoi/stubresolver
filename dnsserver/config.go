package dnsserver

type Config struct {
	Address  string `mapstructure:"address" json:"address"`
	Protocol string `mapstructure:"protocol" json:"protocol"`
}

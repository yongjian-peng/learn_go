package config

type Server struct {
	Zap      Zap      `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis    Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	AutoCode Autocode `mapstructure:"autoCode" json:"autoCode" yaml:"autoCode"`
	Mysql    Mysql    `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Notify   Notify   `mapstructure:"notify" json:"notify" yaml:"notify"`
}

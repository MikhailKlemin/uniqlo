package config

//DefaultConfig is
type DefaultConfig struct {
	DB string
}

//NewDefaultConfig is
func NewDefaultConfig() DefaultConfig {
	var d DefaultConfig
	d.DB = "/media/mike/WDC4_1/Neo/uniclo.uk/assets/uniqlo.v4.sqlite"
	//d.Cache
	return d
}

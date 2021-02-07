package config

import "time"

//DefaultConfig is
type DefaultConfig struct {
	DB   string
	Date time.Time
	Age  int //Age to process
}

//NewDefaultConfig is
func NewDefaultConfig() DefaultConfig {
	var d DefaultConfig
	d.DB = "/media/mike/WDC4_1/Neo/uniclo.uk/assets/uniqlo.v4.sqlite"
	d.Date = time.Now().Add(-24 * time.Hour)
	d.Age = 40
	//d.Cache
	return d
}

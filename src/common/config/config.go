package config

type Config interface {
	Init()
}

type config struct {
	cfg []Config
}

func NewConfig(handlers []func() Config) Config {
	var cfg []Config
	for _, handler := range handlers {
		cfg = append(cfg, handler())
	}

	return &config{cfg: cfg}
}

func (c *config) Init() {
	for _, cfg := range c.cfg {
		cfg.Init()
	}
}

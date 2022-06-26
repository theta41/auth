package cfg

type Config struct {
	JWTSecret string `yaml:"jwt_secret"`
	// JWTTTL stored in seconds.
	JWTTTL int `yaml:"jwt_ttl"`

	HostAddress    string `yaml:"host_address"`
	MetricsAddress string `yaml:"metrics_address"`

	SentryDSN       string `yaml:"sentry_dsn"`
	JaegerCollector string `yaml:"jaeger_collector"`

	Profiling bool `yaml:"-"`
}

func NewConfig(yamlFile string) (*Config, error) {
	conf := &Config{}
	err := loadYaml(yamlFile, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

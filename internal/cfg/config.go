package cfg

type Config struct {
	HostAddress    string `json:"host_address" yaml:"host_address"`
	MetricsAddress string `json:"metrics_address" yaml:"metrics_address"`

	SentryDSN       string `json:"sentry_dsn" yaml:"sentry_dsn"`
	JaegerCollector string `json:"jaeger_collector" yaml:"jaeger_collector"`
}

func NewConfig(yamlFile string) (*Config, error) {
	conf := &Config{}
	err := loadYaml(yamlFile, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

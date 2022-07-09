package cfg

type Config struct {
	JWTSecret string `yaml:"jwt_secret"`
	// JWTTTL stored in seconds.
	JWTTTL int `yaml:"jwt_ttl"`

	HostAddress    string `yaml:"host_address"`
	MetricsAddress string `yaml:"metrics_address"`
	GrpcAddress    string `yaml:"grpc_address"`

	SentryDSN       string `yaml:"sentry_dsn"`
	JaegerCollector string `yaml:"jaeger_collector"`

	Profiling bool `yaml:"-"`

	DB struct {
		Login    string `yaml:"login"`
		Password string `yaml:"pass"`
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}

func (c *Config) GetJWTSecret() string {
	return c.JWTSecret
}
func (c *Config) GetJWTTTL() int {
	return c.JWTTTL
}
func (c *Config) GetProfiling() bool {
	return c.Profiling
}

func NewConfig(yamlFile string) (*Config, error) {
	conf := &Config{}
	err := loadYaml(yamlFile, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

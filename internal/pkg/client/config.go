package client

type FSClientConfig struct {
	Location string
}

type HTTPClientConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

type Config struct {
	Type string
	Log  bool
	FS   *FSClientConfig   `mapstructure:"fs"`
	HTTP *HTTPClientConfig `mapstructure:"http"`
}

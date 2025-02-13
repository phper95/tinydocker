package main

type HttpClientConfig struct {
	Timeout int
	Retry   int
}
type Option func(*HttpClientConfig)

func WithTimeout(timeout int) Option {
	return func(config *HttpClientConfig) {
		config.Timeout = timeout
	}
}

func WithRetry(retry int) Option {
	return func(config *HttpClientConfig) {
		config.Retry = retry
	}
}

func NewHttpClient(opts ...Option) {
	httpConfig := HttpClientConfig{}
	for _, opt := range opts {
		opt(&httpConfig)
	}
}

func main() {
	NewHttpClient(WithTimeout(10), WithRetry(3))

}

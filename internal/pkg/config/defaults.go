package config


// DefaultLogLevel defines default log level for server
const DefaultLogLevel = "info"

// DefaultServerAddr is the default address for http server
const DefaultServerAddr = "0.0.0.0"

// DefaultServerPort is the default port for http server
const DefaultServerPort = "8080"

// DefaultServerGracefulTimeout is the timeout in seconds for server to wait all connections to
// be finished
const DefaultServerGracefulTimeout = 5

// DefaultProcessStrategy is the default strategy to process image from URL parameters/path
const DefaultProcessStrategy = "query"

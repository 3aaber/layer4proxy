package env

var GeoipDatabase string

/**
 * Config file top-level object
 */
type Config struct {
	Title            string               `json:"title"` // Title of the application
	Name             string               `toml:"name"`  // Dealer Name
	Host             string               `toml:"host"`  // Dealer Host IP Address
	GeoipDatabase    string               `toml:"geoip_database"`
	Logging          *LoggingConfig       `toml:"logging"`
	Metrics          *MetricsConfig       `toml:"metrics"`
	Redis            *RedisDatabaseConfig `toml:"redis_server"` // Redis Config to recieve log channels config in runtime
	Defaults         *ConnectionOptions   `toml:"defaults"`
	Profiler         *ProfilerConfig      `toml:"profiler"`
	PrometheusConfig *PrometheusConfig    `toml:"prometheus"`
}

/**
 * Logging config section
 */
type LoggingConfig struct {
	Level  int    `toml:"level"`
	Output string `toml:"output"`
	Format string `toml:"format"`
	Sentry string `toml:"sentry"`
}

/**
 * Metrics config section
 */
type MetricsConfig struct {
	Enabled bool   `toml:"enabled"`
	Bind    string `toml:"bind"`
}

/**
 * Default values can be overridden in server
 */
type ConnectionOptions struct {
	MaxConnections           *int    `toml:"max_connections" json:"max_connections"`
	ClientIdleTimeout        *string `toml:"client_idle_timeout" json:"client_idle_timeout"`
	BackendIdleTimeout       *string `toml:"backend_idle_timeout" json:"backend_idle_timeout"`
	BackendConnectionTimeout *string `toml:"backend_connection_timeout" json:"backend_connection_timeout"`
}

/**
 * Pprof profiler config
 */
type ProfilerConfig struct {
	// enabled
	Enabled bool `toml:"enabled"`
	// hostname:port
	Bind string `toml:"bind"`
}

// PrometheusConfig Prometheus Config to export metrics
type PrometheusConfig struct {
	Address string `toml:"address"`
	Port    string `toml:"port"`
	Path    string `toml:"path"`
}

// RedisDatabaseConfig Redis DatabaseConfigurations exported
type RedisDatabaseConfig struct {
	Addr     string `toml:"addr"`
	Password string `toml:"pass"`
	Channel  string `toml:"channel"`
	Number   int    `toml:"db"`
	Prefix   string `toml:"prefix"`
}

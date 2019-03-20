package kvstore

import (
	"fmt"
)

// Config represents configuration of key/value store
type Config struct {
	Type         string
	Prefix       string
	Redis        RedisConfig        `mapstructure:"redis"`
	RedisCluster RedisClusterConfig `mapstructure:"redis-cluster"`
	Cache        CacheConfig        `mapstructure:"cache"`
}

// RedisConfig represents configuration for redis k/v storage
type RedisConfig struct {
	Host       string
	Port       int
	Password   string
	DB         int
	Expiration int
}

// Addr returns string representing address:port of redis server
func (r RedisConfig) Addr() string {
	return fmt.Sprint(r.Host, ":", r.Port)
}

// RedisClusterConfig represents configuration struct for redis k/v store in cluster mode
type RedisClusterConfig struct {
	Expiration int
	Password   string
	Addresses  []string
}

// CacheConfig represents configuration for in-memory kv storage
type CacheConfig struct {
	Expiration      int
	CleanupInterval int `mapstructure:"cleanup_interval"`
}

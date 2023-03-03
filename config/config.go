package config

import (
	"os"
	"strconv"
)

type AquiductConfiguration struct {
	ClusterAdvertiseAddr string
	ClusterAdvertisePort int
	ClusterBindAddr      string
	ClusterBindPort      int
}

var Configuration *AquiductConfiguration

func Load() {
	Configuration = &AquiductConfiguration{
		ClusterAdvertiseAddr: getEnv("CLUSTER_ADVERTISE_ADDR", ""),
		ClusterAdvertisePort: getInt(getEnv("CLUSTER_ADVERTISE_PORT", "")),
		ClusterBindAddr:      getEnv("CLUSTER_BIND_ADDR", ""),
		ClusterBindPort:      getInt(getEnv("CLUSTER_BIND_PORT", "")),
	}
}

func getInt(val string) int {
	v, _ := strconv.Atoi(val)
	return v
}
func getEnv(key string, default_val string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return default_val
}

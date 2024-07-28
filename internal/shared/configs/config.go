package configs

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Server struct {
	HttpPort int `yaml:"httpPort"`
	GrpcPort int `yaml:"grpcPort"`
}

// Config represents the configuration for the calculator server.
type Config struct {
	Server               Server `yaml:"server"`
	OrchestratorURL      string `yaml:"orchestratorURL"`
	ComputingPower       int    `yaml:"computingPower"`
	TimeAdditionMS       int    `yaml:"timeAdditionMS"`
	TimeSubtractionMS    int    `yaml:"timeSubtractionMS"`
	TimeMultiplicationMS int    `yaml:"timeMultiplicationMS"`
	TimeDivisionMS       int    `yaml:"timeDivisionMS"`
}

// LoadConfig loads the configuration from a YAML file.
func LoadConfig(path string) (*Config, error) {
	defaultConfig := &Config{
		Server:               Server{HttpPort: 8080, GrpcPort: 8081},
		OrchestratorURL:      "localhost:8081",
		ComputingPower:       4,
		TimeAdditionMS:       100,
		TimeSubtractionMS:    200,
		TimeMultiplicationMS: 300,
		TimeDivisionMS:       400,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return defaultConfig, nil
	}

	config, err := ConfigFromData(data)
	if err != nil {
		return defaultConfig, nil
	}

	ConfigFromEnvironment(config)

	return config, nil
}

// ConfigFromEnvironment loads the configuration from environment variables.
func ConfigFromEnvironment(cfg *Config) {

	cfg.TimeAdditionMS = getEnvAsInt("TIME_ADDITION_MS", cfg.TimeAdditionMS)
	cfg.TimeSubtractionMS = getEnvAsInt("TIME_SUBTRACTION_MS", cfg.TimeSubtractionMS)
	cfg.TimeMultiplicationMS = getEnvAsInt("TIME_MULTIPLICATIONS_MS", cfg.TimeMultiplicationMS)
	cfg.TimeDivisionMS = getEnvAsInt("TIME_DIVISIONS_MS", cfg.TimeDivisionMS)
	cfg.ComputingPower = getEnvAsInt("COMPUTING_POWER", cfg.ComputingPower)
	cfg.OrchestratorURL = getEnvAsString("ORCHESTRATOR_URL", cfg.OrchestratorURL)
	cfg.Server.HttpPort = getEnvAsInt("SERVER_PORT", cfg.Server.HttpPort)
}

// ConfigFromData loads the configuration from a YAML byte array.
func ConfigFromData(data []byte) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnvAsInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvAsString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

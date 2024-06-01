package configs

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestConfigFromEnvironment(t *testing.T) {
	// Test case 1: TimeAdditionMS environment variable is set
	t.Run("TimeAdditionMS environment variable is set", func(t *testing.T) {
		os.Setenv("TIME_ADDITION_MS", "100")
		cfg := &Config{}
		ConfigFromEnvironment(cfg)
		if cfg.TimeAdditionMS != 100 {
			t.Errorf("Expected TimeAdditionMS to be 100, got %d", cfg.TimeAdditionMS)
		}
		os.Unsetenv("TIME_ADDITION_MS")
	})

	// Test case 2: TimeSubtractionMS environment variable is set
	t.Run("TimeSubtractionMS environment variable is set", func(t *testing.T) {
		os.Setenv("TIME_SUBTRACTION_MS", "200")
		cfg := &Config{}
		ConfigFromEnvironment(cfg)
		if cfg.TimeSubtractionMS != 200 {
			t.Errorf("Expected TimeSubtractionMS to be 200, got %d", cfg.TimeSubtractionMS)
		}
		os.Unsetenv("TIME_SUBTRACTION_MS")
	})

	// Test case 3: TimeMultiplicationMS environment variable is set
	t.Run("TimeMultiplicationMS environment variable is set", func(t *testing.T) {
		os.Setenv("TIME_MULTIPLICATIONS_MS", "300")
		cfg := &Config{}
		ConfigFromEnvironment(cfg)
		if cfg.TimeMultiplicationMS != 300 {
			t.Errorf("Expected TimeMultiplicationMS to be 300, got %d", cfg.TimeMultiplicationMS)
		}
		os.Unsetenv("TIME_MULTIPLICATIONS_MS")
	})

	// Test case 4: TimeDivisionMS environment variable is set
	t.Run("TimeDivisionMS environment variable is set", func(t *testing.T) {
		os.Setenv("TIME_DIVISIONS_MS", "400")
		cfg := &Config{}
		ConfigFromEnvironment(cfg)
		if cfg.TimeDivisionMS != 400 {
			t.Errorf("Expected TimeDivisionMS to be 400, got %d", cfg.TimeDivisionMS)
		}
		os.Unsetenv("TIME_DIVISIONS_MS")
	})

	// Test case 5: ComputingPower environment variable is set
	t.Run("ComputingPower environment variable is set", func(t *testing.T) {
		os.Setenv("COMPUTING_POWER", "5")
		cfg := &Config{}
		ConfigFromEnvironment(cfg)
		if cfg.ComputingPower != 5 {
			t.Errorf("Expected ComputingPower to be 5, got %d", cfg.ComputingPower)
		}
		os.Unsetenv("COMPUTING_POWER")
	})

	// Test case 6: OrchestratorURL environment variable is set
	t.Run("OrchestratorURL environment variable is set", func(t *testing.T) {
		os.Setenv("ORCHESTRATOR_URL", "http://example.com")
		cfg := &Config{}
		ConfigFromEnvironment(cfg)
		if cfg.OrchestratorURL != "http://example.com" {
			t.Errorf("Expected OrchestratorURL to be http://example.com, got %s", cfg.OrchestratorURL)
		}
		os.Unsetenv("ORCHESTRATOR_URL")
	})

}

func TestConfigFromData(t *testing.T) {
	// Test case: Valid YAML
	validConfig := &Config{
		Server:               Server{Port: 8080},
		OrchestratorURL:      "http://localhost:8080",
		ComputingPower:       4,
		TimeAdditionMS:       100,
		TimeSubtractionMS:    200,
		TimeMultiplicationMS: 300,
		TimeDivisionMS:       400,
	}
	data, err := yaml.Marshal(validConfig)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}
	config, err := ConfigFromData(data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(config, validConfig) {
		t.Errorf("Expected config %v, got %v", validConfig, config)
	}

	// Test case: Invalid YAML
	_, err = ConfigFromData([]byte("invalid yaml"))
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestGetEnvAsInt(t *testing.T) {
	// Test case 1: Environment variable not set
	key := "TEST_KEY"
	defaultVal := 10
	os.Unsetenv(key)
	result := getEnvAsInt(key, defaultVal)
	if result != defaultVal {
		t.Errorf("Expected %d, got %d", defaultVal, result)
	}

	// Test case 2: Environment variable set to valid integer
	key = "TEST_KEY"
	expected := 42
	os.Setenv(key, "42")
	result = getEnvAsInt(key, defaultVal)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}

	// Test case 3: Environment variable set to non-integer value
	key = "TEST_KEY"
	os.Setenv(key, "not a number")
	result = getEnvAsInt(key, defaultVal)
	if result != defaultVal {
		t.Errorf("Expected %d, got %d", defaultVal, result)
	}
}

func TestGetEnvAsString(t *testing.T) {
	// Test case: environment variable exists
	os.Setenv("TEST_KEY", "test_value")
	result := getEnvAsString("TEST_KEY", "default_value")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test case: environment variable does not exist
	result = getEnvAsString("NON_EXISTENT_KEY", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}

	// Test case: empty default value
	result = getEnvAsString("NON_EXISTENT_KEY", "")
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

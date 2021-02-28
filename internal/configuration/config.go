package configuration

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Environment string
	Server      struct {
		Binance struct {
			WebSocket struct {
				URL          string `yaml:"url"`
				LimitRequest int    `yaml:"limitRequest"`
			} `yaml:"websocket"`
			Restful struct {
				URL string `yaml:"url"`
			} `yaml:"restful"`
		} `yaml:"binance"`
	} `yaml:"server"`
}

func NewConfig(env string, configPath string) (*Config, error) {
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	config.Environment = env
	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func ParseFlags() (string, string, error) {
	var env string
	flag.StringVar(&env, "env", "local", "path to config file")

	configPath := fmt.Sprintf("./configurations/%s.yaml", env)

	flag.Parse()

	if err := ValidateConfigPath(configPath); err != nil {
		return "", "", err
	}

	return env, configPath, nil
}

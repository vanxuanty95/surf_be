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
		PfitMgmt struct {
			APIPort string `yaml:"api_port"`
		} `yaml:"pfit_mgmt"`
		Binance struct {
			WebSocket struct {
				WSURL        string `yaml:"ws_url"`
				StreamURL    string `yaml:"stream_url"`
				LimitRequest int    `yaml:"limitRequest"`
			} `yaml:"websocket"`
			Restful struct {
				URL string `yaml:"url"`
			} `yaml:"restful"`
		} `yaml:"binance"`
		DataBase struct {
			Redis struct {
				Host     string `yaml:"host"`
				Port     int    `yaml:"port"`
				Password string `yaml:"password"`
			} `yaml:"redis"`
			Mongo struct {
				Host         []string `yaml:"host"`
				Port         int      `yaml:"port"`
				AuthDatabase string   `yaml:"auth_database"`
				Username     string   `yaml:"username"`
				Password     string   `yaml:"password"`
				Database     string   `yaml:"database"`
				Collection   struct {
					Bot         string `yaml:"bot"`
					Transaction string `yaml:"transaction"`
				} `yaml:"collection"`
			} `yaml:"mongo"`
		} `yaml:"database"`
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

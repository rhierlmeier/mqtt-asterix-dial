package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Broker   string `yaml:"broker"`
	ClientId string `yaml:"client_id"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	Paths Paths `yaml:"paths"`

	Calls []CallTemplate `yaml:"calls"`
}

type Paths struct {
	CallFileDir string `yaml:"call_file_dir"`
	TempDir     string `yaml:"tmp_dir"`
}

type CallTemplate struct {
	Name             string         `yaml:"name"`
	Topic            string         `yaml:"topic"`
	Value            string         `yaml:"value"`
	CallFileTemplate string         `yaml:"template"`
	Variables        []CallVariable `yaml:"variables"`
}

type CallVariable struct {
	Topic string `yaml:"topic"`
	Name  string `yaml:"name"`
}

func (c *Config) LoadFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(c)
}

func (c *Config) LoadFromEnv() {
	c.Broker = os.Getenv("MQTT_BROKER")
	c.ClientId = os.Getenv("MQTT_CLIENT_ID")
	c.Username = os.Getenv("MQTT_USERNAME")
	c.Password = os.Getenv("MQTT_PASSWORD")
}

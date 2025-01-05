package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Broker   string `yaml:"broker"`
	ClientId string `yaml:"client_id"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	CallFileDir string `yaml:"call_file_dir"`

	Calls []CallTemplate `yaml:"calls"`
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
	ret := decoder.Decode(c)

	if err := c.Validate(); err != nil {
		return err
	}

	return ret
}

func checkDirExists(dir string) error {
	fileInfo, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", dir)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("path %s is not a directory", dir)
	}
	return nil
}

func (ct *CallTemplate) Validate() error {
	if ct.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if ct.Topic == "" {
		return fmt.Errorf("topic cannot be empty")
	}

	if ct.CallFileTemplate == "" {
		return fmt.Errorf("template cannot be empty")
	}

	for varIndex, variable := range ct.Variables {
		if variable.Name == "" {
			return fmt.Errorf("variables[%d].name cannot be empty", varIndex)
		}
		if variable.Topic == "" {
			return fmt.Errorf("variables[%d].topic cannot be empty", varIndex)
		}
	}
	return nil
}

func (c *Config) Validate() error {
	if c.Broker == "" {
		return fmt.Errorf("broker cannot be empty")
	}

	if c.ClientId == "" {
		c.ClientId = "mqtt-dial"
	}

	if c.CallFileDir == "" {
		return fmt.Errorf("call_file_dir cannot be empty")
	}

	if err := checkDirExists(c.CallFileDir); err != nil {
		return fmt.Errorf("invalid call_file_dir: %v", err)
	}

	if len(c.Calls) == 0 {
		return fmt.Errorf("calls cannot be empty")
	}

	for callIndex, call := range c.Calls {
		if err := call.Validate(); err != nil {
			return fmt.Errorf("calls[%d]: %v", callIndex, err)
		}
	}

	return nil
}

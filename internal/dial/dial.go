package dial

import (
	"bytes"
	"fmt"
	"log"
	"mqtt-dial/internal/config"
	"os"
	"text/template"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Dialer struct {
	mqttClient   mqtt.Client
	paths        config.Paths
	callTemplate config.CallTemplate

	subscribeToken mqtt.Token

	variableValues map[string]interface{}
}

func NewDialer(mqttClient mqtt.Client, paths config.Paths, callTemplate config.CallTemplate) (*Dialer, error) {

	if mqttClient == nil {
		return nil, fmt.Errorf("mqttClient cannot be nil")
	}
	if callTemplate.Topic == "" {
		return nil, fmt.Errorf("callTemplate.Topic cannot be empty")
	}

	return &Dialer{
		mqttClient:     mqttClient,
		paths:          paths,
		callTemplate:   callTemplate,
		variableValues: make(map[string]interface{}),
	}, nil
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

func (d *Dialer) Start() error {

	if d.callTemplate.Topic == "" {
		return fmt.Errorf("callTemplate.Topic cannot be empty")
	}

	if d.paths.CallFileDir == "" {
		return fmt.Errorf("paths.callFileDir cannot be empty")
	}
	if err := checkDirExists(d.paths.CallFileDir); err != nil {
		return err
	}

	if d.paths.TempDir == "" {
		return fmt.Errorf("paths.tempDir cannot be empty")
	}
	if err := checkDirExists(d.paths.TempDir); err != nil {
		return err
	}

	for index, variable := range d.callTemplate.Variables {
		if variable.Name == "" {
			return fmt.Errorf("variables[%d].name cannot be empty", index)
		}
		if variable.Topic == "" {
			return fmt.Errorf("variables[%d].value cannot be empty", index)
		}
	}

	for _, variable := range d.callTemplate.Variables {
		d.mqttClient.Subscribe(variable.Topic, 0, func(client mqtt.Client, msg mqtt.Message) {
			d.onVariableChanged(variable.Name, string(msg.Payload()))
		})
		log.Printf("Call %s: Subscribed to topic %s for variable %s", d.callTemplate.Name, variable.Topic, variable.Name)
	}

	d.subscribeToken = d.mqttClient.Subscribe(d.callTemplate.Topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		d.onValueChanged(string(msg.Payload()))
	})
	log.Printf("Call %s: Subscribed to topic %s", d.callTemplate.Name, d.callTemplate.Topic)

	return nil
}

func (d *Dialer) onVariableChanged(name string, value string) {
	log.Printf("Call %s: Variable %s received value %s", d.callTemplate.Name, name, value)
	d.variableValues[name] = value
}

func (d *Dialer) onValueChanged(mqttValue string) {

	if d.callTemplate.Value == mqttValue {

		log.Printf("Call %s: Value %s received, creating call file", d.callTemplate.Name, mqttValue)

		tmpl, err := template.New("callfile").Parse(d.callTemplate.CallFileTemplate)
		if err != nil {
			log.Printf("Could not parse call template: %v", err)
			return
		}

		var callFileContent bytes.Buffer
		err = tmpl.Execute(&callFileContent, d.variableValues)
		if err != nil {
			log.Printf("Could not execute call template: %v", err)
			return
		}
		tempFile, err := os.CreateTemp(d.paths.CallFileDir, "callfile-*.call")
		if err != nil {
			log.Printf("Could not create temp file %s: %v", tempFile.Name(), err)
			return
		}
		defer tempFile.Close()

		os.Chmod(tempFile.Name(), 0644)

		_, err = tempFile.Write(callFileContent.Bytes())
		if err != nil {
			log.Printf("Error writing call file %s: %v", tempFile.Name(), err)
			return
		}
		tempFile.Close()
	}

}

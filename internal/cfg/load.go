package cfg

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func loadYaml(yamlFile string, yamlCfg interface{}) error {
	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, yamlCfg)
}

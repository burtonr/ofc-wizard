package actions

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/burtonr/ofc-wizard/types"
	"gopkg.in/yaml.v2"
)

// CreateInitFile marshalls the input answers into a yml file to be used with ofc-bootstrap
func CreateInitFile() *types.InitYaml {
	fmt.Println("Creating file")
	init := &types.InitYaml{}

	if _, err := os.Stat("init.yml"); os.IsNotExist(err) {
		if fErr := ioutil.WriteFile("init.yml", nil, 0644); fErr != nil {
			fmt.Printf("Trouble creating new init.yml file: %s\n", fErr.Error())
			panic(fErr)
		}
	} else {
		fmt.Println("Discovered existing init.yml. Loading existing values.")
		init = LoadInitFile()
	}

	return init
}

// LoadInitFile marshalls the values from the init.yml file in the local directory
func LoadInitFile() *types.InitYaml {
	fmt.Println("Loading existing init.yml file")
	yamlBytes, yamlErr := ioutil.ReadFile("init.yml")
	if yamlErr != nil {
		fmt.Fprintf(os.Stderr, "-yaml file gave error: %s\n", yamlErr.Error())
		os.Exit(1)
	}

	init := types.InitYaml{}
	unmarshalErr := yaml.Unmarshal(yamlBytes, &init)
	if unmarshalErr != nil {
		fmt.Fprintf(os.Stderr, "-yaml file gave error: %s\n", unmarshalErr.Error())
		os.Exit(1)
	}

	return &init
}

// WriteInitFile writes the values to the init.yml file in the local directory
func WriteInitFile(yml types.InitYaml) {
	fmt.Println("Writing the file")
}

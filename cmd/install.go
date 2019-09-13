/*
Copyright © 2019 Burton Rheutan

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/burtonr/ofc-wizard/actions"
	"github.com/spf13/cobra"
)

// generateCmd represents the install command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Creates a new OpenFaaS Cloud init.yml file to be used with ofc-bootstrap",
	Long: `This will step through all of the configuration required to configure
the init.yml file used by the OpenFaaS Cloud bootstrap tool.

The wizard will ask relevant questions and adjust as you enter your values
to ensure that your new OpenFaaS Cloud installation will be successful!`,
	Run: func(cmd *cobra.Command, args []string) {
		actions.GenerateYaml()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"fmt"

	"github.com/spf13/cobra"
)

const initConfigData := `
{
    "type":"artifact/component",
    "replacements":[
        {
            "placeholder":"VAR__TO_BE_REPLACED",
            "desc":"This is description that will be prompted on artifact/component initialization"
        }
    ],
    "injections":[
        {
            "file_name":"File or directory's name (relative path)",
            "injection_point":"Point where this file should be injected (relative path)"
        }
    ],
    "commands":[
        "These will run in newly created project's directory",
        "Add initialization commands here ...",
        "npm install, go get etc..."
    ]
}
`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Adds a fire-up.json file to current working directory.",
	Long: `Adds a fire-up.json file to current working directory.
	Adding a fire-up.json to a project makes it identifiable as an artifact template by fire-up.
	`,
This application is a tool to generate the needed files
	Run: func(cmd *cobra.Command, args []string) {
		
		fmt.Println("Please configure fire-up.json that was added to this directory.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

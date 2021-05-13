package cmd

import (
	"fire-up/utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const fireUp = "fire-up.json"
const initConfigData = `
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
	Run: func(cmd *cobra.Command, args []string) {
		CreateFireUpJsonIfNotPresent()
	},
}

func CreateFireUpJsonIfNotPresent() {
	if _, err := os.Stat(fireUp); os.IsNotExist(err) {
		_, err := os.Create(fireUp)
		if err != nil {
			log.Fatal("Error while creating fire-up.json", err)
		}
		utils.Must(ioutil.WriteFile(fireUp,[]byte(initConfigData),os.ModePerm),
			"Error writing to fire-up.config")
		fmt.Println("Please configure fire-up.json that was added to this directory.")
	}
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

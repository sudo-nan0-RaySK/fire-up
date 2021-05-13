package cmd

import (
	"fire-up/utils"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

const (
	// Error message for wrong usage of this sub-command
	addLocalErrStr = `No src directory provided!
Usage: fire-up add-local <src-directory> --alias <artifact-alias>`
	aliasUsage = `Usage: fire-up add-local <src-directory> --alias <artifact-alias>`
)

// Value of artifactAliasFlag
var artifactAliasFlag = ""


// addLocalCmd represents the addLocal command
var addLocalCmd = &cobra.Command{
	Use:   "add-local",
	Short: "Add a project artifact to local machine",
	Long:  `Add a project artifact to local machine`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add-local called")
		if len(args) <= 0 {
			log.Fatal(addLocalErrStr)
		}
		if artifactAliasFlag==""{
			log.Fatal("Giving an alias (--alias <alias-name>) is a must for adding artifacts or components")
		}
		artifactPath := args[0]
		utils.AddArtifact(artifactPath, artifactAliasFlag)
	},
}

func init() {
	rootCmd.AddCommand(addLocalCmd)
	addLocalCmd.PersistentFlags().StringVar(&artifactAliasFlag, "alias", "", aliasUsage)
}


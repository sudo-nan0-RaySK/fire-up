package cmd

import (
	"fire-up/utils"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var aliasName string

// addRemoteCmd represents the addRemote command
var addRemoteCmd = &cobra.Command{
	Use:   "add-remote",
	Short: "Add a project artifact to github remote repository",
	Long:  `Add a project artifact to github remoter repository `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addRemote called")
		if len(args) != 1 {
			log.Fatal("Wrong usage!, use like: fire-up add-remote <path-to-artifact> --alias <artifact-alias>")
		}
		utils.CheckAuthenticated()
		artifactsRepo := utils.GetOrCreateArtifactRepositoryIfNotPresent()
		utils.CreateResource(args[0], aliasName, *artifactsRepo.ContentsURL)
	},
}

func init() {
	rootCmd.AddCommand(addRemoteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	addRemoteCmd.PersistentFlags().StringVar(&aliasName, "alias", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addRemoteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

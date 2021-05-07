package cmd

import (
	"encoding/json"
	"fire-up/types"
	"fire-up/utils"
	"fmt"
	copy2 "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	// Error message for wrong usage of this sub-command
	addLocalErrStr = `No src directory provided!
Usage: fire-up add-local <src-directory> --alias <artifact-alias>`
	aliasUsage = `Usage: fire-up add-local <src-directory> --alias <artifact-alias>`
)

// Value of artifactAliasFlag
var artifactAliasFlag string


// addLocalCmd represents the addLocal command
var addLocalCmd = &cobra.Command{
	Use:   "add-local",
	Short: "Add a project artifact from local machine",
	Long:  `Add a project artifact from local machine`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add-local called")
		if len(args) <= 0 {
			log.Fatal(addLocalErrStr)
		}
		artifactPath := args[0]
		AddArtifact(artifactPath, artifactAliasFlag)
	},
}

func AddArtifact(artifactPath string, artifactAlias string) {
	CreateDirectoryIfNotPresent()
	CreateArtifactsFileIfNotPresent()
	// Getting artifact's directory name
	artifactDirName := strings.Split(artifactPath, "/")
	constructedPath := utils.ConfigDir + artifactDirName[len(artifactDirName)-1]

	// Set artifact alias to be equivalent to artifactDirName if not specified
	if artifactAlias == "" {
		artifactAlias = artifactDirName[len(artifactDirName)-1]
	}
	log.Printf("constructed path is :- %s", constructedPath)
	err := copy2.Copy(artifactPath, constructedPath)
	if err != nil {
		log.Fatal("Error while copying artifact!", err)
	}
	newArtifact := types.ArtifactEntry{ArtifactAlias: artifactAlias, ArtifactPath: constructedPath}
	// Check if artifact.json is empty
	file, _ := os.Stat(utils.ArtifactsFile)
	// If file is empty
	if file.Size() == 0 {
		artifactEntries := make(types.Artifacts, 0)
		artifactEntries = append(artifactEntries, newArtifact)
		log.Println(artifactEntries)
		writeToArtifactsFile(artifactEntries)
	} else { // If file contains previous entries
		artifactRecords, _ := ioutil.ReadFile(utils.ArtifactsFile)
		var artifacts types.Artifacts
		utils.Must(json.Unmarshal(artifactRecords, &artifacts),
			"Error while un-marshalling artifacts")
		artifacts.CheckDuplicateRecords(newArtifact)
		writeToArtifactsFile(append(artifacts, newArtifact))
	}
}

func writeToArtifactsFile(artifactEntries types.Artifacts) {
	artifactRecords, marshalErr := json.MarshalIndent(artifactEntries,"", "    ")
	if marshalErr != nil {
		log.Fatal("Error while marshalling artifacts")
	}
	log.Println("artifact records :-" + string(artifactRecords))
	writeErr := ioutil.WriteFile(utils.ArtifactsFile, artifactRecords, 0777)
	if writeErr != nil {
		log.Fatal("Error while marshalling artifacts")
	}
}

func CreateArtifactsFileIfNotPresent() {
	if _, err := os.Stat(utils.ArtifactsFile); os.IsNotExist(err) {
		_, err := os.Create(utils.ArtifactsFile)
		if err != nil {
			log.Fatal("Error while creating config directory", err)
		}
	}
}

func CreateDirectoryIfNotPresent() {
	if _, err := os.Stat(utils.ConfigDir); os.IsNotExist(err) {
		err := os.Mkdir(utils.ConfigDir, 0777)
		if err != nil {
			log.Fatal("Error while creating config directory", err)
		}
	}
}

func init() {
	rootCmd.AddCommand(addLocalCmd)
	addLocalCmd.PersistentFlags().StringVar(&artifactAliasFlag, "alias", "", aliasUsage)
}


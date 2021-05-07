package utils

import (
	"bufio"
	"encoding/json"
	"fire-up/types"
	"fmt"
	copy2 "github.com/otiai10/copy"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ConfigDir Location of the config directory
var ConfigDir = os.Getenv("HOME") + "/.fire-up/"

// ArtifactsFile Location of artifacts.json file
var ArtifactsFile = ConfigDir + "artifacts.json"

var FireUpConfig = "/fire-up.json"

func InitializeProjectFromArtifact(artifactName string) {
	artifactsData, err := ioutil.ReadFile(ArtifactsFile)
	Must(err, "Error while reading artifacts.json")
	var artifactRecords types.Artifacts
	Must(json.Unmarshal(artifactsData, &artifactRecords),
		"Error while unmarshalling artifacts.json")
	for _, artifact := range artifactRecords {
		if artifact.ArtifactAlias == artifactName {
			createProject(artifactName, artifact.ArtifactPath)
			return
		}
	}
	// TODO: Fetch from github
	log.Fatal("No such artifact found!")
}

func createProject(artifactName string, artifactPath string) {
	pathNodes := strings.Split(artifactPath,"/")
	nodeName := pathNodes[len(pathNodes)-1]
	Must(copy2.Copy(artifactPath, nodeName, copy2.Options{AddPermission: os.ModePerm}),
		"Error copying artifact")
	configDataRaw, err := ioutil.ReadFile(artifactPath + FireUpConfig)
	Must(err, "Error occurred while reading fire-up.json")
	var configData types.Config
	Must(json.Unmarshal(configDataRaw, &configData),
		"Error while unmarshalling fire-up.json")
	replacementMap := configData.ReplacementList.ConstructReplacementMap()
	var fileObjects = make([]string, 0)
	log.Printf("replacementMap %v", replacementMap)
	Must(filepath.WalkDir(nodeName, getWalkAndCollect(&fileObjects)), "Error initiating WalkDir()")
	for index := range fileObjects{
		replaceFileOrDir(fileObjects[len(fileObjects)-1-index],replacementMap)
	}
}

func getWalkAndCollect(fileObjects *[]string) fs.WalkDirFunc {
	return func(name string, dir fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal("An error occurred while walking project directory recursively", err)
		}
		*fileObjects = append(*fileObjects, name)
		return nil
	}
}

func getReplacement(original string, replacementMap map[string]string) (string, bool) {
	var replacedPath = original
	for replacementTarget,replacement := range replacementMap{
		replacedPath = strings.Replace(replacedPath, replacementTarget, replacement, -1)
	}
	if replacedPath!=original{
		return replacedPath,true
	}
	return "",false
}

func replaceFileOrDir(fileOrDir string, replacementMap map[string]string) {
	pathElements := strings.Split(fileOrDir,"/")
	nodeName := pathElements[len(pathElements)-1]
	newName, doReplace := getReplacement(nodeName, replacementMap)
	if doReplace {
		pathElements[len(pathElements)-1] = newName
		newPath := strings.Join(pathElements,"/")
		fmt.Printf("%s -> %s \n", fileOrDir, newPath)
		Must(os.Rename(fileOrDir, newPath), "Error renaming file/dir!")
		replaceFileContents(newPath, replacementMap)
	} else {
		replaceFileContents(fileOrDir, replacementMap)
	}
}

func replaceFileContents(fileName string, replacementMap map[string]string){
	info, err := os.Stat(fileName)
	Must(err, "Error opening file/directory")
	if !info.IsDir() {
		fd, err := os.Open(fileName)
		Must(err, "Error opening file/directory")
		scanner := bufio.NewScanner(fd)
		scanner.Split(bufio.ScanLines)
		var buffer = ""
		for scanner.Scan() {
			nextLine := scanner.Text()
			replacedLine, wasReplaced := getReplacement(nextLine, replacementMap)
			if wasReplaced{
				buffer += replacedLine + "\n"
			} else {
				buffer += nextLine + "\n"
			}
		}
		Must(ioutil.WriteFile(fileName, []byte(buffer), fs.ModePerm),
			"Error replacing file contents")
		defer Must(fd.Close(), "Error closing file/directory")
	}
}

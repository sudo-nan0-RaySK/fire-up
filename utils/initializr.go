package utils

import (
	"bufio"
	"encoding/json"
	"fire-up/types"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lithammer/shortuuid"
	copy2 "github.com/otiai10/copy"
)

// ConfigDir Location of the config directory
var ConfigDir = os.Getenv("HOME") + "/.fire-up/"

// ArtifactsFile Location of artifacts.json file
var ArtifactsFile = ConfigDir + "artifacts.json"

var FireUpConfig = "/fire-up.json"

func InitializeProjectFromArtifact(artifactName string, newName string) {
	// Searching in local
	artifactsData, err := ioutil.ReadFile(ArtifactsFile)
	Must(err, "Error while reading artifacts.json")
	var artifactRecords types.Artifacts
	Must(json.Unmarshal(artifactsData, &artifactRecords),
		"Error while unmarshalling artifacts.json")
	for _, artifact := range artifactRecords {
		if artifact.ArtifactAlias == artifactName {
			createProject(newName, artifact.ArtifactPath)
			return
		}
	}
	// Searching in private remote
	CheckAuthenticated()
	artifactsRepo := GetOrCreateArtifactRepositoryIfNotPresent()
	artifactPath, err := GetResource(artifactName, *artifactsRepo.ContentsURL)
	if err!=nil{
		log.Fatal(err)
	}
	if artifactPath!=""{
		AddArtifact(artifactPath, artifactName)
		createProject(newName, artifactPath)
		Must(os.RemoveAll(artifactPath), "Error removing tmp data")
		return
	}

	// TODO: Search in global remote
	log.Fatal("No such artifact found!")
}

func createProject(artifactName string, artifactPath string) {
	//pathNodes := strings.Split(artifactPath,"/")
	nodeName := artifactName
	cpyOpts := copy2.Options{
		AddPermission: os.ModePerm,
		Skip: func(src string) (bool, error) {
			pathParts := strings.Split(src, "/")
			file  := pathParts[len(pathParts)-1]
			if file == "fire-up.json" {
				return true,nil
			}
			return false,nil
		},
	}
	Must(copy2.Copy(artifactPath, nodeName, cpyOpts),
		"Error copying artifact")
	configDataRaw, err := ioutil.ReadFile(artifactPath + FireUpConfig)
	Must(err, "Error occurred while reading fire-up.json")
	var configData types.Config
	Must(json.Unmarshal(configDataRaw, &configData),
		"Error while unmarshalling fire-up.json")
	replacementMap := configData.ReplacementList.ConstructReplacementMap()
	var fileObjects = make([]string, 0)
	log.Printf("replacementMap %v", replacementMap)
	Must(filepath.WalkDir(nodeName, GetWalkAndCollect(&fileObjects)), "Error initiating WalkDir()")
	for index := range fileObjects{
		replaceFileOrDir(fileObjects[len(fileObjects)-1-index],replacementMap)
	}
	// Running all the initialization commands
	currDir,err:= os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if configData.CommandList!=nil && len(configData.CommandList)>0{
		RunCommands(configData.CommandList, currDir+"/"+nodeName)
	}
}

func GetWalkAndCollect(fileObjects *[]string) fs.WalkDirFunc {
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

func AddArtifact(artifactPath string, artifactAlias string) {
	CreateDirectoryIfNotPresent()
	CreateArtifactsFileIfNotPresent()
	CheckForConfigFile(artifactPath)
	// Set artifact alias to be some random UUID
	if artifactAlias == "" {
		artifactAlias = shortuuid.New()
	}
	// Getting artifact's directory name
	constructedPath := ConfigDir + artifactAlias
	//log.Printf("constructed path is :- %s", constructedPath)
	err := copy2.Copy(artifactPath, constructedPath)
	if err != nil {
		log.Fatal("Error while copying artifact!", err)
	}
	newArtifact := types.ArtifactEntry{ArtifactAlias: artifactAlias, ArtifactPath: constructedPath}
	// Check if artifact.json is empty
	file, _ := os.Stat(ArtifactsFile)
	// If file is empty
	if file.Size() == 0 {
		artifactEntries := make(types.Artifacts, 0)
		artifactEntries = append(artifactEntries, newArtifact)
		log.Println(artifactEntries)
		writeToArtifactsFile(artifactEntries)
	} else { // If file contains previous entries
		artifactRecords, _ := ioutil.ReadFile(ArtifactsFile)
		var artifacts types.Artifacts
		Must(json.Unmarshal(artifactRecords, &artifacts),
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
	//log.Println("artifact records :-" + string(artifactRecords))
	writeErr := ioutil.WriteFile(ArtifactsFile, artifactRecords, 0777)
	if writeErr != nil {
		log.Fatal("Error while marshalling artifacts")
	}
}

func CreateArtifactsFileIfNotPresent() {
	if _, err := os.Stat(ArtifactsFile); os.IsNotExist(err) {
		_, err := os.Create(ArtifactsFile)
		if err != nil {
			log.Fatal("Error while creating config directory", err)
		}
	}
}

func CreateDirectoryIfNotPresent() {
	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		err := os.Mkdir(ConfigDir, 0777)
		if err != nil {
			log.Fatal("Error while creating config directory", err)
		}
	}
}

func CheckForConfigFile(artifactPath string){
	//log.Printf("Searching for ~/.fire-up.json in %s\n", artifactPath + FireUpConfig)
	configDataRaw, err := ioutil.ReadFile(artifactPath + FireUpConfig)
	Must(err, "Error occurred while reading fire-up.json")
	var configData types.Config
	Must(json.Unmarshal(configDataRaw, &configData),
		"Error while unmarshalling fire-up.json")
}

func RunCommands(commandList []string, targetDir string){
	numberOfCommands := len(commandList)
	for index,command := range commandList {
		fmt.Printf("(#%d of %d) Running command %s",index+1,numberOfCommands,command)
		cmdArray := strings.Split(command, " ")
		var cmd *exec.Cmd
		if len(cmdArray)>1{
			cmd  = exec.Command(cmdArray[0],cmdArray[1:]...)
		} else {
			cmd = exec.Command(cmdArray[0])
		}
		cmd.Dir = targetDir
		Must(cmd.Run(),"Unable to run the initialization command!")
		//stdoutFromCmd,err := cmd.StdoutPipe()
		//Must(err,"Error while streaming output for this command")
		//streamProcessStdOut(stdoutFromCmd)
	}
}

func streamProcessStdOut(outStream io.ReadCloser){
	oneByte := make([]byte, 8)
	var fullOut string
	for {
		_, err := outStream.Read(oneByte)
		if err != nil {
			fmt.Printf(err.Error())
			break
		}
		fullOut += string(oneByte)
	}
	fmt.Print(fullOut)
}


package utils

import (
	"encoding/json"
	"fire-up/types"
	"fmt"
	"io/ioutil"
)

func ListAllArtifacts() {
	artifactRecords, _ := ioutil.ReadFile(ArtifactsFile)
	var artifacts types.Artifacts
	Must(json.Unmarshal(artifactRecords, &artifacts),
		"Error while un-marshalling artifacts")
	for _,artifactAlias := range artifacts{
		fmt.Printf("%v \n",artifactAlias.ArtifactAlias)
	}
}

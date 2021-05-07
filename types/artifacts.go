package types

import "log"

type ArtifactEntry struct {
	ArtifactAlias string `json:"artifact_alias"`
	ArtifactPath  string `json:"artifact_path"`
}

type Artifacts []ArtifactEntry

func (records Artifacts) CheckDuplicateRecords(entry ArtifactEntry) {
	for _, record := range records {
		if record == entry {
			log.Fatal("Err! Artifact with same name is already present")
		}
	}
}

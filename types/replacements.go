package types

import (
	"fmt"
	"log"
)

type Replacement struct {
	Placeholder string `json:"placeholder"`
	Description string `json:"desc"`
}

type Replacements []Replacement

func (replacements Replacements) ConstructReplacementMap() map[string]string {
	replacementMap := make(map[string]string)
	for _, replacement := range replacements {
		fmt.Print(replacement.Description+" ("+replacement.Placeholder+") ")
		var replaceWith string
		_, err := fmt.Scanln(&replaceWith)
		if err!=nil{
			log.Fatal("Error reading input")
		}
		replacementMap[replacement.Placeholder] = replaceWith
	}
	return replacementMap
}


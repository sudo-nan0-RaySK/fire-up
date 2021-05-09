package types

type Config struct {
	Type string `json:"type"`
	ReplacementList Replacements `json:"replacements"`
}

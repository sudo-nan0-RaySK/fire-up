package types

type Config struct {
	Type string `json:"type"`
	ReplacementList Replacements `json:"replacements"`
	InjectionList []InjectionMap `json:"injection_map"`
	CommandList []string `json:"command_list"`
}

package types

type GithubContentResponses []GithubContentResponse

type GithubContentResponse struct{
	Name string `json:"name"`
	Path string `json:"path"`
	SHA string `json:"sha"`
	Size int32 `json:"size"`
	Url string `json:"url"`
	HtmlUrl string `json:"html_url"`
	GitUrl string `json:"git_url"`
	DownloadUrl string `json:"download_url"`
	Type string `json:"type"`
	Links LinkType `json:"_links"`
}

type LinkType struct{
	Self string `json:"self"`
	Git string `json:"git"`
	Html string `json:"html"`
}
package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fire-up/types"
	"fmt"
	"github.com/cli/oauth/device"
	"github.com/google/go-github/v35/github"
	"github.com/gosuri/uilive"
	copy2 "github.com/otiai10/copy"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	deviceCodeEndpoint = "https://github.com/login/device/code"
	artifactsRepoName  = "fire-up-artifacts"
)

var tokenFile = ConfigDir + ".github/token.json"

var GithubClient *github.Client

var oauthConfig = &oauth2.Config{
	ClientID: "d48bf431f2c7ef866ec3",
	Scopes:   []string{"repo", "user"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
}

type LocalTokenSource struct{}

func (tokenSrc LocalTokenSource) Token() (*oauth2.Token, error) {
	token, err := getAccessToken()
	Must(err, "Error authenticating user")
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)
	if err != nil {
		return nil, err
	}
	oauth2Token, err := tokenSource.Token()
	Must(err, "Error getting oauth2Token")
	saveTokenInFileSystem(oauth2Token)
	return oauth2Token, nil
}

func CheckAuthenticated() {
	if GithubClient == nil {
		authenticateWithGithub()
	}
}

func GetOrCreateArtifactRepositoryIfNotPresent() *github.Repository {
	ctx := context.Background()
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	for {
		repoList, resp, err := GithubClient.Repositories.List(ctx, "", opts)
		Must(err, "Error fetching repository list")
		for _, repo := range repoList {
			//fmt.Printf("%s \n", *repo.Name)
			if *repo.Name == artifactsRepoName {
				fmt.Printf("%s \n", *repo.ContentsURL)
				return repo
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	// Create an artifacts repository
	fmt.Printf("Artifacts repository is not present, creating one...\n")
	artifactsRepo := &github.Repository{
		Name:    github.String(artifactsRepoName),
		Private: github.Bool(true),
	}
	artifactsRepRef, _, err := GithubClient.Repositories.Create(ctx, "", artifactsRepo)
	Must(err, "Error creating artifacts repository!")
	fmt.Printf("%s \n", *artifactsRepRef.ContentsURL)
	return artifactsRepRef
}

func CreateResource(artifactPath string, artifactName string, postUrl string) {
	CheckForConfigFile(artifactPath)
	client := getOAuthClient()
	Must(
		copy2.Copy(artifactPath, "/tmp/"+artifactName, copy2.Options{AddPermission: os.ModePerm}),
		"Error while copying artifact to /tmp",
	)
	fileObjects, err := getFiles("/tmp/" + artifactName)
	Must(err, "Error reading artifact!")
	for _, file := range fileObjects {
		postFiles(client, file, postUrl)
	}
	Must(os.RemoveAll("/tmp/"+artifactName),
		"Error cleaning up from /tmp")
}

func getOAuthClient() *http.Client {
	token, _ := fetchAuthTokenFromFileSystem()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token.AccessToken,
		},
	)
	client := oauth2.NewClient(context.Background(), tokenSource)
	return client
}

func postFiles(client *http.Client, filePath string, postUrl string) {
	actualPath := "/" + strings.Join(strings.Split(filePath, "/")[2:], "/")
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	Must(err, "Error reading file from the artifact")
	encData := base64.StdEncoding.EncodeToString(fileData)
	//fmt.Printf("filepath : %s \n", actualPath)
	//fmt.Printf("b64 string: %s \n", encData)
	urlParts := strings.Split(postUrl, "/")
	postUrl = strings.Join(urlParts[:len(urlParts)-1], "/")
	postURI := postUrl + actualPath
	fmt.Printf("url: %s\n", postURI)
	content := types.GithubContent{
		Content: encData,
		Message: "Added " + actualPath,
	}
	jsonString, marshallingErr := json.Marshal(content)
	Must(marshallingErr, "Error converting to json")
	request, reqErr := http.NewRequest(http.MethodPut, postURI, bytes.NewBuffer(jsonString))
	Must(reqErr, "Error while creating request!")
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	fmt.Printf("response: %s \n", response.Status)
}

func getFiles(artifactPath string) ([]string, error) {
	var fileObjects []string
	var onlyFiles []string
	Must(filepath.WalkDir(artifactPath, GetWalkAndCollect(&fileObjects)), "Error initiating WalkDir()")
	for _, fileObject := range fileObjects {
		if info, err := os.Stat(fileObject); err == nil && info != nil {
			if info.IsDir() {
				continue
			}
			onlyFiles = append(onlyFiles, fileObject)
		} else {
			return nil, err
		}
	}
	return onlyFiles, nil
}

func GetResource(artifactAlias string, getUrl string) (string, error) {
	client := getOAuthClient()
	urlParts := strings.Split(getUrl, "/")
	getUrl = strings.Join(urlParts[:len(urlParts)-1], "/")
	getURI := getUrl + "/" + artifactAlias
	writer := uilive.New()
	writer.Start()
	tmpPath, err := getResourcesRecursive("/tmp/"+artifactAlias, getURI, client, writer)
	writer.Stop()
	if err != nil {
		return "", err
	}
	return tmpPath, nil
}

func getResourcesRecursive(currFolder string, contentUrl string, client *http.Client, writer *uilive.Writer) (string, error) {
	Must(os.Mkdir(currFolder, os.ModePerm), "Error creating folder "+currFolder)
	response, errResponse := client.Get(contentUrl)
	Must(errResponse, "Error fetching resource meta-data")
	if response.StatusCode == 404 {
		return "", nil
	}
	responseBody, readErr := ioutil.ReadAll(response.Body)
	Must(readErr, "Error reading resource meta-data")
	var contentResponses types.GithubContentResponses
	Must(json.Unmarshal(responseBody, &contentResponses),
		"Error unmarshalling response")
	var folder string
	for _, contentResponse := range contentResponses {
		if contentResponse.Type == "dir" {
			folder = currFolder + "/" + contentResponse.Name
			_, _ = getResourcesRecursive(folder, contentResponse.Links.Self, client, writer)
		} else {
			_, err := fmt.Fprintf(writer, "Downloading %s \n", contentResponse.Name)
			if err != nil {
				return "", err
			}
			downloadFile(currFolder+"/"+contentResponse.Name, contentResponse.DownloadUrl, client)
		}
	}
	return currFolder, nil
}

func authenticateWithGithub() {
	ctx := context.Background()

	if token, foundStatus := fetchAuthTokenFromFileSystem(); foundStatus {
		// Checks if token needs to be refreshed
		tokenSource := oauth2.ReuseTokenSource(token, LocalTokenSource{})
		client := oauth2.NewClient(ctx, tokenSource)
		GithubClient = github.NewClient(client)
	} else {
		log.Printf("authConfig: %v\n", oauthConfig)
		ghClient, err := getGithubClient()
		GithubClient = ghClient
		if err != nil {
			panic(err)
		}
	}

	user, _, ghErr := GithubClient.Users.Get(ctx, "")
	if ghErr != nil {
		// Try once again (When access is revoked)
		ghClient, err := getGithubClient()
		GithubClient = ghClient
		if err != nil {
			panic(err)
		}
		user, _, ghErr = GithubClient.Users.Get(ctx, "")
	}
	Must(ghErr, "Error fetching github user!")
	log.Printf("Logged in as %s \n", *user.Name)

}

func fetchAuthTokenFromFileSystem() (*oauth2.Token, bool) {
	if _, err := os.Stat(tokenFile); !os.IsNotExist(err) {
		var token oauth2.Token
		byteToken, readErr := ioutil.ReadFile(tokenFile)
		Must(readErr, "Error reading token file!")
		Must(json.Unmarshal(byteToken, &token),
			"Error while unmarshalling token file!")
		return &token, true
	} else {
		return nil, false
	}
}

func saveTokenInFileSystem(token *oauth2.Token) {
	if _, foundErr := os.Stat(ConfigDir + ".github/"); os.IsNotExist(foundErr) {
		Must(os.Mkdir(ConfigDir+".github/", os.ModePerm),
			"Error while saving github access-token")
	}
	fmt.Printf("token: %v\n", *token)
	byteToken, err := json.Marshal(*token)
	Must(err, "Error in marshalling access-token")
	Must(ioutil.WriteFile(tokenFile, byteToken, os.ModePerm),
		"Error while writing github-access token to file-system!")
}

func getGithubClient() (*github.Client, error) {
	token, err := getAccessToken()
	Must(err, "Error authenticating user")
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)
	ctx := context.Background()
	client := oauth2.NewClient(ctx, tokenSource)
	oauth2Token, err := tokenSource.Token()
	Must(err, "Error getting oauth2Token")
	saveTokenInFileSystem(oauth2Token)
	if err != nil {
		log.Fatal("Error while authenticating", err)
	}
	return github.NewClient(client), nil
}

func getAccessToken() (string, error) {
	httpClient := http.DefaultClient
	code, err := device.RequestCode(httpClient, deviceCodeEndpoint, oauthConfig.ClientID, oauthConfig.Scopes)
	if err != nil {
		return "", err
	}
	fmt.Printf("Copy the device-access code: %s\n", code.UserCode)
	fmt.Printf("And then put it here (open manually if not opened automatically) %s\n", code.VerificationURI)
	actionErr := open.Run(code.VerificationURI)
	if actionErr != nil {
		fmt.Printf("Can't open browser automatically!\n" +
			"Please go to the above link and add the device-access code")
	}
	accessToken, err := device.PollToken(httpClient, oauthConfig.Endpoint.TokenURL, oauthConfig.ClientID, code)
	if err != nil {
		return "", err
	}
	return accessToken.Token, nil
}

func downloadFile(filePath string, downloadURL string, client *http.Client) {
	response, errResponse := client.Get(downloadURL)
	Must(errResponse, "Error downloading resource")
	fileData, readErr := ioutil.ReadAll(response.Body)
	Must(readErr, "Error reading resource")
	Must(ioutil.WriteFile(filePath, fileData, os.ModePerm),
		"Error writing resource")
}

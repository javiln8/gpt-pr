package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type PRCreateRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

func CreatePullRequest(owner, repo, branch, baseBranch, title, body string) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)

	pr := PRCreateRequest{
		Title: title,
		Body:  body,
		Head:  fmt.Sprintf("%s:%s", owner, branch), // Add the owner prefix to the branch name
		Base:  baseBranch,
	}

	jsonData, err := json.Marshal(pr)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to create PR, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

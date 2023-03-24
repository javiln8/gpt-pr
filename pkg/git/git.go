package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func GetGitDiff() (string, error) {
	return runGitCommand("diff")
}

func GetRepoOwnerAndName() (string, string, error) {
	remoteOutput, err := runGitCommand("config", "--get", "remote.origin.url")
	if err != nil {
		return "", "", fmt.Errorf("error getting repository owner and name: %v", err)
	}

	parts := strings.Split(strings.TrimSpace(remoteOutput), ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid remote URL format")
	}

	ownerAndRepo := strings.TrimSuffix(parts[1], ".git")
	ownerAndRepoParts := strings.Split(ownerAndRepo, "/")
	if len(ownerAndRepoParts) != 2 {
		return "", "", fmt.Errorf("invalid remote URL format")
	}

	return ownerAndRepoParts[0], ownerAndRepoParts[1], nil
}

func GetCurrentBranch() (string, error) {
	branchOutput, err := runGitCommand("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("error getting current branch: %v", err)
	}

	return strings.TrimSpace(branchOutput), nil
}

func CreateBranch(branchName string) (string, error) {
	createOutput, err := runGitCommand("checkout", "-b", branchName)
	if err != nil {
		return "", fmt.Errorf("error creating branch: %v", err)
	}

	return strings.TrimSpace(createOutput), nil
}

func PushBranch(branchName string) (string, error) {
	pushOutput, err := runGitCommand("push", "-u", "origin", branchName)
	if err != nil {
		return "", fmt.Errorf("error pushing branch: %v", err)
	}

	return strings.TrimSpace(pushOutput), nil
}

func runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running git command: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

func CommitChanges(commitMessage string) (string, error) {
	_, err := runGitCommand("add", ".")
	if err != nil {
		return "", fmt.Errorf("error staging changes: %v", err)
	}

	commitOutput, err := runGitCommand("commit", "-m", commitMessage)
	if err != nil {
		return "", fmt.Errorf("error committing changes: %v", err)
	}

	return strings.TrimSpace(commitOutput), nil
}

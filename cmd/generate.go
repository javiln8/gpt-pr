package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"regexp"

	"gpt-pr/pkg/chatgpt"
	"gpt-pr/pkg/git"
	"gpt-pr/pkg/github"
)

var verbose bool

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a PR title, branch name, and PR description based on a git diff",
	Long: `Generate a PR title, branch name, and PR description based on a git diff
using OpenAI's GPT-3 model.`,
	Run: func(cmd *cobra.Command, args []string) {
		gitDiff, err := git.GetGitDiff()
		if err != nil {
			fmt.Printf("Error getting git diff: %v\n", err)
			os.Exit(1)
		}

		if gitDiff == "" {
			fmt.Println("No git diff found. Please make changes to your files before running this command.")
			os.Exit(1)
		}

		owner, repo, err := git.GetRepoOwnerAndName()
		if err != nil {
			fmt.Printf("Error getting repository owner and name: %v\n", err)
			os.Exit(1)
		}

		baseBranch, err := git.GetCurrentBranch()
		if err != nil {
			fmt.Printf("Error getting current branch: %v\n", err)
			os.Exit(1)
		}

		apiKey := os.Getenv("GPT_API_KEY")
		if apiKey == "" {
			fmt.Println("Error: GPT_API_KEY environment variable not set.")
			os.Exit(1)
		}

		client := chatgpt.NewChatGPTClient(apiKey)
		generatedBranchName, prTitle, prDescription, err := client.GeneratePRDetailsGPT3(gitDiff)
		if err != nil {
			fmt.Printf("Error generating response: %v\n", err)
			os.Exit(1)
		}

		branchName, err := extractBranchName(generatedBranchName)
		if err != nil {
			fmt.Printf("Error extracting branch name from response: %v\n", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Println("Generated values:")
			fmt.Printf("Branch Name: %s\n", branchName)
			fmt.Printf("PR Title: %s\n", prTitle)
			fmt.Printf("PR Description:\n%s\n", prDescription)
		}

		if _, err := git.CreateBranch(branchName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = git.CommitChanges(prTitle)
		if err != nil {
			fmt.Printf("Error committing changes: %v\n", err)
			os.Exit(1)
		}

		if _, err := git.PushBranch(branchName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := github.CreatePullRequest(owner, repo, branchName, baseBranch, prTitle, prDescription); err != nil {
			fmt.Printf("Error creating pull request: %v\n", err)
			os.Exit(1)
		}

		err = github.CreatePullRequest(owner, repo, branchName, baseBranch, prTitle, prDescription)
		if err != nil {
			fmt.Printf("Error creating pull request: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Pull request created successfully.")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show the ChatGPT response")
}

func extractBranchName(response string) (string, error) {
	log.Print(response)

	// This regular expression matches branch names that follow the format <type>/<short-description>
	branchNamePattern := regexp.MustCompile(`\b[\w-]+\/[\w-]+`)

	branchName := branchNamePattern.FindString(response)
	if branchName != "" {
		return branchName, nil
	}

	return "", errors.New("no branch name found in response")
}

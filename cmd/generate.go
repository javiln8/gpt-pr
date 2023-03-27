package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"regexp"

	"gpt-pr/pkg/chatgpt"
	"gpt-pr/pkg/git"
	"gpt-pr/pkg/github"
)

type generateFlags struct {
	verbose        bool
	gptVersion     int
	addAttribution bool
}

var flags generateFlags

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

		var generatedBranchName, prTitle, prDescription string
		client := chatgpt.NewChatGPTClient(apiKey)
		if flags.gptVersion == 4 {
			generatedBranchName, prTitle, prDescription, err = client.GeneratePRDetailsGPT4(gitDiff)
		} else {
			generatedBranchName, prTitle, prDescription, err = client.GeneratePRDetailsGPT3(gitDiff)
		}

		// Add the attribution message if the flag is set
		if flags.addAttribution {
			prDescription += "\n\n---\n*This PR has been created with [gpt-pr](https://github.com/javiln8/gpt-pr).*"
		}

		branchName, err := extractBranchName(generatedBranchName)
		if err != nil {
			fmt.Printf("Error extracting branch name from response: %v\n", err)
			os.Exit(1)
		}

		if flags.verbose {
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
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().BoolVarP(&flags.verbose, "verbose", "v", false, "Show the ChatGPT response")
	generateCmd.Flags().IntVarP(&flags.gptVersion, "gpt-version", "g", 3, "Choose the GPT version (3 or 4), default is 3")
	generateCmd.Flags().BoolVarP(&flags.addAttribution, "add-attribution", "a", false, "Add an attribution message at the end of the PR summary")
}
func extractBranchName(response string) (string, error) {
	// This regular expression matches branch names that follow the format <type>/<short-description>
	branchNamePattern := regexp.MustCompile(`\b[\w-]+\/[\w-]+`)

	branchName := branchNamePattern.FindString(response)
	if branchName != "" {
		return branchName, nil
	}

	return "", errors.New("no branch name found in response")
}

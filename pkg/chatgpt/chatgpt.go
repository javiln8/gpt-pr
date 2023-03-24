package chatgpt

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/sashabaranov/go-openai"
)

type ChatGPTClient struct {
	client *openai.Client
}

func NewChatGPTClient(apiKey string) *ChatGPTClient {
	return &ChatGPTClient{
		client: openai.NewClient(apiKey),
	}
}

func (c *ChatGPTClient) generateGPT4Prompt(gitDiff string) string {
	return fmt.Sprintf(`I have the following git diff output from a code project:

%s

Analyze the git diff and generate a suitable PR title following the Conventional Commit format, a branch name, and a PR description based on the changes made. Please strictly follow this format for the response:

PR Title (Conventional Commit): <PR title>
Branch Name: <branch name>
PR Description (in Markdown):
<PR description>

For example:

PR Title (Conventional Commit): feat: add generate command for PR title, branch name, and description
Branch Name: feature/generate-command
PR Description (in Markdown):
- Added a generate command that generates a PR title, branch name, and description.
- Implemented ChatGPTClient for generating responses using OpenAI's GPT-4.

Do not include any additional information or analysis of the git diff in the response. Ensure the response is formatted correctly.`, gitDiff)
}

func (c *ChatGPTClient) generateBranchNamePrompt(gitDiff string) string {
	return fmt.Sprintf(`Analyze the following git diff output from a code project and generate a branch name based on the changes made:

%s

Ensure the branch name strictly follows the Conventional Commit format: <type>/<short-description>
For example: feature/add-login

Your response should be a single line containing only the branch name. Do not include any other information or context.

IMPORTANT: Please begin your response with "Branch Name: " followed by the actual branch name.`, gitDiff)
}

func (c *ChatGPTClient) generatePrTitlePrompt(gitDiff string) string {
	return fmt.Sprintf(`Analyze the following git diff output from a code project and generate a PR title following the Conventional Commit format, based on the changes made:

%s

Ensure the PR title strictly follows the Conventional Commit format: <type>: <short-description>
For example: feat: add login functionality
`, gitDiff)
}

func (c *ChatGPTClient) generatePrDescriptionPrompt(gitDiff string) string {
	return fmt.Sprintf(`Analyze the following git diff output from a code project and generate a PR description in Markdown format, based on the changes made:

%s

Ensure the PR description includes a clear and concise summary of the changes made, formatted as a bullet-point list in Markdown. Focus on the analysis of the git diff and avoid any personal language in the response.
`, gitDiff)
}

func (c *ChatGPTClient) GeneratePRDetailsGPT3(gitDiff string) (branchName, prTitle, prDescription string, err error) {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		prompt := c.generateBranchNamePrompt(gitDiff)
		branchName, err = c.generateResponseWithPrompt(prompt, openai.GPT3Dot5Turbo)
		if err != nil {
			log.Printf("Error generating branch name: %v\n", err)
		}
		log.Printf("Successfully generated branch name.")
	}()

	go func() {
		defer wg.Done()
		prompt := c.generatePrTitlePrompt(gitDiff)
		prTitle, err = c.generateResponseWithPrompt(prompt, openai.GPT3Dot5Turbo)
		if err != nil {
			log.Printf("Error generating PR title: %v\n", err)
		}
		log.Printf("Successfully generated PR title.")

	}()

	go func() {
		defer wg.Done()
		prompt := c.generatePrDescriptionPrompt(gitDiff)
		prDescription, err = c.generateResponseWithPrompt(prompt, openai.GPT3Dot5Turbo)
		if err != nil {
			log.Printf("Error generating PR description: %v\n", err)
		}
		log.Printf("Successfully generated PR description.")
	}()

	wg.Wait()

	if err != nil {
		return "", "", "", err
	}
	return branchName, prTitle, prDescription, nil
}

func (c *ChatGPTClient) GeneratePRDetailsGPT4(gitDiff string) (branchName, prTitle, prDescription string, err error) {
	prompt := c.generateGPT4Prompt(gitDiff)
	response, err := c.generateResponseWithPrompt(prompt, openai.GPT4)

	if err != nil {
		return "", "", "", err
	}

	// Parse the response to extract the PR title, branch name, and PR description
	// This code assumes the response is strictly formatted as described in the generateGPT4Prompt function
	lines := strings.Split(response, "\n")
	if len(lines) >= 4 {
		prTitle = strings.TrimSpace(strings.TrimPrefix(lines[0], "PR Title (Conventional Commit):"))
		branchName = strings.TrimSpace(strings.TrimPrefix(lines[1], "Branch Name:"))
		prDescription = strings.TrimSpace(lines[3])
		for _, line := range lines[4:] {
			prDescription += "\n" + strings.TrimSpace(line)
		}
	} else {
		return "", "", "", errors.New("unable to parse GPT-4 response")
	}

	return branchName, prTitle, prDescription, nil
}

func (c *ChatGPTClient) generateResponseWithPrompt(prompt string, model string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

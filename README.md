# GPT-PR

🤖 GPT-PR is a command-line tool that leverages the power of OpenAI's GPT-3 and GPT-4 models to automatically generate pull
request (PR) details, including the title, branch name, and description, based on a given Git diff.

All PRs of this repository have been created with GPT-PR.

## Features

* Generate PR title, branch name, and description using GPT-3 and GPT-4 models
* Automatically create and push a new branch with the generated details
* Create a pull request on GitHub using the generated details

## Installation

1. Clone the repository:

```bash
git clone https://github.com/javiln8/gpt-pr.git
```

2. Build the project:

```bash
cd gpt-pr
go build -o gpt-pr
```

3. Add the gpt-pr binary to your system's PATH.

## Setting up Environment Variables

Set up the `GPT_API_KEY` environment variable with your OpenAI API key:

```bash
export GPT_API_KEY=<your_openai_api_key>
```

Set up the `GITHUB_TOKEN` environment variable with your GitHub Personal Access Token:

```bash
export GITHUB_TOKEN=<your_github_personal_access_token>
```

## Usage

To use GPT-PR, you can run the following commands:

```bash
Usage:
  gpt-pr generate [flags]

Flags:
  -g, --gpt-version int   Choose the GPT version (3 or 4), default is 3 (default 3)
  -h, --help              help for generate
  -v, --verbose           Show the ChatGPT response
```

## GPT-3 Limitations

Although GPT-PR automates the PR creation process, it might be less reliable due to GPT-3's limitations. The tool makes
multiple calls to ChatGPT for better results, but the generated PR details may still vary in accuracy. However, the tool
serves as a valuable starting point for streamlining the PR creation process and can be further enhanced with more
sophisticated prompts or by utilizing newer AI models.

## Obtaining GPT API Keys

To get API keys for GPT-3 or GPT-4, you need to sign up for an account with OpenAI. Visit the OpenAI website to create
an account and obtain the necessary API keys ([link](https://platform.openai.com/account/api-keys)).

## Acknowledgments

This project was built with the help of OpenAI's GPT-4, which contributed to the code generation and decision-making
processes.

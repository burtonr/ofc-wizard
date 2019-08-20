package actions

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

type initialAnswers struct {
	Orchestrator  string
	RootDomain    string
	Registry      string
	SourceControl string
}

type githubAnswers struct {
	AppID          string
	WebhookSecret  string
	PrivateKeyFrom string
}

type gitlabAnswers struct {
	WebhookSecret string
	Instance      string
}

type oauthAnswers struct {
	ClientID string
	BaseURL  string
}

var (
	kubernetes          = "kubernetes"
	swarm               = "swarm"
	github              = "github"
	gitlab              = "gitlab"
	createAppHelpText   = "Create a Github app by following the instructions in the docs: https://docs.openfaas.com/openfaas-cloud/self-hosted/github/"
	createOAuthHelpText = "Create the OAuth App on your source control management system"
)

// StartInstall will create and ask the survey questions to generate a yml file for use with the ofc-bootstrap tool
func StartInstall() {
	fmt.Println("Starting wizard")
	initAnswers, err := askInitialQuestions()

	if err != nil {
		fmt.Println(err.Error())
	}

	if initAnswers.SourceControl == github {
		ghAnswers := askGithubQuestions()

		fmt.Println(ghAnswers)
	} else if initAnswers.SourceControl == gitlab {
		glAnswers := askGitLabQuestions()

		fmt.Println(glAnswers)
	}

	oauthAnswers := askOAuthQuestions(initAnswers.SourceControl)

	// Println statements to avoid "unused" errors (temporary)
	fmt.Println("Orchestration:", initAnswers.Orchestrator)
	fmt.Println("Root Domain:", initAnswers.RootDomain)
	fmt.Println("Registry:", initAnswers.Registry)
	fmt.Println("Source Control:", initAnswers.SourceControl)
	fmt.Println("OAuth Client ID:", oauthAnswers.ClientID)
}

type answers struct {
	RootDomain string
}

func askInitialQuestions() (*initialAnswers, error) {
	var questions = []*survey.Question{
		{
			Name: "Orchestrator",
			Prompt: &survey.Select{
				Message: "Select an orchestration provider",
				Options: []string{kubernetes, swarm},
			},
			Validate: survey.Required,
		},
		{
			Name:     "RootDomain",
			Prompt:   &survey.Input{Message: "Root Domain (eg: faas.example.com):"},
			Validate: survey.Required,
		},
		{
			Name:   "Registry",
			Prompt: &survey.Input{Message: "Registry to publish images (eg: docker.io/your-name/):"},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || !strings.HasSuffix(str, "/") {
					return errors.New("The registry address must end with a '/'")
				}
				return nil
			},
		},
		{
			Name: "SourceControl",
			Prompt: &survey.Select{
				Message: "Select an source control management",
				Options: []string{github, gitlab},
			},
			Validate: survey.Required,
		},
	}

	a := &initialAnswers{}

	if err := survey.Ask(questions, a); err != nil {
		return nil, err
	}
	return a, nil
}

func askGithubQuestions() *githubAnswers {
	var preReqQuestion = &survey.Confirm{Message: "Do you have your Github App setup already?"}

	appCreated := false
	survey.AskOne(preReqQuestion, &appCreated, nil)

	if !appCreated {
		fmt.Printf("\n%s\n\n", createAppHelpText)
	}

	var questions = []*survey.Question{
		{
			Name:     "AppID",
			Prompt:   &survey.Input{Message: "Github App ID:"},
			Validate: survey.Required,
		},
		{
			Name:   "WebhookSecret",
			Prompt: &survey.Input{Message: "Enter your webhook secret (leave blank for a random value):"},
		},
		{
			Name: "PrivateKeyFrom",
			Prompt: &survey.Input{
				Message: "File location of the private key:",
				Help:    "Enter the full path of the private key downloaded from the Github app (eg: ~/Downloads/private-key.pem)",
			},
			Validate: survey.Required,
		},
	}

	a := &githubAnswers{}

	if err := survey.Ask(questions, a); err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return a
}

func askGitLabQuestions() *gitlabAnswers {
	var questions = []*survey.Question{
		{
			Name: "WebhookSecret",
			Prompt: &survey.Input{
				Message: "Enter your webhook secret (leave blank for a random value):",
				Help:    "Enter the full path of the private key downloaded from GitLab (eg: ~/Downloads/private-key.pem)",
			},
		},
		{
			Name: "Instance",
			Prompt: &survey.Input{
				Message: "Enter the public URL for your GitLab instance (with trailing slash):",
				Help:    "Enter the full URL of your public GitLab (eg: https://gitlab.example.com/)",
			},
			Validate: survey.Required,
		},
	}

	a := &gitlabAnswers{}

	if err := survey.Ask(questions, a); err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return a
}

func askOAuthQuestions(scm string) *oauthAnswers {
	var preReqQuestion = &survey.Confirm{Message: "Have you created your OAuth App already?"}

	appCreated := false
	survey.AskOne(preReqQuestion, &appCreated, nil)

	if !appCreated {
		fmt.Printf("\n%s\n\n", createOAuthHelpText)
	}

	var questions = []*survey.Question{
		{
			Name:   "ClientID",
			Prompt: &survey.Input{Message: "Enter the OAuth App ID:"},
		},
	}

	if scm == gitlab {
		baseQ := survey.Question{Name: "BaseURL", Prompt: &survey.Input{Message: "Enter your OAuth Provider Base URL:"}}
		questions = append(questions, &baseQ)
	}

	a := &oauthAnswers{}

	if err := survey.Ask(questions, a); err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return a
}

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

type configAnswers struct {
	AuditURL        string
	CustomersURL    string
	UseDockerfile   bool
	OFVersion       string
	ScaleZero       bool
	NetworkPolicies bool
}

var (
	kubernetes          = "kubernetes"
	swarm               = "swarm"
	github              = "github"
	gitlab              = "gitlab"
	defaultVersion      = "0.9.5"
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

	// s3
	// tls
	// dns-service

	finalConfigAnswers := askFinalConfigQuestions()

	// Println statements to avoid "unused" errors (temporary)
	fmt.Println("Orchestration:", initAnswers.Orchestrator)
	fmt.Println("Root Domain:", initAnswers.RootDomain)
	fmt.Println("Registry:", initAnswers.Registry)
	fmt.Println("Source Control:", initAnswers.SourceControl)
	fmt.Println("OAuth Client ID:", oauthAnswers.ClientID)
	fmt.Println("Final Answers:", finalConfigAnswers)
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

func askFinalConfigQuestions() *configAnswers {
	answers := &configAnswers{}
	answers.AuditURL = "http://gateway.openfaas:8080/function/echo"
	answers.OFVersion = defaultVersion

	var customAuditQuestion = &survey.Confirm{Message: "Would you like to use a custom audit trail URL (ie post to Slack)?"}

	customAudit := false
	survey.AskOne(customAuditQuestion, &customAudit, nil)

	if customAudit {
		var auditURLQuestion = &survey.Input{Message: "URL to post audit trail message to:"}
		survey.AskOne(auditURLQuestion, &answers.AuditURL, nil)
	}

	// customers
	var custURLQuestion = &survey.Input{
		Message: "URL of the customers access control list:",
		Help:    "The raw text file, or Github raw URL of allowed users. This must be a public repo",
	}

	survey.AskOne(custURLQuestion, &answers.CustomersURL, nil)

	// dockerfile
	var dockerfileQuestion = &survey.Confirm{
		Message: "Would you like to enable the Dockerfile template?",
		Help:    "This will allow templates built using dockerfile to be deployed which will allow ANY workload to be built and run. Use with caution",
	}

	survey.AskOne(dockerfileQuestion, &answers.UseDockerfile, nil)

	// scale-zero
	var scaleZeroQuestion = &survey.Confirm{
		Message: "Would you like to enable scale-to-zero as the default?",
		Help:    "With this enabled, all functions will scale to zero. To turn off, add a label 'com.openfaas.scale.zero: false'",
	}

	survey.AskOne(scaleZeroQuestion, &answers.ScaleZero, nil)

	// ofc version
	var versionQuestion = &survey.Input{
		Message: fmt.Sprintf("Enter the version of OpenFaaS Cloud to use (default: %s)", defaultVersion),
		Help:    "See available versions here: https://github.com/openfaas/openfaas-cloud/releases/",
	}

	survey.AskOne(versionQuestion, &answers.OFVersion, nil)

	// network policies
	var netPoliciesQuestion = &survey.Confirm{
		Message: "Would you like to enable network policies (restrict functions from calling the openfaas namespace)",
		Help:    "Prevents functions from talkking to the openfaas namespace, and to each other. Use the ingress address for the gateway or external IP instead",
	}

	survey.AskOne(netPoliciesQuestion, &answers.NetworkPolicies, nil)

	return answers
}

package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/burtonr/ofc-wizard/types"
	"gopkg.in/AlecAivazis/survey.v1"
)

type initialAnswers struct {
	Orchestrator  string
	RootDomain    string
	Registry      string
	SourceControl string
	EnableOAuth   bool
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

type storageAnswers struct {
	URL       string
	Region    string
	EnableTLS bool
	Bucket    string
}

type dnsAnswers struct {
	Name       string
	AccessFile types.Literal
	Filters    []string
	Namespace  string
}

type tlsAnswers struct {
	Enabled      bool
	IssuerType   string
	EmailAddress string
	DNSService   string
	ProjectID    string // Used only for GCP DNS
	Region       string // Used only for AWS DNS
	AccessKey    string // Used only for AWS DNS
}

type configAnswers struct {
	AuditURL        string
	CustomersURL    string
	UseDockerfile   bool
	OFVersion       string
	ScaleZero       bool
	NetworkPolicies bool
	Ingress         string
}

type dnsProvider struct {
	FriendlyName string
	Name         string
	Filter       []string
	File         string
	HelpText     string
}

var (
	kubernetes          = "kubernetes"
	swarm               = "swarm"
	github              = "github"
	gitlab              = "gitlab"
	defaultVersion      = "0.9.7"
	createAppHelpText   = "Create a Github app by following the instructions in the docs: https://docs.openfaas.com/openfaas-cloud/self-hosted/github/"
	createOAuthHelpText = "Create the OAuth App on your source control management system"
	digOceanDNS         = dnsProvider{
		FriendlyName: "DigitalOcean",
		Name:         "digitalocean-dns",
		Filter:       []string{"do_dns01"},
		File:         "access-token",
		HelpText:     "Create a Personal Access Token and save it into a file, with no new lines",
	}
	gCloudDNS = dnsProvider{
		FriendlyName: "Google Cloud",
		Name:         "clouddns-service-account",
		Filter:       []string{"gcp_dns01"},
		File:         "service-account.json",
		HelpText:     "Create a service account for DNS management and export it",
	}
	awsDNS = dnsProvider{
		FriendlyName: "AWS Route 53",
		Name:         "route53-credentials-secret",
		Filter:       []string{"route53_dns01"},
		File:         "secret-access-key",
		HelpText:     "Create a role and download it's secret access key",
	}
)

// GenerateYaml will create and ask the survey questions to generate a yml file for use with the ofc-bootstrap tool
func GenerateYaml() {
	yml := CreateInitFile()

	initAnswers, err := askInitialQuestions()

	if err != nil {
		fmt.Println(err.Error())
	}

	yml.Orchestration = initAnswers.Orchestrator
	yml.RootDomain = initAnswers.RootDomain
	yml.Registry = initAnswers.Registry
	yml.SCM = initAnswers.SourceControl
	yml.EnableOAuth = initAnswers.EnableOAuth

	if initAnswers.SourceControl == github {
		ghAnswers := askGithubQuestions()
		yml.Github = types.Github{
			AppID: ghAnswers.AppID,
		}
	} else if initAnswers.SourceControl == gitlab {
		glAnswers := askGitLabQuestions()
		yml.GitLab = types.GitLab{
			GitLabInstance: glAnswers.Instance,
		}
	}

	if initAnswers.EnableOAuth {
		oAuthAnswers := askOAuthQuestions(initAnswers.SourceControl)
		yml.OAuth = types.OAuth{
			ClientID:             oAuthAnswers.ClientID,
			OAuthProviderBaseURL: oAuthAnswers.BaseURL,
		}
	}

	storageAnswers := askStorageQuestions()
	yml.S3 = types.Storage{
		S3URL:    storageAnswers.URL,
		S3Region: storageAnswers.Region,
		S3Bucket: storageAnswers.Bucket,
		S3TLS:    storageAnswers.EnableTLS,
	}

	dnsAnswers := askDNSQuestions()
	tlsAnswers := askTLSQuestions(dnsAnswers.Name)

	yml.TLS = tlsAnswers.Enabled

	if yml.TLS {
		yml.TLSConfig = types.TLSConfig{
			DNSService:  tlsAnswers.DNSService,
			Email:       tlsAnswers.EmailAddress,
			IssuerType:  tlsAnswers.IssuerType,
			ProjectID:   tlsAnswers.ProjectID,
			Region:      tlsAnswers.Region,
			AccessKeyID: tlsAnswers.AccessKey,
		}
	}

	// finalConfigAnswers := askFinalConfigQuestions()

	// // Println statements to avoid "unused" errors (temporary)
	// fmt.Println("Orchestration:", initAnswers.Orchestrator)
	// fmt.Println("Root Domain:", initAnswers.RootDomain)
	// fmt.Println("Registry:", initAnswers.Registry)
	// fmt.Println("Source Control:", initAnswers.SourceControl)
	// fmt.Println("Storage Answers:", storageAnswers)
	// fmt.Println("DNS Answers:", dnsAnswers)
	// fmt.Println("TLS Answers:", tlsAnswers)
	// fmt.Println("Final Answers:", finalConfigAnswers)
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
		{
			Name:   "EnableOAuth",
			Prompt: &survey.Confirm{Message: "Would you like to enable OAuth so only those with Github/Gitlab accounts may log in (recommended)"},
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
				Message: "Enter the path of the file for your private key:",
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
				Message: "Enter the path to the file with your webhook secret (leave blank for a random value):",
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

func askStorageQuestions() *storageAnswers {
	answers := &storageAnswers{
		URL:       "cloud-minio.openfaas.svc.cluster.local:9000",
		Region:    "us-east-1",
		EnableTLS: false,
		Bucket:    "pipeline"}

	customStorageQuestion := &survey.Confirm{Message: "Would you like to use custom storage (S3 compatible) for logs from buildkit? (not recommended)"}
	customStorage := false
	survey.AskOne(customStorageQuestion, &customStorage, nil)

	if customStorage {
		storageQuestions := []*survey.Question{
			{
				Name:   "URL",
				Prompt: &survey.Input{Message: "Enter the Base URL for your storage location:"},
			},
			{
				Name:   "Region",
				Prompt: &survey.Input{Message: "Enter the S3 region associated with the storage location:"},
			},
			{
				Name:   "Bucket",
				Prompt: &survey.Input{Message: "Enter the bucket name to store the buildkit logs:"},
			},
			{
				Name:   "EnableTLS",
				Prompt: &survey.Confirm{Message: "Would you like to enable TLS encryption on the requests to the storage?"},
			},
		}

		if err := survey.Ask(storageQuestions, answers); err != nil {
			fmt.Println(err.Error())
			return nil
		}
	}

	return answers
}

func askDNSQuestions() *dnsAnswers {
	var providers = map[string]dnsProvider{digOceanDNS.FriendlyName: digOceanDNS, gCloudDNS.FriendlyName: gCloudDNS, awsDNS.FriendlyName: awsDNS}
	dnsNames := []string{digOceanDNS.FriendlyName, gCloudDNS.FriendlyName, awsDNS.FriendlyName}

	nameQuestion := &survey.Select{Message: "Select a DNS provider:", Options: dnsNames}
	var name string
	survey.AskOne(nameQuestion, &name, nil)

	fileQuestion := &survey.Input{
		Message: "Enter the path to the file containing the DNS provider credentials:",
		Help:    providers[name].HelpText,
	}
	var fileName string
	survey.AskOne(fileQuestion, &fileName, nil)

	selectedProvider := providers[name]
	resultFileLit := types.Literal{Name: selectedProvider.File, Value: fileName}
	result := &dnsAnswers{Name: selectedProvider.Name, Filters: selectedProvider.Filter, AccessFile: resultFileLit}

	return result
}

func askTLSQuestions(dnsService string) *tlsAnswers {
	answers := &tlsAnswers{Enabled: false}

	enableTLSQuestion := &survey.Confirm{Message: "Would you like to enable TLS? (recommended)"}
	survey.AskOne(enableTLSQuestion, &answers.Enabled, nil)

	if !answers.Enabled {
		return answers
	}

	tlsConfigQuestions := []*survey.Question{
		{
			Name:   "EmailAddress",
			Prompt: &survey.Input{Message: "Enter the email address to use for registering the domain:"},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || !strings.Contains(str, "@") {
					return errors.New("You must provide a valid email address")
				}
				return nil
			},
		},
		{
			Name: "IssuerType",
			Prompt: &survey.Select{
				Message: "Choose which type of certificate to issue (recommend prod):",
				Options: []string{"prod", "staging"},
			},
		},
	}

	survey.Ask(tlsConfigQuestions, answers)

	switch dnsService {
	case gCloudDNS.Name:
		survey.AskOne(&survey.Input{Message: "Enter the Project ID:"}, &answers.ProjectID, nil)
		break
	case awsDNS.Name:
		awsConfigQuestions := []*survey.Question{
			{Name: "Region", Prompt: &survey.Input{Message: "Enter the AWS Region:"}},
			{Name: "AccessKey", Prompt: &survey.Input{Message: "Enter the Access Key ID:"}},
		}
		survey.Ask(awsConfigQuestions, answers)
		break
	case digOceanDNS.Name:
		return answers
	}

	return answers
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
		Help:    "The raw text file, or Github raw URL of allowed users. This must be a public endpoint",
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
		Message: fmt.Sprintf("Enter the version of OpenFaaS Cloud to use (blank for default: %s)", defaultVersion),
		Help:    "See available versions here: https://github.com/openfaas/openfaas-cloud/releases/",
	}

	survey.AskOne(versionQuestion, &answers.OFVersion, nil)

	// network policies
	var netPoliciesQuestion = &survey.Confirm{
		Message: "Would you like to enable network policies (restrict functions from calling the openfaas namespace, recommended)",
		Help:    "Prevents functions from talkking to the openfaas namespace, and to each other. Use the ingress address for the gateway or external IP instead",
	}

	survey.AskOne(netPoliciesQuestion, &answers.NetworkPolicies, nil)

	var ingressQuestion = &survey.Select{
		Message: "Choose which type of ingress to use:",
		Options: []string{"loadbalancer", "host"},
	}

	survey.AskOne(ingressQuestion, &answers.Ingress, nil)
	// TODO: custom templates

	return answers
}

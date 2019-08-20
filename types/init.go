package types

type InitYaml struct {
	Orchestration        string         `yaml:"orchestration"`
	Secrets              []Secret       `yaml:"secrets"`
	Registry             string         `yaml:"registry"`
	RootDomain           string         `yaml:"root_domain"`
	Ingress              string         `yaml:"ingress"`
	Deployment           DeploymentOpts `yaml:"deployment"`
	SCM                  string         `yaml:"scm"`
	Github               Github         `yaml:"github"`
	GitLab               GitLab         `yaml:"gitlab"`
	OAuth                OAuth          `yaml:"oauth"`
	Slack                Slack          `yaml:"slack"`
	CustomersURL         string         `yaml:"cusomter_url"`
	S3                   Storage        `yaml:"s3"`
	EnableOAuth          bool           `yaml:"enable_oauth"`
	TLS                  bool           `yaml:"tls"`
	TLSConfig            TLSConfig      `yaml:"tls_config"`
	EnableDockerFile     bool           `yaml:"enable_dockerfile_lang"`
	ScaleToZero          bool           `yaml:"scale_to_zero"`
	OpenFaaSCloudVersion string         `yaml:"openfaas_cloud_version"`
	NetworkPolicies      bool           `yaml:"network_policies"`
}

type Secret struct {
	Name      string      `yaml:"name"`
	Literals  []Literal   `yaml:"literals"`
	Filters   []string    `yaml:",omitempty"`
	Namespace string      `yaml:",omitempty"`
	Files     []FileValue `yaml:",omitempty"`
}

type Literal struct {
	Name  string `yaml:"name"`
	Value string `yaml:",omitempty"`
}

type FileValue struct {
	Name         string `yaml:"name"`
	ValueFrom    string `yaml:"value_from"`
	ValueCommand string `yaml:"value_command"`
}

type DeploymentOpts struct {
	CustomTemplates []string `yaml:"custom_templates"`
}

type Github struct {
	AppID string `yaml:"app_id"`
}

type GitLab struct {
	GitLabInstance string `yaml:"gitlab_instance"`
}

type OAuth struct {
	ClientID             string `yaml:"client_id"`
	OAuthProviderBaseURL string `yaml:"oauth_provider_base_url,omitempty"`
}

type Slack struct {
	URL string `yaml:"url"`
}

type Storage struct {
	S3URL    string `yaml:"s3_url"`
	S3Region string `yaml:"s3_region"`
	S3TLS    bool   `yaml:"s3_tls"`
	S3Bucket string `yaml:"s3_bucket"`
}

type TLSConfig struct {
	IssuerType  string `yaml:"issuer_type"`
	Email       string `yaml:"email"`
	DNSService  string `yaml:"dns_service"`
	ProjectID   string `yaml:"project_id,omitempty"`
	Region      string `yaml:"region,omitempty"`
	AccessKeyID string `yaml:"access_key_id,omitempty"`
}

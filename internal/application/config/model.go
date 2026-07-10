package config

type Configuration struct {
	Instances []Instance
}

type Instance struct {
	Name       string
	Endpoint   Endpoint
	Credential CredentialReference
}

type Endpoint string

type CredentialReference string

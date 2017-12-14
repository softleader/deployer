package datamodels

type Deployment struct {
	Project     string `json:"project"`
	EurekaPort  int    `json:"eurekaPort"`
	KibanaPort  int    `json:"kibanaPort"`
	GatewayPort int    `json:"gatewayPort"`
	PublishPort int    `json:"publishPort"`
	Volume0     string `json:"volume0"`
	Net0        string `json:"net0"`
	Yaml        string `json:"yaml"`
}

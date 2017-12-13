package datamodels

type Deployment struct {
	EurekaPort  int    `json:"eurekaPort"`
	KibanaPort  int    `json:"kibanaPort"`
	GatewayPort int    `json:"gatewayPort"`
	PublishPort int    `json:"publishPort"`
	Mount       string `json:"mount"`
	Network     string `json:"network"`
	Yaml        string `json:"yaml"`
}

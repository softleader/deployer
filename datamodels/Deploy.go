package datamodels

type Deploy struct {
	CleanUp     bool     `json:cleanUp`
	Project     string   `json:"project"`
	EurekaPort  int      `json:"eurekaPort"`
	KibanaPort  int      `json:"kibanaPort"`
	GatewayPort int      `json:"gatewayPort"`
	Dev         string   `json:"dev"`
	Volume0     string   `json:"volume0"`
	Net0        string   `json:"net0"`
	Yaml        []string `json:"yaml"`
	Silence     bool     `json:"silence"`
}

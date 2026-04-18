package config

type Config struct {
	Language string `json:"language" mapstructure:"language"`
	Format   string `json:"format" mapstructure:"format"`
	MaxNodes int    `json:"max_nodes" mapstructure:"max_nodes"`
}

func Default() Config {
	return Config{
		Language: "auto",
		Format:   "mermaid",
		MaxNodes: 200,
	}
}

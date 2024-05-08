package config

type Language struct {
	Name    string         `mapstructure:"name"`
	Options map[string]any `mapstructure:"extends"`
	Extends []string       `mapstructure:"extends"`
	Models  []Model        `mapstructure:"models"`
}

type Model struct {
	Name    string         `mapstructure:"name"`
	Tags    []string         `mapstructure:"tags"`
	Options map[string]any `mapstructure:"options"`
}

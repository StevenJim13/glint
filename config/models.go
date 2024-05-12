package config

import "fmt"

type Language struct {
	Name    string
	Options map[string]any
	Extends []string
	Models  []Model
}

type Model struct {
	Name    string
	Tags    []string
	Options map[string]any
}

func (m Model) String() string {
	return fmt.Sprintf("<Model: %s>", m.Name)
}

package config

import "fmt"

type Language struct {
	Name    string
	Options map[string]any
	Extends []string
	Models  []Model
}

func (l Language) String() string {
	return fmt.Sprintf("<Language: %s, extends:%s>", l.Name, l.Extends)
}

type Model struct {
	Name    string
	Tags    []string
	Options map[string]any
}

func (m Model) String() string {
	return fmt.Sprintf("<Model: %s>", m.Name)
}

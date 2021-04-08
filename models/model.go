package models

type Instructions struct {
	CommonOverlays []Overlay  `yaml:"commonOverlays,omitempty"`
	YamlFiles      []YamlFile `yaml:"yamlFiles,omitempty"`
}

type Overlay struct {
	Name          string          `yaml:"name,omitempty"`
	Query         string          `yaml:"query,omitempty"`
	Value         interface{}     `yaml:"value,omitempty"`
	Action        string          `yaml:"action,omitempty"`
	DocumentQuery []DocumentQuery `yaml:"documentQuery,omitempty"`
	OnMissing     OnMissing       `yaml:"onMissing,omitempty"`
	DocumentIndex []string        `yaml:"documentIndex,omitempty"`
}

type DocumentQuery struct {
	Conditions []Condition `yaml:"conditions,omitempty"`
}

type Condition struct {
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type YamlFile struct {
	Name      string    `yaml:"name,omitempty"`
	Path      string    `yaml:"path,omitempty"`
	Overlays  []Overlay `yaml:"overlays,omitempty"`
	Documents []Overlay `yaml:"documents,omitempty"`
}

type OnMissing struct {
	Action     string `yaml:"action,omitempty"`
	InjectPath string `yaml:"injectPath,omitempty"`
}

type YamlDocuments struct {
	Path string `yaml:"path,omitempty"`
}

package models

type Instructions struct {
	CommonOverlays []Overlay       `yaml:"commonOverlays"`
	YamlFiles      []YamlFile      `yaml:"yamlFiles"`
	Documents      []YamlDocuments `yaml:"documents"`
}

type Overlay struct {
	Name          string          `yaml:"name"`
	Query         string          `yaml:"query"`
	Value         interface{}     `yaml:"value"`
	Action        string          `yaml:"action"`
	DocumentQuery []DocumentQuery `yaml:"documentQuery"`
	OnMissing     OnMissing       `yaml:"onMissing"`
	DocumentIndex []string        `yaml:"documentIndex"`
}

type DocumentQuery struct {
	Conditions []Condition `yaml:"conditions"`
}

type Condition struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type YamlFile struct {
	Name      string    `yaml:"name"`
	Path      string    `yaml:"path"`
	Overlays  []Overlay `yaml:"overlays"`
	Documents []Overlay `yaml:"documents"`
}

type OnMissing struct {
	Action     string `yaml:"action"`
	InjectPath string `yaml:"injectPath"`
}

type YamlDocuments struct {
	Path string `yaml:"path"`
}

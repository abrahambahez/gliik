package instruction

type Instruction struct {
	Name       string
	Path       string
	SystemText string
	Meta       Meta
}

type Meta struct {
	Version     string   `yaml:"version"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags,omitempty"`
	Lang        string   `yaml:"lang,omitempty"`
	Type        string   `yaml:"type,omitempty"`
	Name        string   `yaml:"name,omitempty"`
	Assets      []string `yaml:"assets,omitempty"`
}

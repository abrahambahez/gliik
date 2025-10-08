package instruction

type Instruction struct {
	Name       string
	Path       string
	SystemText string
	Meta       Meta
}

type Meta struct {
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

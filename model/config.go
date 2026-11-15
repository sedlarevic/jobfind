package model

type CategoryConfig struct {
	Fallback   string     `json:"fallback"`
	Categories []Category `json:"categories"`
}

type ExperienceConfig struct {
	Fallback    string       `json:"fallback"`
	Experiences []Experience `json:"experiences"`
}

type WorkingModelConfig struct {
	Fallback      string         `json:"fallback"`
	WorkingModels []WorkingModel `json:"working_models"`
}

type SourceConfig struct {
	Fallback string   `json:"fallback"`
	Sources  []Source `json:"sources"`
}

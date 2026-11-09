package resolver

import (
	"jobfind/model"
	"strings"
)

func RetrieveExperience(title string, description string) model.Experience {
	title = strings.ToLower(title)
	description = strings.ToLower(description)

	if experience := retrieveExperienceFromText(title); experience != model.ExperienceUnknown {
		return experience
	}

	return retrieveExperienceFromText(description)
}

func retrieveExperienceFromText(text string) model.Experience {
	switch {
	case strings.Contains(text, "executive"):
		return model.ExperienceExecutive
	case strings.Contains(text, "lead"):
		return model.ExperienceLead
	case strings.Contains(text, "senior"),
		strings.Contains(text, "5+ years"):
		return model.ExperienceSenior
	case strings.Contains(text, "mid"),
		strings.Contains(text, "medior"),
		strings.Contains(text, "2+ years"),
		strings.Contains(text, "3+ years"),
		strings.Contains(text, "4+ years"),
		strings.Contains(text, "intermediate"):
		return model.ExperienceMid
	case strings.Contains(text, "junior"),
		strings.Contains(text, "associate"):
		return model.ExperienceJunior
	case strings.Contains(text, "intern"),
		strings.Contains(text, "entry"):
		return model.ExperienceIntern
	default:
		return model.ExperienceUnknown
	}
}

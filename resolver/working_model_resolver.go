package resolver

import (
	"jobfind/model"
	"strings"
)

func RetrieveWorkingModel(title string, description string) model.WorkingModel {
	title = strings.ToLower(title)
	description = strings.ToLower(description)

	if workingModel := retrieveWorkingModelFromText(title); workingModel != model.WorkingModelUnknown {
		return workingModel
	}

	return retrieveWorkingModelFromText(description)
}

func retrieveWorkingModelFromText(text string) model.WorkingModel {
	switch {
	case strings.Contains(text, "hybrid"),
		strings.Contains(text, "flexible working hours"):
		return model.WorkingModelHybrid
	case strings.Contains(text, "remote"):
		return model.WorkingModelRemote
	case strings.Contains(text, "on-site"):
		return model.WorkingModelOnSite
	default:
		return model.WorkingModelUnknown
	}
}

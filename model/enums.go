package model

type Source string

const (
	SourceUnknown         Source = "unknown"
	SourceDirect          Source = "direct"
	SourceLever           Source = "lever"
	SourceGreenhouse      Source = "greenhouse"
	SourceBambooHR        Source = "bamboohr"
	SourceAshby           Source = "ashby"
	SourceSmartRecruiters Source = "smart_recruiters"
)

type WorkingModel string

const (
	WorkingModelUnknown WorkingModel = "unknown"
	WorkingModelOnSite  WorkingModel = "on_site"
	WorkingModelHybrid  WorkingModel = "hybrid"
	WorkingModelRemote  WorkingModel = "remote"
)

type Experience string

const (
	ExperienceUnknown   Experience = "unknown"
	ExperienceIntern    Experience = "intern"
	ExperienceJunior    Experience = "junior"
	ExperienceMid       Experience = "mid"
	ExperienceSenior    Experience = "senior"
	ExperienceStaff     Experience = "staff"
	ExperiencePrincipal Experience = "principal"
	ExperienceLead      Experience = "lead"
	ExperienceExecutive Experience = "executive"
)

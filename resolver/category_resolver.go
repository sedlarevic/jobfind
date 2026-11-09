package resolver

import (
	"jobfind/model"
	"strings"
)

// CATEGORIES

var kwCategorySoftware = []string{
	"backend",
	"frontend",
	"fullstack",
	"software engineer",
	"software developer",
}

var kwCategoryData = []string{
	"data analyst",
	"data scientist",
	"data engineer",
}

var kwCategoryAI = []string{
	"machine learning",
	"ml engineer",
	"ai engineer",
	"artificial intelligence",
}

var kwCategoryMarketing = []string{
	"marketing",
	"seo",
	"copywriter",
	"social media",
}

var kwCategoryManagement = []string{
	"product manager",
	"project manager",
	"engineering manager",
	"scrum master",
}

//SUBCATEGORIES

var kwSubCategoryBackend = []string{
	"backend",
	"back-end",
	"server-side",
	"server side",
	"api developer",
	"api engineer",
}

var kwSubCategoryFrontend = []string{
	"frontend",
	"front-end",
	"front end",
	"ui developer",
	"web developer",
}

var kwSubCategoryFullstack = []string{
	"fullstack",
	"full-stack",
	"full stack",
}

var kwSubCategoryDevOps = []string{
	"devops",
	"site reliability",
	"sre",
	"platform engineer",
	"infrastructure engineer",
	"cloud engineer",
}

var kwSubCategoryMobile = []string{
	"mobile developer",
	"mobile engineer",
	"android developer",
	"android engineer",
	"ios developer",
	"ios engineer",
	"flutter developer",
	"react native",
}

var kwSubCategoryQA = []string{
	"quality assurance",
	"qa engineer",
	"qa developer",
	"test engineer",
	"software tester",
	"test automation",
	"automation engineer",
}

var kwSubCategorySecurity = []string{
	"security engineer",
	"cybersecurity",
	"cyber security",
	"application security",
	"appsec",
	"security analyst",
	"penetration tester",
	"pentester",
}

// Data

var kwSubCategoryDataScience = []string{
	"data scientist",
	"data science",
	"statistical modeling",
	"predictive modeling",
}

var kwSubCategoryDataAnalysis = []string{
	"data analyst",
	"data analytics",
	"business intelligence analyst",
	"bi analyst",
	"analytics specialist",
}

var kwSubCategoryDataEngineering = []string{
	"data engineer",
	"data engineering",
	"etl engineer",
	"data pipeline",
	"data platform engineer",
}

var kwSubCategoryDataDBA = []string{
	"database administrator",
	"database administration",
	"dba",
	"database engineer",
}

// AI

var kwSubCategoryAIEngineering = []string{
	"ai engineer",
	"ai engineering",
	"artificial intelligence engineer",
	"applied ai engineer",
	"generative ai engineer",
	"genai engineer",
	"llm engineer",
}

var kwSubCategoryMLEngineering = []string{
	"machine learning engineer",
	"ml engineer",
	"machine learning engineering",
	"mlops engineer",
}

var kwSubCategoryComputerVision = []string{
	"computer vision",
	"computer vision engineer",
	"vision engineer",
	"image recognition",
	"image processing",
}

var kwSubCategoryAIResearch = []string{
	"ai researcher",
	"ai research",
	"machine learning researcher",
	"ml researcher",
	"research scientist",
	"research engineer",
}

// Marketing

var kwSubCategoryMarketAnalysis = []string{
	"market analyst",
	"market analysis",
	"market research",
	"marketing analyst",
	"marketing analytics",
}

var kwSubCategorySEO = []string{
	"seo",
	"search engine optimization",
	"seo specialist",
	"seo manager",
}

var kwSubCategoryCopywriting = []string{
	"copywriter",
	"copywriting",
	"content writer",
	"content writing",
}

var kwSubCategorySocialMedia = []string{
	"social media",
	"social media manager",
	"social media specialist",
	"community manager",
}

var kwSubCategoryBrandManagement = []string{
	"brand manager",
	"brand management",
	"brand strategist",
	"branding manager",
}

var kwSubCategoryMarketingOperations = []string{
	"marketing operations",
	"marketing ops",
	"marketing automation",
	"marketing operations manager",
}

// Management

var kwSubCategorySoftwareManagement = []string{
	"engineering manager",
	"software engineering manager",
	"development manager",
	"software development manager",
	"technical manager",
}

var kwSubCategoryDataManagement = []string{
	"data manager",
	"data management",
	"head of data",
	"director of data",
}

var kwSubCategoryProjectManagement = []string{
	"project manager",
	"project management",
	"technical project manager",
	"program manager",
	"producer",
}

var kwSubCategoryProductManagement = []string{
	"product manager",
	"product management",
	"technical product manager",
	"product owner",
}

var kwSubCategoryCSuite = []string{
	"chief executive officer",
	"ceo",
	"chief technology officer",
	"cto",
	"chief operating officer",
	"coo",
	"chief product officer",
	"cpo",
	"chief marketing officer",
	"cmo",
}

var kwSubCategoryOperationsManagement = []string{
	"operations manager",
	"operations management",
	"business operations manager",
	"head of operations",
}

var kwSubCategoryScrumMaster = []string{
	"scrum master",
	"agile coach",
}

var kwSubCategorySalesManagement = []string{
	"sales manager",
	"sales management",
	"sales director",
	"head of sales",
	"regional sales manager",
}

// Design

func containsAny(text string, keywords []string) bool {

	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

func retrieveCategoryFromText(text string) model.Category {
	text = strings.ToLower(text)

	switch {
	case containsAny(text, kwCategoryAI):
		return model.CategoryAI
	case containsAny(text, kwCategorySoftware):
		return model.CategorySoftwareEngineering
	case containsAny(text, kwCategoryData):
		return model.CategoryData
	case containsAny(text, kwCategoryMarketing):
		return model.CategoryMarketing
	case containsAny(text, kwCategoryManagement):
		return model.CategoryManagement
	default:
		return model.CategoryOther
	}

}

func RetrieveCategory(title string, description string) model.Category {
	retrievedCategory := retrieveCategoryFromText(title)

	if retrievedCategory == model.CategoryOther {
		retrievedCategory = retrieveCategoryFromText(description)
	}

	return retrievedCategory
}

func retrieveAISubCategory(text string) model.SubCategory {

	switch {
	case containsAny(text, kwSubCategoryAIEngineering):
		return model.SubCategoryAIEngineering
	case containsAny(text, kwSubCategoryAIResearch):
		return model.SubCategoryAIResearch
	case containsAny(text, kwSubCategoryComputerVision):
		return model.SubCategoryAIComputerVision
	case containsAny(text, kwSubCategoryMLEngineering):
		return model.SubCategoryAIMLEngineering
	default:
		return model.SubCategoryOther
	}

}

func retrieveDataSubCategory(text string) model.SubCategory {
	switch {
	case containsAny(text, kwSubCategoryDataAnalysis):
		return model.SubCategoryDataAnalysis
	case containsAny(text, kwSubCategoryDataDBA):
		return model.SubCategoryDataDBA
	case containsAny(text, kwSubCategoryDataEngineering):
		return model.SubCategoryDataEngineering
	case containsAny(text, kwSubCategoryDataScience):
		return model.SubCategoryDataScience
	case containsAny(text, kwSubCategoryDataManagement):
		return model.SubCategoryDataManagement
	default:
		return model.SubCategoryOther
	}
}

func retrieveSoftwareSubCategory(text string) model.SubCategory {
	switch {
	case containsAny(text, kwSubCategoryFullstack):
		return model.SubCategoryFullstack
	case containsAny(text, kwSubCategoryQA):
		return model.SubCategoryQA
	case containsAny(text, kwSubCategoryBackend):
		return model.SubCategoryBackend
	case containsAny(text, kwSubCategoryFrontend):
		return model.SubCategoryFrontend
	case containsAny(text, kwSubCategoryDevOps):
		return model.SubCategoryDevOps
	case containsAny(text, kwSubCategorySecurity):
		return model.SubCategorySecurity
	case containsAny(text, kwSubCategoryMobile):
		return model.SubCategoryMobile
	default:
		return model.SubCategoryOther
	}
}

func retrieveMarketingSubCategory(text string) model.SubCategory {
	switch {
	case containsAny(text, kwSubCategoryMarketAnalysis):
		return model.SubCategoryMarketAnalysis
	case containsAny(text, kwSubCategorySEO):
		return model.SubCategorySEO
	case containsAny(text, kwSubCategoryCopywriting):
		return model.SubCategoryCopywriting
	case containsAny(text, kwSubCategorySocialMedia):
		return model.SubCategorySocialMedia
	case containsAny(text, kwSubCategoryBrandManagement):
		return model.SubCategoryBrandManagement
	case containsAny(text, kwSubCategoryMarketingOperations):
		return model.SubCategoryMarketingOperations
	default:
		return model.SubCategoryOther
	}
}

func retrieveManagementSubCategory(text string) model.SubCategory {
	switch {
	case containsAny(text, kwSubCategorySoftwareManagement):
		return model.SubCategorySoftwareManagement
	case containsAny(text, kwSubCategoryDataManagement):
		return model.SubCategoryDataManagement
	case containsAny(text, kwSubCategoryProjectManagement):
		return model.SubCategoryProjectManagement
	case containsAny(text, kwSubCategoryProductManagement):
		return model.SubCategoryProductManagement
	case containsAny(text, kwSubCategoryCSuite):
		return model.SubCategoryCSuite
	case containsAny(text, kwSubCategoryOperationsManagement):
		return model.SubCategoryOperationsManagement
	case containsAny(text, kwSubCategoryScrumMaster):
		return model.SubCategoryScrumMaster
	case containsAny(text, kwSubCategorySalesManagement):
		return model.SubCategorySalesManagement
	default:
		return model.SubCategoryOther
	}

}

func retrieveSubCategoryFromText(category model.Category, text string) model.SubCategory {
	text = strings.ToLower(text)

	switch category {
	case model.CategoryAI:
		return retrieveAISubCategory(text)
	case model.CategoryData:
		return retrieveDataSubCategory(text)
	case model.CategorySoftwareEngineering:
		return retrieveSoftwareSubCategory(text)
	case model.CategoryMarketing:
		return retrieveMarketingSubCategory(text)
	case model.CategoryManagement:
		return retrieveManagementSubCategory(text)
	case model.CategoryOther:
		// TODO: handle this better -> mby return the title of job posting
		return model.SubCategoryOther
	default:
		return model.SubCategoryOther
	}
}

func RetrieveSubCategory(category model.Category, title string, description string) model.SubCategory {

	retrievedSubCategory := retrieveSubCategoryFromText(category, title)

	if retrievedSubCategory == model.SubCategoryOther {
		retrievedSubCategory = retrieveSubCategoryFromText(category, description)
	}

	return retrievedSubCategory
}

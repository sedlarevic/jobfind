package model

type Category string

const (
	CategoryOther               Category = "other"
	CategorySoftwareEngineering Category = "software_engineering"
	CategoryData                Category = "data"
	CategoryAI                  Category = "ai"
	CategoryMarketing           Category = "marketing"
	CategoryManagement          Category = "management"
	CategoryDesign              Category = "design"
)

type SubCategory string

const (
	SubCategoryOther SubCategory = "other"

	// Software
	SubCategoryBackend   SubCategory = "backend"
	SubCategoryFrontend  SubCategory = "frontend"
	SubCategoryDevOps    SubCategory = "devops"
	SubCategoryMobile    SubCategory = "mobile"
	SubCategoryQA        SubCategory = "qa"
	SubCategoryFullstack SubCategory = "fullstack"
	SubCategorySecurity  SubCategory = "security"

	// Data
	SubCategoryDataScience     SubCategory = "data_science"
	SubCategoryDataAnalysis    SubCategory = "data_analysis"
	SubCategoryDataEngineering SubCategory = "data_engineering"
	SubCategoryDataDBA         SubCategory = "database_administration"

	// AI
	SubCategoryAIEngineering    SubCategory = "ai_engineering"
	SubCategoryAIMLEngineering  SubCategory = "ml_engineering"
	SubCategoryAIComputerVision SubCategory = "computer_vision"
	SubCategoryAIResearch       SubCategory = "ai_research"

	// Marketing
	SubCategoryMarketAnalysis      SubCategory = "market_analysis"
	SubCategorySEO                 SubCategory = "seo"
	SubCategoryCopywriting         SubCategory = "copywriting"
	SubCategorySocialMedia         SubCategory = "social_media_management"
	SubCategoryBrandManagement     SubCategory = "brand_management"
	SubCategoryMarketingOperations SubCategory = "marketing_operations"

	// Management
	SubCategorySoftwareManagement   SubCategory = "software_management"
	SubCategoryDataManagement       SubCategory = "data_management"
	SubCategoryProjectManagement    SubCategory = "project_management"
	SubCategoryProductManagement    SubCategory = "product_management"
	SubCategoryCSuite               SubCategory = "c_suite"
	SubCategoryOperationsManagement SubCategory = "operations_management"
	SubCategoryScrumMaster          SubCategory = "scrum_master"
	SubCategorySalesManagement      SubCategory = "sales_management"

	// Design
	SubCategoryUXDesign      SubCategory = "ux_design"
	SubCategoryUIDesign      SubCategory = "ui_design"
	SubCategoryProductDesign SubCategory = "product_design"
	SubCategoryGraphicDesign SubCategory = "graphic_design"
	SubCategoryUXResearch    SubCategory = "ux_research"
)

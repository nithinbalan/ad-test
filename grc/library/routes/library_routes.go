package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	libraryRoutes := router.Group("/")

	libraryRoutes.Use(middlewares.AuthMiddleware(app.DB))

	questionLibraryController := app.LibraryModule.QuestionLibraryController
	KRILibraryController := app.LibraryModule.KriLibraryController
	RiskLibraryController := app.LibraryModule.RiskLibraryController
	KnowledgeCenterResourceController := app.LibraryModule.KnowledgeCenterResourceController
	ComplianceLibraryController := app.LibraryModule.ComplianceLibraryController
	KpiLibraryController := app.LibraryModule.KpiLibraryController
	StrategicObjectiveLibraryController := app.LibraryModule.StrategicObjectiveLibraryController
	VulnerabilityLibraryController := app.LibraryModule.VulnerabilityLibraryController
	ArticleLibraryController := app.LibraryModule.ArticleLibraryController

	libraryRoutes.GET("/question-library", questionLibraryController.GetAllQuestionLibrary)
	libraryRoutes.POST("/question-library", questionLibraryController.CreateQuestionLibrary)
	libraryRoutes.GET("/question-library/:id", questionLibraryController.GetQuestionLibraryDetails)
	libraryRoutes.PUT("/question-library/:id", questionLibraryController.UpdateQuestionLibrary)
	libraryRoutes.POST("/question-library/delete", questionLibraryController.DeleteQuestionLibraries)
	libraryRoutes.PUT("/question-library/activate", questionLibraryController.ActivateQuestionLibraries)
	libraryRoutes.PUT("/question-library/deactivate", questionLibraryController.DeactivateQuestionLibraries)

	libraryRoutes.GET("/kri-library/export", KRILibraryController.GetKRILibraryExport)
	libraryRoutes.GET("/kri-library", KRILibraryController.GetAllKRILibrary)
	libraryRoutes.GET("/kri-library/category-count", KRILibraryController.KRICountByCategory)
	libraryRoutes.POST("/kri-library", KRILibraryController.CreateKRILibrary)
	libraryRoutes.GET("/kri-library/:id", KRILibraryController.GetKRILibraryDetails)
	libraryRoutes.PUT("/kri-library/:id", KRILibraryController.UpdateKRILibrary)
	libraryRoutes.POST("/kri-library/delete", KRILibraryController.DeleteKRILibraries)
	libraryRoutes.PUT("/kri-library/activate", KRILibraryController.ActivateKRILibraries)
	libraryRoutes.PUT("/kri-library/deactivate", KRILibraryController.DeactivateKRILibraries)

	libraryRoutes.GET("/risk-library/export", RiskLibraryController.GetRiskLibraryExport)
	libraryRoutes.GET("/risk-library", RiskLibraryController.GetAllRiskLibrary)
	libraryRoutes.GET("/risk-library/category-count", RiskLibraryController.RiskCountByCategory)
	libraryRoutes.GET("/risk-library/likelihood-count", RiskLibraryController.RiskCountByLikelihood)
	libraryRoutes.GET("/risk-library/impact-count", RiskLibraryController.RiskCountByImpact)

	libraryRoutes.GET("/knowledge-center", KnowledgeCenterResourceController.GetAllResources)

	libraryRoutes.GET("/compliance-library/nca", ComplianceLibraryController.GetNCADetails)
	libraryRoutes.GET("/compliance-library/sama", ComplianceLibraryController.GetSAMADetails)

	libraryRoutes.GET("/compliance-library/:complianceLibraryID/versions/:versionID/clauses", ComplianceLibraryController.GetClauseControls)

	// libraryRoutes.GET("/control-mapping", ComplianceLibraryController.GetAllMappings)
	libraryRoutes.GET("/control-mapping", ComplianceLibraryController.GetControlMappings)
	libraryRoutes.GET("/control-mapping/export", ComplianceLibraryController.GetControlMappingsExport)

	libraryRoutes.GET("/kpi-library/export", KpiLibraryController.GetKpiExport)
	libraryRoutes.GET("/kpi-library", KpiLibraryController.GetAllKpi)
	libraryRoutes.GET("/kpi-library/category-count", KpiLibraryController.KPICountByCategory)

	libraryRoutes.GET("/strategic-objective/export", StrategicObjectiveLibraryController.Export)
	libraryRoutes.GET("/strategic-objective", StrategicObjectiveLibraryController.GetAll)
	libraryRoutes.GET("/strategic-objective/category-count", StrategicObjectiveLibraryController.CountByCategory)

	libraryRoutes.GET("/vulnerability-library/export", VulnerabilityLibraryController.Export)
	libraryRoutes.GET("/vulnerability-library", VulnerabilityLibraryController.GetAll)
	libraryRoutes.GET("/vulnerability-library/category-count", VulnerabilityLibraryController.CountByCategory)

	libraryRoutes.GET("/article-library", ArticleLibraryController.GetAll)
	libraryRoutes.GET("/article-library/:id", ArticleLibraryController.GetArticleLibraryDetails)
	libraryRoutes.GET("/article-categories", ArticleLibraryController.ArticleCategories)
}

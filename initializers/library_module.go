package initializers

import (
	"jethur/grc/library/controllers"
	"jethur/grc/library/repositories"
	"jethur/grc/library/services"

	"gorm.io/gorm"
)

type LibraryModule struct {
	QuestionLibraryController         *controllers.QuestionLibraryController
	KriLibraryController              *controllers.KRILibraryController
	RiskLibraryController             *controllers.RiskLibraryController
	KnowledgeCenterResourceController *controllers.KnowledgeCenterResourceController
	ComplianceLibraryController       *controllers.ComplianceLibraryController

	KriLibrarySvc services.KRILibraryServiceInterface
}

func NewLibraryModule(db *gorm.DB, sharedModule *SharedModule, MasterModule *MasterModule) *LibraryModule {
	// Question Library
	questionLibraryRepo := repositories.NewQuestionLibraryRepository(db)
	questionLibrarySvc := services.NewQuestionLibraryService(questionLibraryRepo, sharedModule.ActivityLogService, MasterModule.QuestionGroupService)
	questionLibraryController := controllers.NewQuestionLibraryController(questionLibrarySvc)

	// KRI Library
	kriLibraryRepository := repositories.NewKRILibraryRepository(db)
	kriLibrarySvc := services.NewKRILibraryService(kriLibraryRepository, MasterModule.ReviewFrequencySvc, sharedModule.ActivityLogService)
	kriLibraryController := controllers.NewKRILibraryController(kriLibrarySvc)

	RiskLibraryRepository := repositories.NewRiskLibraryRepository(db)
	RiskLibrarySvc := services.NewRiskLibraryService(RiskLibraryRepository)
	RiskLibraryController := controllers.NewRiskLibraryController(RiskLibrarySvc)

	knowledgeCenterResourceRepository := repositories.NewKnowledgeCenterResourceRepository(db)
	knowledgeCenterResourceSvc := services.NewKnowledgeCenterResourceService(knowledgeCenterResourceRepository)
	knowledgeCenterResourceController := controllers.NewKnowledgeCenterResourceController(knowledgeCenterResourceSvc, db)

	ComplianceLibraryRepository := repositories.NewComplianceLibraryRepository(db)
	ComplianceLibrarySvc := services.NewComplianceLibraryService(ComplianceLibraryRepository)
	ComplianceLibraryController := controllers.NewComplianceLibraryController(ComplianceLibrarySvc, db)

	return &LibraryModule{
		QuestionLibraryController: questionLibraryController,
		KriLibraryController:      kriLibraryController,

		KriLibrarySvc:                     kriLibrarySvc,
		RiskLibraryController:             RiskLibraryController,
		KnowledgeCenterResourceController: knowledgeCenterResourceController,
		ComplianceLibraryController:       ComplianceLibraryController,
	}
}

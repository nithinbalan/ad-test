package initializers

import (
	"jethur/grc/library/controllers"
	"jethur/grc/library/repositories"
	"jethur/grc/library/services"

	"gorm.io/gorm"
)

type LibraryModule struct {
	QuestionLibraryController *controllers.QuestionLibraryController
	KriLibraryController      *controllers.KRILibraryController
	RiskLibraryController     *controllers.RiskLibraryController

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

	return &LibraryModule{
		QuestionLibraryController: questionLibraryController,
		KriLibraryController:      kriLibraryController,

		KriLibrarySvc:         kriLibrarySvc,
		RiskLibraryController: RiskLibraryController,
	}
}

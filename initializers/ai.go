package initializers

import (
	"jethur/grc/jethur-ai/controllers"
	"jethur/grc/jethur-ai/services"

	"gorm.io/gorm"
)

type JethurAIModule struct {
	DocumentAIController *controllers.DocumentAIController

	DocumentAIService services.DocumentAIServiceInterface
}

func NewJethurAIModule(db *gorm.DB, documentModule *DocumentModule) *JethurAIModule {
	// documentAIRepo := repositories.NewAIRepository(db)
	documentAISvc := services.NewDocumentAIService(documentModule.DocumentSvc)
	documentAIController := controllers.NewDocumentAIController(documentAISvc)

	return &JethurAIModule{
		DocumentAIController: documentAIController,

		DocumentAIService: documentAISvc,
	}
}

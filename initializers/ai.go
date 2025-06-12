package initializers

import (
	"jethur/grc/jethur-ai/controllers"
	"jethur/grc/jethur-ai/services"

	"gorm.io/gorm"
)

type JethurAIModule struct {
	DocumentAIController *controllers.DocumentAIController
	RiskAIController     *controllers.RiskAIController

	DocumentAIService services.DocumentAIServiceInterface
}

func NewJethurAIModule(db *gorm.DB, documentModule *DocumentModule, organizationModule *OrganizationModule) *JethurAIModule {
	// documentAIRepo := repositories.NewAIRepository(db)
	documentAISvc := services.NewDocumentAIService(documentModule.DocumentSvc, organizationModule.UserSvc)
	documentAIController := controllers.NewDocumentAIController(documentAISvc)

	riskAISvc := services.NewRiskAIService()
	riskAIController := controllers.NewRiskAIController(riskAISvc)

	return &JethurAIModule{
		DocumentAIController: documentAIController,
		RiskAIController:     riskAIController,

		DocumentAIService: documentAISvc,
	}
}

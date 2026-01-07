package initializers

import (
	"jethur/grc/jethur-ai/controllers"
	"jethur/grc/jethur-ai/services"
     
	"gorm.io/gorm"
)

type JethurAIModule struct {
	DocumentAIController   *controllers.DocumentAIController
	CommonAIController     *controllers.CommonAIController
	RiskAIController       *controllers.RiskAIController
	ComplianceAIController *controllers.ComplianceAIController

	DocumentAIService services.DocumentAIServiceInterface
}

func NewJethurAIModule(db *gorm.DB, documentModule *DocumentModule, organizationModule *OrganizationModule, settingsModule *SettingsModule, complianceModule *ComplianceModule,) *JethurAIModule {
	// documentAIRepo := repositories.NewAIRepository(db)
	documentAISvc := services.NewDocumentAIService(documentModule.DocumentSvc, organizationModule.UserSvc, settingsModule.SystemPreferenceSvc)
	documentAIController := controllers.NewDocumentAIController(documentAISvc)

	commonAISvc := services.NewCommonAIService(settingsModule.SystemPreferenceSvc)
	commonAIController := controllers.NewCommonAIController(commonAISvc)

	riskAISvc := services.NewRiskAIService(settingsModule.SystemPreferenceSvc)
	riskAIController := controllers.NewRiskAIController(riskAISvc)

	complianceAISvc := services.NewComplianceAIService(settingsModule.SystemPreferenceSvc,complianceModule.ComplianceRegisterSvc)
	complianceAIController := controllers.NewComplianceAIController(complianceAISvc)

	return &JethurAIModule{
		DocumentAIController:   documentAIController,
		RiskAIController:       riskAIController,
		CommonAIController:     commonAIController,
		ComplianceAIController: complianceAIController,

		DocumentAIService: documentAISvc,
	}
}

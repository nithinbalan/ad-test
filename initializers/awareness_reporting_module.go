package initializers

import (
	"jethur/grc/awareness-reporting/controllers"
	"jethur/grc/awareness-reporting/repositories"
	"jethur/grc/awareness-reporting/services"

	"gorm.io/gorm"
)

type AwarenessReportingModule struct {
	AwarenessCampaignTypeController *controllers.AwarenessCampaignTypeController
}

func NewAwarenessReportingModule(
	db *gorm.DB,
	sharedModule *SharedModule,
) *AwarenessReportingModule {

	awarenessCampaignTypeRepo := repositories.NewAwarenessCampaignTypeRepository(db)
	awarenessCampaignTypeSvc := services.NewAwarenessCampaignTypeService(awarenessCampaignTypeRepo,sharedModule.ActivityLogService,)
	awarenessCampaignTypeController := controllers.NewAwarenessCampaignTypeController(awarenessCampaignTypeSvc,)

	return &AwarenessReportingModule{
		AwarenessCampaignTypeController: awarenessCampaignTypeController,
	}
}

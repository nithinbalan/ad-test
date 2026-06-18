package initializers

import (
	"jethur/grc/awareness-reporting/controllers"
	"jethur/grc/awareness-reporting/repositories"
	"jethur/grc/awareness-reporting/services"

	"gorm.io/gorm"
)

type AwarenessReportingModule struct {
	AwarenessCampaignTypeController   *controllers.AwarenessCampaignTypeController
	AwarenessChannelController        *controllers.AwarenessChannelController
	AwarenessCampaignStatusController *controllers.AwarenessCampaignStatusController
}

func NewAwarenessReportingModule(
	db *gorm.DB,
	sharedModule *SharedModule,
) *AwarenessReportingModule {

	awarenessCampaignTypeRepo := repositories.NewAwarenessCampaignTypeRepository(db)
	awarenessCampaignTypeSvc := services.NewAwarenessCampaignTypeService(awarenessCampaignTypeRepo, sharedModule.ActivityLogService)
	awarenessCampaignTypeController := controllers.NewAwarenessCampaignTypeController(awarenessCampaignTypeSvc)

	awarenessChannelRepo := repositories.NewAwarenessChannelRepository(db)
	awarenessChannelSvc := services.NewAwarenessChannelService(awarenessChannelRepo, sharedModule.ActivityLogService)
	awarenessChannelController := controllers.NewAwarenessChannelController(awarenessChannelSvc)

	awarenessCampaignStatusRepo := repositories.NewAwarenessCampaignStatusRepository(db)
	awarenessCampaignStatusSvc := services.NewAwarenessCampaignStatusService(awarenessCampaignStatusRepo)
	awarenessCampaignStatusController := controllers.NewAwarenessCampaignStatusController(awarenessCampaignStatusSvc)

	return &AwarenessReportingModule{
		AwarenessCampaignTypeController:   awarenessCampaignTypeController,
		AwarenessChannelController:        awarenessChannelController,
		AwarenessCampaignStatusController: awarenessCampaignStatusController,
	}
}

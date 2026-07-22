package initializers

import (
	"jethur/grc/awareness-reporting/controllers"
	"jethur/grc/awareness-reporting/repositories"
	"jethur/grc/awareness-reporting/services"

	"gorm.io/gorm"
)

type AwarenessReportingModule struct {
	AwarenessCampaignTypeController       *controllers.AwarenessCampaignTypeController
	AwarenessChannelController            *controllers.AwarenessChannelController
	AwarenessCampaignStatusController     *controllers.AwarenessCampaignStatusController
	AwarenessPriorityController           *controllers.AwarenessPriorityController
	AwarenessEvaluationRatingController   *controllers.AwarenessEvaluationRatingController
	AwarenessCampaignController           *controllers.AwarenessCampaignController
	AwarenessMaturityLevelController      *controllers.AwarenessMaturityLevelController
	AwarenessAssessmentThemeController    *controllers.AwarenessAssessmentThemeController
	AwarenessPostCampaignSurveyController *controllers.AwarenessPostCampaignSurveyController
	AwarenessHandlingStatusController     *controllers.AwarenessHandlingStatusController
	AwarenessReportSubCategoryController  *controllers.AwarenessReportSubCategoryController
	AwarenessReportController             *controllers.AwarenessReportController
}

func NewAwarenessReportingModule(
	db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	organizationModule *OrganizationModule,
	commonModule *CommonModule,
	masterModule *MasterModule,
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

	awarenessPriorityRepo := repositories.NewAwarenessPriorityRepository(db)
	awarenessPrioritySvc := services.NewAwarenessPriorityService(awarenessPriorityRepo)
	awarenessPriorityController := controllers.NewAwarenessPriorityController(awarenessPrioritySvc)

	awarenessEvaluationRatingRepo := repositories.NewAwarenessEvaluationRatingRepository(db)
	awarenessEvaluationRatingSvc := services.NewAwarenessEvaluationRatingService(awarenessEvaluationRatingRepo)
	awarenessEvaluationRatingController := controllers.NewAwarenessEvaluationRatingController(awarenessEvaluationRatingSvc)

	awarenessCampaignRepo := repositories.NewAwarenessCampaignRepository(db)
	awarenessCampaignSvc := services.NewAwarenessCampaignService(
		awarenessCampaignRepo,
		sharedModule.ActivityLogService,
		awarenessCampaignTypeSvc,
		awarenessPrioritySvc,
		awarenessCampaignStatusSvc,
		awarenessEvaluationRatingSvc,
		settingsModule.OrganizationUnitSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		organizationModule.UserSettingSvc,
	)
	awarenessCampaignController := controllers.NewAwarenessCampaignController(awarenessCampaignSvc)

	awarenessMaturityLevelRepo := repositories.NewAwarenessMaturityLevelRepository(db)
	awarenessMaturityLevelSvc := services.NewAwarenessMaturityLevelService(awarenessMaturityLevelRepo)
	awarenessMaturityLevelController := controllers.NewAwarenessMaturityLevelController(awarenessMaturityLevelSvc)

	awarenessAssessmentThemeRepo := repositories.NewAwarenessAssessmentThemeRepository(db)
	awarenessAssessmentThemeSvc := services.NewAwarenessAssessmentThemeService(awarenessAssessmentThemeRepo, sharedModule.ActivityLogService)
	awarenessAssessmentThemeController := controllers.NewAwarenessAssessmentThemeController(awarenessAssessmentThemeSvc)

	awarenessPostCampaignSurveyRepo := repositories.NewAwarenessPostCampaignSurveyRepository(db)
	awarenessPostCampaignSurveySvc := services.NewAwarenessPostCampaignSurveyService(
		awarenessPostCampaignSurveyRepo,
		sharedModule.ActivityLogService,
		awarenessMaturityLevelSvc,
		awarenessCampaignSvc,
		awarenessAssessmentThemeSvc,
		organizationModule.UserSettingSvc,
	)
	awarenessPostCampaignSurveyController := controllers.NewAwarenessPostCampaignSurveyController(awarenessPostCampaignSurveySvc)

	awarenessHandlingStatusRepo := repositories.NewAwarenessHandlingStatusRepository(db)
	awarenessHandlingStatusSvc := services.NewAwarenessHandlingStatusService(awarenessHandlingStatusRepo)
	awarenessHandlingStatusController := controllers.NewAwarenessHandlingStatusController(awarenessHandlingStatusSvc)

	awarenessReportSubCategoryRepo := repositories.NewAwarenessReportSubCategoryRepository(db)
	awarenessReportSubCategorySvc := services.NewAwarenessReportSubCategoryService(
		awarenessReportSubCategoryRepo,
		sharedModule.ActivityLogService,
	)
	awarenessReportSubCategoryController := controllers.NewAwarenessReportSubCategoryController(
		awarenessReportSubCategorySvc,
	)

	awarenessReportRepo := repositories.NewAwarenessReportRepository(db)
	awarenessReportSvc := services.NewAwarenessReportService(
		awarenessReportRepo,
		sharedModule.ActivityLogService,
		awarenessCampaignTypeSvc,
		awarenessReportSubCategorySvc,
		masterModule.SeverityLevelSvc,
		awarenessHandlingStatusSvc,
		organizationModule.UserSvc,
		organizationModule.UserSettingSvc,
	)
	awarenessReportController := controllers.NewAwarenessReportController(awarenessReportSvc)

	return &AwarenessReportingModule{
		AwarenessCampaignTypeController:       awarenessCampaignTypeController,
		AwarenessChannelController:            awarenessChannelController,
		AwarenessCampaignStatusController:     awarenessCampaignStatusController,
		AwarenessPriorityController:           awarenessPriorityController,
		AwarenessEvaluationRatingController:   awarenessEvaluationRatingController,
		AwarenessCampaignController:           awarenessCampaignController,
		AwarenessMaturityLevelController:      awarenessMaturityLevelController,
		AwarenessAssessmentThemeController:    awarenessAssessmentThemeController,
		AwarenessPostCampaignSurveyController: awarenessPostCampaignSurveyController,
		AwarenessHandlingStatusController:     awarenessHandlingStatusController,
		AwarenessReportSubCategoryController:  awarenessReportSubCategoryController,
		AwarenessReportController:             awarenessReportController,
	}
}

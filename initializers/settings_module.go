package initializers

import (
	orgRepo "jethur/grc/organization/repositories"
	"jethur/grc/settings/controllers"
	"jethur/grc/settings/repositories"
	"jethur/grc/settings/services"

	"gorm.io/gorm"
)

type SettingsModule struct {
	OrganizationLevelController      *controllers.OrganizationLevelController
	SystemPreferenceController       *controllers.SystemPreferenceController
	RoleController                   *controllers.RoleController
	OrganizationUnitController       *controllers.OrganizationUnitController
	LogSettingsController            *controllers.LogSettingsController
	WorkflowLevelTypeController      *controllers.WorkflowLevelTypeController
	WorkflowController               *controllers.WorkflowController
	RolePermissionTypeController     *controllers.RolePermissionTypeController
	ModuleController                 *controllers.ModuleController
	EmailTemplateSettingsController  *controllers.EmailTemplateSettingsController
	IntegrationController            *controllers.IntegrationController
	AuditTemplateSettingController   *controllers.AuditTemplateSettingController
	RiskTemplateSettingController    *controllers.RiskTemplateSettingController
	ReportTemplateSettingController  *controllers.ReportTemplateSettingController
	ActiveDirectorySettingController *controllers.ActiveDirectorySettingController
	SIEMSettingController            *controllers.SIEMSettingController
	AISettingController              *controllers.AISettingController
	TriggerConditionController       *controllers.TriggerConditionController
	ReminderController               *controllers.ReminderController
	ReminderRuleController           *controllers.ReminderRuleController
	EscalationController             *controllers.EscalationController
	EscalationLevelController        *controllers.EscalationLevelController
	ReminderContextController        *controllers.ReminderContextController
	RAGSettingController             *controllers.RAGSettingController
	TenableSettingsController        *controllers.TenableSettingsController
	HealthCheckController            *controllers.HealthCheckController
	ThemeController                  *controllers.ThemeController

	SystemPreferenceSvc       services.SystemPreferenceServiceInterface
	OrganizationUnitSvc       services.OrganizationUnitServiceInterface
	OrganizationLevelSvc      services.OrganizationLevelServiceInterface
	RoleSvc                   services.RoleServiceInterface
	ModuleSvc                 services.ModuleServiceInterface
	WorkflowSvc               services.WorkflowServiceInterface
	RolePermissionTypeSvc     services.RolePermissionTypeServiceInterface
	EmailTemplateSettingsSvc  services.EmailTemplateSettingsServiceInterface
	IntegrationSvc            services.IntegrationServiceInterface
	AuditTemplateSettingSvc   services.AuditTemplateSettingServiceInterface
	RiskTemplateSettingSvc    services.RiskTemplateSettingServiceInterface
	ReportTemplateSettingSvc  services.ReportTemplateSettingServiceInterface
	ActiveDirectorySettingSvc services.ActiveDirectorySettingServiceInterface
	ThemeSvc                  services.ThemeServiceInterface
	TenableSettingsSvc        services.TenableSettingsServiceInterface
}

func NewSettingsModule(db *gorm.DB, sharedModule *SharedModule, masterModule *MasterModule) *SettingsModule {
	organizationLevelRepo := repositories.NewOrganizationLevelRepository(db)
	organizationLevelSvc := services.NewOrganizationLevelService(organizationLevelRepo, sharedModule.ActivityLogService)
	organizationLevelController := controllers.NewOrganizationLevelController(organizationLevelSvc)

	integrationRepo := repositories.NewIntegrationRepository(db)
	integrationSvc := services.NewIntegrationService(integrationRepo, sharedModule.ActivityLogService, masterModule.IntegrationCategorySvc)
	integrationController := controllers.NewIntegrationController(integrationSvc)

	moduleRepo := repositories.NewModuleRepository(db)
	moduleSvc := services.NewModuleService(moduleRepo)
	moduleController := controllers.NewModuleController(moduleSvc)

	systemPreferenceRepo := repositories.NewSystemPreferenceRepository(db)
	systemPreferenceSvc := services.NewSystemPreferenceService(systemPreferenceRepo, sharedModule.ActivityLogService, masterModule.CurrencySvc, integrationSvc, moduleSvc)
	systemPreferenceController := controllers.NewSystemPreferenceController(systemPreferenceSvc)

	emailTemplateSettingsRepo := repositories.NewEmailTemplateSettingsRepository(db)
	emailTemplateSettingsService := services.NewEmailTemplateSettingsService(emailTemplateSettingsRepo, sharedModule.ActivityLogService)
	emailTemplateSettingsController := controllers.NewEmailTemplateSettingsController(emailTemplateSettingsService)

	rolePermissionTypeRepo := repositories.NewRolePermissionTypeRepository(db)
	rolePermissionTypeSvc := services.NewRolePermissionTypeService(rolePermissionTypeRepo)
	rolePermissionTypeController := controllers.NewRolePermissionTypeController(rolePermissionTypeSvc)

	roleRepo := repositories.NewRoleRepository(db)
	roleSvc := services.NewRoleService(roleRepo, sharedModule.ActivityLogService, moduleSvc, rolePermissionTypeRepo, masterModule.LicenseTypeSvc)
	roleController := controllers.NewRoleController(roleSvc)

	organizationUnitRepo := repositories.NewOrganizationUnitRepository(db)
	organizationUnitSvc := services.NewOrganizationUnitService(organizationUnitRepo, sharedModule.ActivityLogService, organizationLevelSvc)
	organizationUnitController := controllers.NewOrganizationUnitController(organizationUnitSvc)

	logSettingsRepo := repositories.NewLogSettingsRepository(db)
	logSettingsSvc := services.NewLogSettingsService(logSettingsRepo, sharedModule.ActivityLogService)
	logSettingsController := controllers.NewLogSettingsController(logSettingsSvc)

	workflowLevelTypeRepo := repositories.NewWorkflowLevelTypeRepository(db)
	workflowLevelTypeSvc := services.NewWorkflowLevelTypeService(workflowLevelTypeRepo)
	workflowLevelTypeController := controllers.NewWorkflowLevelTypeController(workflowLevelTypeSvc)

	workflowRepo := repositories.NewWorkflowRepository(db)
	workflowSvc := services.NewWorkflowService(workflowRepo, sharedModule.ActivityLogService, workflowLevelTypeSvc)
	workflowController := controllers.NewWorkflowController(workflowSvc)

	auditTemplateSettingRepo := repositories.NewAuditTemplateSettingRepository(db)
	auditTemplateSettingSvc := services.NewAuditTemplateSettingService(
		auditTemplateSettingRepo,
		sharedModule.ActivityLogService)
	auditTemplateSettingController := controllers.NewAuditTemplateSettingController(auditTemplateSettingSvc)

	riskTemplateSettingRepo := repositories.NewRiskTemplateSettingRepository(db)
	riskTemplateSettingSvc := services.NewRiskTemplateSettingService(
		riskTemplateSettingRepo,
		sharedModule.ActivityLogService)
	riskTemplateSettingController := controllers.NewRiskTemplateSettingController(riskTemplateSettingSvc)

	reportTemplateSettingRepo := repositories.NewReportTemplateSettingRepository(db)
	reportTemplateSettingSvc := services.NewReportTemplateSettingService(
		reportTemplateSettingRepo,
		sharedModule.ActivityLogService)
	reportTemplateSettingController := controllers.NewReportTemplateSettingController(reportTemplateSettingSvc)

	activeDirectorySettingRepo := repositories.NewActiveDirectorySettingRepository(db)
	orgRepo := orgRepo.NewUserRepository(db)
	activeDirectorySettingSvc := services.NewActiveDirectorySettingService(
		activeDirectorySettingRepo, orgRepo, organizationLevelRepo, integrationSvc)
	activeDirectorySettingController := controllers.NewActiveDirectorySettingController(activeDirectorySettingSvc)

	siemSettingRepo := repositories.NewSIEMSettingRepository(db)
	siemSettingSvc := services.NewSIEMSettingService(
		siemSettingRepo, integrationSvc)
	siemSettingController := controllers.NewSIEMSettingController(siemSettingSvc)

	aiSettingRepo := repositories.NewAISettingRepository(db)
	aiSettingSvc := services.NewAISettingService(
		aiSettingRepo, systemPreferenceSvc)
	aiSettingController := controllers.NewAISettingController(aiSettingSvc)

	triggerConditionRepo := repositories.NewTriggerConditionRepository(db)
	triggerConditionSvc := services.NewTriggerConditionService(triggerConditionRepo)
	triggerConditionController := controllers.NewTriggerConditionController(triggerConditionSvc)

	reminderRepo := repositories.NewReminderRepository(db)
	reminderSvc := services.NewReminderService(reminderRepo)
	reminderController := controllers.NewReminderController(reminderSvc)

	reminderContextRepo := repositories.NewReminderContextRepository(db)
	reminderContextSvc := services.NewReminderContextService(reminderContextRepo)
	reminderContextController := controllers.NewReminderContextController(reminderContextSvc)

	reminderRuleRepo := repositories.NewReminderRuleRepository(db)
	reminderRuleSvc := services.NewReminderRuleService(reminderRuleRepo, reminderSvc, triggerConditionSvc, sharedModule.ActivityLogService, reminderContextSvc)
	reminderRuleController := controllers.NewReminderRuleController(reminderRuleSvc)

	escalationRepo := repositories.NewEscalationRepository(db)
	escalationSvc := services.NewEscalationService(escalationRepo)
	escalationController := controllers.NewEscalationController(escalationSvc)

	escalationLevelRepo := repositories.NewEscalationLevelRepository(db)
	escalationLevelSvc := services.NewEscalationLevelService(escalationLevelRepo, escalationSvc, sharedModule.ActivityLogService)
	escalationLevelController := controllers.NewEscalationLevelController(escalationLevelSvc)

	ragSettingRepo := repositories.NewRAGSettingRepository(db)
	ragSettingSvc := services.NewRAGSettingService(ragSettingRepo)
	ragSettingController := controllers.NewRAGSettingController(ragSettingSvc)

	tenableSettingsRepo := repositories.NewTenableSettingsRepository(db)
	tenableSettingsSvc := services.NewTenableSettingsService(tenableSettingsRepo, integrationSvc, sharedModule.ActivityLogService)
	tenableSettingsController := controllers.NewTenableSettingsController(tenableSettingsSvc)

	healthCheckRepo := repositories.NewHealthCheckRepository(db)
	healthCheckService := services.NewHealthCheckService(healthCheckRepo, systemPreferenceSvc)
	healthCheckController := controllers.NewHealthCheckController(healthCheckService)

	themeRepo := repositories.NewThemeRepository(db)
	themeSvc := services.NewThemeService(themeRepo, sharedModule.ActivityLogService)
	themeController := controllers.NewThemeController(themeSvc)

	return &SettingsModule{
		OrganizationLevelController:      organizationLevelController,
		SystemPreferenceController:       systemPreferenceController,
		RoleController:                   roleController,
		OrganizationUnitController:       organizationUnitController,
		LogSettingsController:            logSettingsController,
		WorkflowLevelTypeController:      workflowLevelTypeController,
		WorkflowController:               workflowController,
		RolePermissionTypeController:     rolePermissionTypeController,
		ModuleController:                 moduleController,
		EmailTemplateSettingsController:  emailTemplateSettingsController,
		IntegrationController:            integrationController,
		AuditTemplateSettingController:   auditTemplateSettingController,
		RiskTemplateSettingController:    riskTemplateSettingController,
		ReportTemplateSettingController:  reportTemplateSettingController,
		ActiveDirectorySettingController: activeDirectorySettingController,
		SIEMSettingController:            siemSettingController,
		AISettingController:              aiSettingController,
		TriggerConditionController:       triggerConditionController,
		ReminderController:               reminderController,
		ReminderRuleController:           reminderRuleController,
		EscalationController:             escalationController,
		EscalationLevelController:        escalationLevelController,
		ReminderContextController:        reminderContextController,
		RAGSettingController:             ragSettingController,
		TenableSettingsController:        tenableSettingsController,
		HealthCheckController:            healthCheckController,
		ThemeController:                  themeController,

		SystemPreferenceSvc:       systemPreferenceSvc,
		OrganizationLevelSvc:      organizationLevelSvc,
		RoleSvc:                   roleSvc,
		ModuleSvc:                 moduleSvc,
		WorkflowSvc:               workflowSvc,
		OrganizationUnitSvc:       organizationUnitSvc,
		RolePermissionTypeSvc:     rolePermissionTypeSvc,
		EmailTemplateSettingsSvc:  emailTemplateSettingsService,
		IntegrationSvc:            integrationSvc,
		AuditTemplateSettingSvc:   auditTemplateSettingSvc,
		RiskTemplateSettingSvc:    riskTemplateSettingSvc,
		ReportTemplateSettingSvc:  reportTemplateSettingSvc,
		ActiveDirectorySettingSvc: activeDirectorySettingSvc,
		ThemeSvc:                  themeSvc,
		TenableSettingsSvc:        tenableSettingsSvc,
	}
}

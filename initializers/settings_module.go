package initializers

import (
	"jethur/grc/settings/controllers"
	"jethur/grc/settings/repositories"
	"jethur/grc/settings/services"

	"gorm.io/gorm"
)

type SettingsModule struct {
	OrganizationLevelController     *controllers.OrganizationLevelController
	SystemPreferenceController      *controllers.SystemPreferenceController
	RoleController                  *controllers.RoleController
	OrganizationUnitController      *controllers.OrganizationUnitController
	LogSettingsController           *controllers.LogSettingsController
	WorkflowLevelTypeController     *controllers.WorkflowLevelTypeController
	WorkflowController              *controllers.WorkflowController
	RolePermissionTypeController    *controllers.RolePermissionTypeController
	ModuleController                *controllers.ModuleController
	EmailTemplateSettingsController *controllers.EmailTemplateSettingsController
	IntegrationController           *controllers.IntegrationController

	SystemPreferenceSvc      services.SystemPreferenceServiceInterface
	OrganizationUnitSvc      services.OrganizationUnitServiceInterface
	OrganizationLevelSvc     services.OrganizationLevelServiceInterface
	RoleSvc                  services.RoleServiceInterface
	ModuleSvc                services.ModuleServiceInterface
	WorkflowSvc              services.WorkflowServiceInterface
	RolePermissionTypeSvc    services.RolePermissionTypeServiceInterface
	EmailTemplateSettingsSvc services.EmailTemplateSettingsServiceInterface
	IntegrationSvc           services.IntegrationServiceInterface
}

func NewSettingsModule(db *gorm.DB, sharedModule *SharedModule, masterModule *MasterModule) *SettingsModule {
	organizationLevelRepo := repositories.NewOrganizationLevelRepository(db)
	organizationLevelSvc := services.NewOrganizationLevelService(organizationLevelRepo, sharedModule.ActivityLogService)
	organizationLevelController := controllers.NewOrganizationLevelController(organizationLevelSvc)

	systemPreferenceRepo := repositories.NewSystemPreferenceRepository(db)
	systemPreferenceSvc := services.NewSystemPreferenceService(systemPreferenceRepo, sharedModule.ActivityLogService, masterModule.CurrencySvc)
	systemPreferenceController := controllers.NewSystemPreferenceController(systemPreferenceSvc)

	emailTemplateSettingsRepo := repositories.NewEmailTemplateSettingsRepository(db)
	emailTemplateSettingsService := services.NewEmailTemplateSettingsService(emailTemplateSettingsRepo, sharedModule.ActivityLogService)
	emailTemplateSettingsController := controllers.NewEmailTemplateSettingsController(emailTemplateSettingsService)

	moduleRepo := repositories.NewModuleRepository(db)
	moduleSvc := services.NewModuleService(moduleRepo)
	moduleController := controllers.NewModuleController(moduleSvc)

	rolePermissionTypeRepo := repositories.NewRolePermissionTypeRepository(db)
	rolePermissionTypeSvc := services.NewRolePermissionTypeService(rolePermissionTypeRepo)
	rolePermissionTypeController := controllers.NewRolePermissionTypeController(rolePermissionTypeSvc)

	roleRepo := repositories.NewRoleRepository(db)
	roleSvc := services.NewRoleService(roleRepo, sharedModule.ActivityLogService, moduleSvc, rolePermissionTypeRepo)
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

	integrationRepo := repositories.NewIntegrationRepository(db)
	integrationSvc := services.NewIntegrationService(integrationRepo, sharedModule.ActivityLogService, masterModule.IntegrationCategorySvc)
	integrationController := controllers.NewIntegrationController(integrationSvc)

	return &SettingsModule{
		OrganizationLevelController:     organizationLevelController,
		SystemPreferenceController:      systemPreferenceController,
		RoleController:                  roleController,
		OrganizationUnitController:      organizationUnitController,
		LogSettingsController:           logSettingsController,
		WorkflowLevelTypeController:     workflowLevelTypeController,
		WorkflowController:              workflowController,
		RolePermissionTypeController:    rolePermissionTypeController,
		ModuleController:                moduleController,
		EmailTemplateSettingsController: emailTemplateSettingsController,
		IntegrationController:           integrationController,

		SystemPreferenceSvc:      systemPreferenceSvc,
		OrganizationLevelSvc:     organizationLevelSvc,
		RoleSvc:                  roleSvc,
		ModuleSvc:                moduleSvc,
		WorkflowSvc:              workflowSvc,
		OrganizationUnitSvc:      organizationUnitSvc,
		RolePermissionTypeSvc:    rolePermissionTypeSvc,
		EmailTemplateSettingsSvc: emailTemplateSettingsService,
		IntegrationSvc:           integrationSvc,
	}
}

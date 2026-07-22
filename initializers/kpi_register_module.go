package initializers

import (
	"jethur/grc/kpi-register/controllers"
	"jethur/grc/kpi-register/repositories"
	"jethur/grc/kpi-register/services"

	"gorm.io/gorm"
)

type KPIRegisterModule struct {
	KPIRegisterController       *controllers.KPIRegisterController
	KPIValueUpdateController    *controllers.KPIValueUpdateController
	KPIAssociatedItemController *controllers.KPIAssociatedItemController
}

func NewKPIRegisterModule(db *gorm.DB, sharedModule *SharedModule, settingsModule *SettingsModule, masterModule *MasterModule, organizationModule *OrganizationModule, commonModule *CommonModule, riskModule *RiskModule, assetModule *AssetModule, tprmModule *TprmModule, complianceModule *ComplianceModule, documentModule *DocumentModule, strategyModule *StrategyModule) *KPIRegisterModule {
	kpiRegisterRepo := repositories.NewKPIRegisterRepository(db)
	kpiRegisterSvc := services.NewKPIRegisterService(kpiRegisterRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, organizationModule.UserSettingSvc, settingsModule.OrganizationUnitSvc, strategyModule.KpiStatusSvc, strategyModule.EntryTypeSvc, strategyModule.DataSourceSvc)

	kpiValueUpdateRepo := repositories.NewKPIValueUpdateRepository(db)
	kpiValueUpdateSvc := services.NewKPIValueUpdateService(kpiValueUpdateRepo, sharedModule.ActivityLogService, kpiRegisterSvc, masterModule.ReviewFrequencySvc, organizationModule.UserSvc, commonModule.NotificationService, strategyModule.StrategyProfileKpiSvc, strategyModule.StrategyProfileStrategicObjectiveSvc)

	kpiAssociatedItemRepo := repositories.NewKPIAssociatedItemRepository(db)
	kpiAssociatedItemSvc := services.NewKPIAssociatedItemService(kpiAssociatedItemRepo, kpiRegisterSvc, organizationModule.ControlSvc, organizationModule.ProcessSvc, assetModule.AssetRegisterSvc, riskModule.RiskRegisterSvc, tprmModule.ThirdPartyRegisterSvc, complianceModule.ComplianceRegisterSvc, documentModule.DocumentSvc, complianceModule.ComplianceClausevc, commonModule.AssociatedItemValidationService)

	return &KPIRegisterModule{
		KPIRegisterController:       controllers.NewKPIRegisterController(kpiRegisterSvc),
		KPIValueUpdateController:    controllers.NewKPIValueUpdateController(kpiValueUpdateSvc),
		KPIAssociatedItemController: controllers.NewKPIAssociatedItemController(kpiAssociatedItemSvc),
	}
}

package initializers

import (
	"jethur/grc/kpi/controllers"
	"jethur/grc/kpi/repositories"
	"jethur/grc/kpi/services"

	"gorm.io/gorm"
)

type KPIModule struct {
	GRCDashboardController        *controllers.GRCDashboardController
	ERMDashboardController        *controllers.ERMDashboardController
	ComplianceDashboardController *controllers.ComplianceDashboardController
	KRIDashboardController        *controllers.KRIDashboardController

	GrcDashboardSvc services.GRCDashboardServiceInterface
}

func NewKPIModule(db *gorm.DB, riskModule *RiskModule, organizationModule *OrganizationModule, documentModule *DocumentModule) *KPIModule {

	grcDashboardRepo := repositories.NewGRCDashboardRepository(db)
	grcDashboardSvc := services.NewGRCDashboardService(grcDashboardRepo, organizationModule.GovernanceReportSvc, documentModule.DocumentSvc)
	grcDashboardController := controllers.NewGRCDashboardController(grcDashboardSvc, riskModule.RiskHeatMapSvc)

	ermDashboardRepo := repositories.NewERMDashboardRepository(db)
	ermDashboardSvc := services.NewERMDashboardService(ermDashboardRepo)
	ermDashboardController := controllers.NewERMDashboardController(ermDashboardSvc, riskModule.RiskHeatMapSvc)

	complianceDashboardRepo := repositories.NewComplianceDashboardRepository(db)
	complianceDashboardSvc := services.NewComplianceDashboardService(complianceDashboardRepo)
	complianceDashboardController := controllers.NewComplianceDashboardController(complianceDashboardSvc)

	kriDashboardRepo := repositories.NewKRIDashboardRepository(db)
	kriDashboardSvc := services.NewKRIDashboardService(kriDashboardRepo, riskModule.KRIRiskToleranceStatusSvc)
	kriDashboardController := controllers.NewKRIDashboardController(kriDashboardSvc)

	return &KPIModule{
		GRCDashboardController:        grcDashboardController,
		ERMDashboardController:        ermDashboardController,
		ComplianceDashboardController: complianceDashboardController,
		KRIDashboardController:        kriDashboardController,

		GrcDashboardSvc: grcDashboardSvc,
	}
}

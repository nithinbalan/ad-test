package initializers

import (
	"jethur/grc/kpi/controllers"
	"jethur/grc/kpi/repositories"
	"jethur/grc/kpi/services"

	"gorm.io/gorm"
)

type KPIModule struct {
	GRCDashboardController *controllers.GRCDashboardController
	ERMDashboardController *controllers.ERMDashboardController

	GrcDashboardSvc services.GRCDashboardServiceInterface
}

func NewKPIModule(db *gorm.DB, riskModule *RiskModule) *KPIModule {

	grcDashboardRepo := repositories.NewGRCDashboardRepository(db)
	grcDashboardSvc := services.NewGRCDashboardService(grcDashboardRepo)
	grcDashboardController := controllers.NewGRCDashboardController(grcDashboardSvc, riskModule.RiskHeatMapSvc)

	// Erm Dashboard
	ermDashboardRepo := repositories.NewERMDashboardRepository(db)
	ermDashboardSvc := services.NewERMDashboardService(ermDashboardRepo)
	ermDashboardController := controllers.NewERMDashboardController(ermDashboardSvc)

	return &KPIModule{
		GRCDashboardController: grcDashboardController,
		ERMDashboardController: ermDashboardController,

		GrcDashboardSvc: grcDashboardSvc,
	}
}

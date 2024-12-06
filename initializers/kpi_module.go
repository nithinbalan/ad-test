package initializers

import (
	"jethur/grc/kpi/controllers"
	"jethur/grc/kpi/repositories"
	"jethur/grc/kpi/services"

	"gorm.io/gorm"
)

type KPIModule struct {
	GRCDashboardController *controllers.GRCDashboardController
}

func NewKPIModule(db *gorm.DB, riskModule *RiskModule) *KPIModule {

	grcDashboardRepo := repositories.NewGRCDashboardRepository(db)
	grcDashboardSvc := services.NewGRCDashboardService(grcDashboardRepo)
	grcDashboardController := controllers.NewGRCDashboardController(grcDashboardSvc, riskModule.RiskHeatMapSvc)

	return &KPIModule{
		GRCDashboardController: grcDashboardController,
	}
}

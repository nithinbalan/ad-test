package initializers

import (
	"jethur/grc/kpi/controllers"
	"jethur/grc/kpi/repositories"
	"jethur/grc/kpi/services"

	"gorm.io/gorm"
)

type KPIModule struct {
	GRCDashboardController           *controllers.GRCDashboardController
	ERMDashboardController           *controllers.ERMDashboardController
	ComplianceDashboardController    *controllers.ComplianceDashboardController
	KRIDashboardController           *controllers.KRIDashboardController
	AwarenessDashboardController     *controllers.AwarenessDashboardController
	VulnerabilityDashboardController *controllers.VulnerabilityDashboardController
	DocumentDashboardController      *controllers.DocumentDashboardController
	IncidentDashboardController      *controllers.IncidentDashboardController

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
	complianceDashboardSvc := services.NewComplianceDashboardService(complianceDashboardRepo, organizationModule.UserSvc)
	complianceDashboardController := controllers.NewComplianceDashboardController(complianceDashboardSvc)

	kriDashboardRepo := repositories.NewKRIDashboardRepository(db)
	kriDashboardSvc := services.NewKRIDashboardService(kriDashboardRepo, riskModule.KRIRiskToleranceStatusSvc)
	kriDashboardController := controllers.NewKRIDashboardController(kriDashboardSvc)

	awarenessDashboardRepo := repositories.NewAwarenessDashboardRepository(db)
	awarenessDashboardSvc := services.NewAwarenessDashboardService(awarenessDashboardRepo)
	awarenessDashboardController := controllers.NewAwarenessDashboardController(awarenessDashboardSvc)

	vulnerabilityDashboardRepo := repositories.NewVulnerabilityDashboardRepository(db)
	vulnerabilityDashboardSvc := services.NewVulnerabilityDashboardService(vulnerabilityDashboardRepo)
	vulnerabilityDashboardController := controllers.NewVulnerabilityDashboardController(vulnerabilityDashboardSvc, db)

	documentDashboardRepo := repositories.NewDocumentDashboardRepository(db)
	documentDashboardSvc := services.NewDocumentDashboardService(documentDashboardRepo)
	documentDashboardController := controllers.NewDocumentDashboardController(documentDashboardSvc, db)

	incidentDashboardRepo := repositories.NewIncidentDashboardRepository(db)
	incidentDashboardSvc := services.NewIncidentDashboardService(incidentDashboardRepo)
	incidentDashboardController := controllers.NewIncidentDashboardController(incidentDashboardSvc, db)

	return &KPIModule{
		GRCDashboardController:           grcDashboardController,
		ERMDashboardController:           ermDashboardController,
		ComplianceDashboardController:    complianceDashboardController,
		KRIDashboardController:           kriDashboardController,
		AwarenessDashboardController:     awarenessDashboardController,
		VulnerabilityDashboardController: vulnerabilityDashboardController,
		DocumentDashboardController:      documentDashboardController,
		IncidentDashboardController:      incidentDashboardController,

		GrcDashboardSvc: grcDashboardSvc,
	}
}

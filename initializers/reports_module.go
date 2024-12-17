package initializers

import (
	"jethur/grc/reports/controllers"
	"jethur/grc/reports/repositories"
	"jethur/grc/reports/services"

	"gorm.io/gorm"
)

type ReportsModule struct {
	ComplianceReportController   *controllers.ComplianceReportController
	OrganizationReportController *controllers.OrganizationReportController
	RiskReportController         *controllers.RiskReportController
}

func NewReportsModule(db *gorm.DB,
	complianceModule *ComplianceModule,
	organizationModule *OrganizationModule,
	riskModule *RiskModule) *ReportsModule {

	organizationReportRepo := repositories.NewOrganizationReportRepository(db)
	organizationReportSvc := services.NewOrganizationReportService(organizationReportRepo)
	organizationReportController := controllers.NewOrganizationReportController(organizationReportSvc, organizationModule.UserSvc, organizationModule.ProcessSvc, organizationModule.ControlSvc)

	complianceReportRepo := repositories.NewComplianceReportRepository(db)
	complianceReportSvc := services.NewComplianceReportService(complianceReportRepo)
	complianceReportController := controllers.NewComplianceReportController(complianceReportSvc,
		complianceModule.ComplianceRegisterSvc,
		complianceModule.SlaRegisterSvc,
		complianceModule.ComplianceRequirementSvc,
		complianceModule.ComplianceActionPlanSvc)

	riskReportRepo := repositories.NewRiskReportRepository(db)
	riskReportSvc := services.NewRiskReportService(riskReportRepo)
	riskReportController := controllers.NewRiskReportController(riskReportSvc, riskModule.kriRegisterSvc)

	return &ReportsModule{
		ComplianceReportController:   complianceReportController,
		OrganizationReportController: organizationReportController,
		RiskReportController:         riskReportController,
	}
}

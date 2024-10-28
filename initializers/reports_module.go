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
}

func NewReportsModule(db *gorm.DB, complianceModule *ComplianceModule, organizationModule *OrganizationModule) *ReportsModule {

	organizationReportRepo := repositories.NewOrganizationReportRepository(db)
	organizationReportSvc := services.NewOrganizationReportService(organizationReportRepo)
	organizationReportController := controllers.NewOrganizationReportController(organizationReportSvc, organizationModule.UserSvc, organizationModule.ProcessSvc)

	complianceReportRepo := repositories.NewComplianceReportRepository(db)
	complianceReportSvc := services.NewComplianceReportService(complianceReportRepo)
	complianceReportController := controllers.NewComplianceReportController(complianceReportSvc, complianceModule.ComplianceRegisterSvc)

	return &ReportsModule{
		ComplianceReportController:   complianceReportController,
		OrganizationReportController: organizationReportController,
	}
}

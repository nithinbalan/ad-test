package initializers

import (
	"jethur/grc/reports/controllers"
	"jethur/grc/reports/repositories"
	"jethur/grc/reports/services"

	"gorm.io/gorm"
)

type ReportsModule struct {
	ComplianceReportController *controllers.ComplianceReportController
}

func NewReportsModule(db *gorm.DB, complianceModule *ComplianceModule) *ReportsModule {

	complianceReportRepo := repositories.NewComplianceReportRepository(db)
	complianceReportSvc := services.NewComplianceReportService(complianceReportRepo)
	complianceReportController := controllers.NewComplianceReportController(complianceReportSvc, complianceModule.ComplianceRegisterSvc)

	return &ReportsModule{
		ComplianceReportController: complianceReportController,
	}
}

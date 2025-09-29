package repositories

import (
	"gorm.io/gorm"
	"jethur/grc/models"
	"jethur/utility"
)

type IncidentReportRepositoryInterface interface {
	GetFullReportDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentRegister, error)
}

type IncidentReportRepository struct {
	db *gorm.DB
}

func NewIncidentReportRepository(db *gorm.DB) *IncidentReportRepository {
	return &IncidentReportRepository{
		db: db,
	}
}

func (repo *IncidentReportRepository) GetFullReportDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incident *models.IncidentRegister

	result := db.
		Preload("IncidentCategory").
		Preload("IncidentType").
		Preload("IncidentStatus").
		Preload("ReportedByUser").
		Preload("SeverityLevel").
		Preload("SLARegisters").
		Preload("AssociatedDocuments").
		Preload("ResponsibleUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Investigation").
		Preload("Investigation.TeamLead").
		Preload("Investigation.Users").
		Preload("Investigation.Users.User").
		Preload("Investigation.Users.User.Designation"). // if you want to show designation
		Preload("Investigation.Documents").
		Preload("Investigation.Findings").
		Preload("Assets").
		Preload("Assets.AssetRegister").
		Preload("Assets.AssetRegister.AssetCategory").
		Preload("Assets.AssetRegister.AssetStatus").
		Preload("Documents").
		Preload("Controls.Control").
		Preload("Controls.Control.ControlImplementationStatus").
		Preload("Controls.Control.ControlEfficiencyStatus").
		Preload("Controls.Control.LifecycleStatus").
		Preload("Controls.Control.ResponsibleUser").
		Preload("Controls.Control.ResponsibleUnit").
		Preload("Controls.Control.CreatedByUser").
		Preload("Controls.Control.UpdatedByUser").
		Preload("Risks.RiskRegister").
		Preload("Risks.RiskRegister.RiskCategory").
		Preload("Risks.RiskRegister.RiskStatus").
		Preload("Risks.RiskRegister.CurrentLikelihood").
		Preload("Risks.RiskRegister.CurrentImpact").
		Preload("Risks.RiskRegister.CurrentRiskRating").
		Preload("Processes").
		Preload("Processes.Process").
		Preload("Processes.Process.LifecycleStatus").
		Preload("Processes.Process.ResponsibleUser").
		Preload("ComplianceRequirements.ComplianceRequirement").
		Preload("ComplianceRequirements.ComplianceRequirement.ComplianceStatus").
		Preload("Timelines").
		Preload("Timelines.PerformedBy").
		Preload("Timelines.TimelineStatus").
		Preload("Timelines.IncidentTimelineType").
		Preload("CommunicationLogs").
		Preload("ActionPlans").
		Preload("ActionPlans.ResponsibleUser").
		Preload("ActionPlans.SeverityLevel").
		Preload("ActionPlans.ActionPlanStatus").
		Preload("IncidentClosure").
		Preload("IncidentPostReview").
		Preload("IncidentClosure.Lessons").
		Preload("IncidentPostReview").
		Preload("IncidentPostReview.Users").
		Preload("IncidentPostReview.Users.User").
		Preload("IncidentPostReview.Recommendations").
		Preload("IncidentPostReview.Lessons").
		Preload("IncidentPostReview.IncidentPostReviewEffectiveness").
		Preload("IncidentRootCauseAnalyses").
		Preload("StrategicObjectives").
		Preload("StrategicObjectives.StrategicObjective").
		Preload("KPIs.KPIRegister").
		Preload("KPIs.KPIRegister.KPICategory").
		Preload("KPIs.KPIRegister.KPIStatus").
		First(&incident, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incident, nil
}

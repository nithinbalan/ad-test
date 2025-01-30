package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	reportsRoutes := router.Group("/")

	reportsRoutes.Use(middlewares.AuthMiddleware())

	complianceReportController := app.ReportsModule.ComplianceReportController
	organizationReportController := app.ReportsModule.OrganizationReportController
	riskReportController := app.ReportsModule.RiskReportController
	auditReportController := app.ReportsModule.AuditReportController
	incidentReportController := app.ReportsModule.IncidentReportController
	meetingReportController := app.ReportsModule.MeetingReportController

	//Compliance
	reportsRoutes.GET("/reports/compliance-registers", complianceReportController.GetAllComplianceRegisters)
	reportsRoutes.GET("/reports/compliance-registers/export", complianceReportController.ComplianceRegisterReportExport)
	reportsRoutes.GET("/reports/compliance-by-compliance-category", complianceReportController.GetComplianceByComplianceCategory)
	reportsRoutes.GET("/reports/compliance-by-severity-level", complianceReportController.GetComplianceBySeverityLevel)
	reportsRoutes.GET("/reports/compliance-by-compliance-status", complianceReportController.GetComplianceByComplianceStatus)
	reportsRoutes.GET("/reports/compliance-by-responsible-unit", complianceReportController.GetComplianceByResponsibleUnit)
	reportsRoutes.GET("/reports/compliance-by-responsible-user", complianceReportController.GetComplianceByResponsibleUser)

	reportsRoutes.GET("/reports/compliance-requirements", complianceReportController.GetAllComplianceRequirements)
	reportsRoutes.GET("/reports/compliance-requirement-by-compliance-status", complianceReportController.GetComplianceRequirementByComplianceStatus)
	reportsRoutes.GET("/reports/compliance-requirement-by-severity-level", complianceReportController.GetComplianceRequirementBySeverityLevel)
	reportsRoutes.GET("/reports/compliance-requirement-by-category", complianceReportController.GetComplianceRequirementByCategory)
	reportsRoutes.GET("/reports/compliance-requirement-by-responsible-user", complianceReportController.GetComplianceRequirementByResponsibleUser)

	reportsRoutes.GET("/reports/compliance-action-plans", complianceReportController.GetAllComplianceActionPlans)
	reportsRoutes.GET("/reports/compliance-action-plans/export", complianceReportController.GetComplianceActionPlanExport)
	reportsRoutes.GET("/reports/compliance-action-plan-by-status", complianceReportController.GetComplianceActionPlanByStatus)
	reportsRoutes.GET("/reports/compliance-action-plan-by-responsible-user", complianceReportController.GetComplianceActionPlanByResponsibleUser)

	reportsRoutes.GET("/reports/slas", complianceReportController.GetAllSLARegisters)
	reportsRoutes.GET("/reports/slas/export", complianceReportController.GetSLARegisterExport)
	reportsRoutes.GET("/reports/sla-by-category", complianceReportController.GetSLAByCategory)
	reportsRoutes.GET("/reports/sla-by-severity-level", complianceReportController.GetSLABySeverityLevel)
	reportsRoutes.GET("/reports/sla-by-compliance-status", complianceReportController.GetSLAByComplianceStatus)

	//Organization
	reportsRoutes.GET("/reports/users", organizationReportController.GetAllUsers)
	reportsRoutes.GET("/reports/users/export", organizationReportController.GetUsersExport)
	reportsRoutes.GET("/reports/users-by-status", organizationReportController.GetUsersByStatus)
	reportsRoutes.GET("/reports/users-by-role", organizationReportController.GetUsersByRole)
	reportsRoutes.GET("/reports/users-by-designation", organizationReportController.GetUsersByDesignation)

	reportsRoutes.GET("/reports/processes", organizationReportController.GetAllProcesses)
	reportsRoutes.GET("/reports/processes/export", organizationReportController.GetProcessesExport)
	reportsRoutes.GET("/reports/processes-by-responsible-unit", organizationReportController.GetProcessesByResponsibleUnit)
	reportsRoutes.GET("/reports/processes-by-risk-rating", organizationReportController.GetProcessesByRiskRating)
	reportsRoutes.GET("/reports/processes-by-responsible-user", organizationReportController.GetProcessesByResponsibleUser)
	reportsRoutes.GET("/reports/processes-by-process-group", organizationReportController.GetProcessesByProcessGroup)

	reportsRoutes.GET("/reports/controls", organizationReportController.GetAllControls)
	reportsRoutes.GET("/reports/controls/export", organizationReportController.GetControlsExport)
	reportsRoutes.GET("/reports/controls-by-control-type", organizationReportController.GetControlsByControlType)
	reportsRoutes.GET("/reports/controls-by-control-category", organizationReportController.GetControlsByCategory)
	reportsRoutes.GET("/reports/controls-by-compliance-register", organizationReportController.GetControlsByComplianceRegister)
	reportsRoutes.GET("/reports/controls-by-risk-rating", organizationReportController.GetControlsByRiskRating)
	reportsRoutes.GET("/reports/controls-by-implementation-status", organizationReportController.GetControlsByImplementationStatus)
	reportsRoutes.GET("/reports/controls-by-efficiency-status", organizationReportController.GetControlsByEfficiencyStatus)
	reportsRoutes.GET("/reports/controls-by-applicability", organizationReportController.GetControlsByApplicability)

	//Risk
	reportsRoutes.GET("/reports/kris", riskReportController.GetAllKRIRegisters)
	reportsRoutes.GET("/reports/kris/export", riskReportController.GetKriRegisterExport)
	reportsRoutes.GET("/reports/kri-by-category", riskReportController.GetKRIByCategory)
	reportsRoutes.GET("/reports/kri-by-responsible-unit", riskReportController.GetKRIByResponsibleUnit)
	reportsRoutes.GET("/reports/kri-by-responsible-user", riskReportController.GetKRIByResponsibleUser)
	reportsRoutes.GET("/reports/kri-by-risk-level", riskReportController.GetKRIByRiskLevel)
	reportsRoutes.GET("/reports/kri-by-review-frequency", riskReportController.GetKRIByReviewFrequency)

	reportsRoutes.GET("/reports/risks", riskReportController.GetAllRiskRegisters)
	reportsRoutes.GET("/reports/risks/export", riskReportController.GetRiskRegisterExport)
	reportsRoutes.GET("/reports/risk-by-category", riskReportController.GetRiskByCategory)
	reportsRoutes.GET("/reports/risk-by-responsible-unit", riskReportController.GetRiskByResponsibleUnit)
	reportsRoutes.GET("/reports/risk-by-responsible-user", riskReportController.GetRiskByResponsibleUser)
	reportsRoutes.GET("/reports/risk-by-risk-source", riskReportController.GetRiskByRiskSource)
	reportsRoutes.GET("/reports/risk-by-inherent-risk-rating", riskReportController.GetRiskByInherentRiskRating)
	reportsRoutes.GET("/reports/risk-by-residual-risk-rating", riskReportController.GetRiskByResidualRiskRating)
	reportsRoutes.GET("/reports/risk-by-risk-status", riskReportController.GetRiskByRiskStatus)
	reportsRoutes.GET("/reports/risk-by-risk-response-strategy", riskReportController.GetRiskByRiskResponseStrategy)
	reportsRoutes.GET("/reports/risk-by-risk-appetite-level", riskReportController.GetRiskByRiskAppetiteToleranceLevel)
	reportsRoutes.GET("/reports/risk-by-kri-enabled-status", riskReportController.GetRiskByKRIEnabledStatus)

	reportsRoutes.GET("/reports/risk-action-plans", riskReportController.GetAllRiskActionPlans)
	reportsRoutes.GET("/reports/risk-action-plans/export", riskReportController.GetRiskActionPlanExport)
	reportsRoutes.GET("/reports/risk-action-plan-by-status", riskReportController.GetRiskActionPlansByStatus)
	reportsRoutes.GET("/reports/risk-action-plan-by-responsible-user", riskReportController.GetRiskActionPlansByResponsibleUser)
	reportsRoutes.GET("/reports/risk-action-plan-by-severity-level", riskReportController.GetRiskActionPlansBySeverityLevel)

	//Audits Plan
	reportsRoutes.GET("/reports/audit-plans", auditReportController.GetAllAuditPlans)
	reportsRoutes.GET("/reports/audit-plans/export", auditReportController.GetAuditPlanExport)
	reportsRoutes.GET("/reports/audit-plan-by-category", auditReportController.GetAuditPlanByCategory)
	reportsRoutes.GET("/reports/audit-plan-by-status", auditReportController.GetAuditPlanByStatus)
	reportsRoutes.GET("/reports/audit-plan-by-team", auditReportController.GetAuditPlanByTeam)
	reportsRoutes.GET("/reports/audit-plan-by-team-lead", auditReportController.GetAuditPlanByTeamLead)
	reportsRoutes.GET("/reports/audit-plan-by-year", auditReportController.GetAuditPlanByYear)

	//Audits
	reportsRoutes.GET("/reports/audits", auditReportController.GetAllAudits)
	reportsRoutes.GET("/reports/audits/export", auditReportController.GetAuditExport)
	reportsRoutes.GET("/reports/audit-by-category", auditReportController.GetAuditByCategory)
	reportsRoutes.GET("/reports/audit-by-status", auditReportController.GetAuditByStatus)
	reportsRoutes.GET("/reports/audit-by-open-vs-closed", auditReportController.GetAuditByOpenVsClosed)
	reportsRoutes.GET("/reports/audit-by-month", auditReportController.GetAuditByMonth)
	reportsRoutes.GET("/reports/audit-by-team-lead", auditReportController.GetAuditByTeamLead)

	//Audit findings
	reportsRoutes.GET("/reports/audit-findings", auditReportController.GetAllAuditFindings)
	reportsRoutes.GET("/reports/audit-findings/export", auditReportController.GetAuditFindingExport)
	reportsRoutes.GET("/reports/audit-finding-by-category", auditReportController.GetAuditFindingByCategory)
	reportsRoutes.GET("/reports/audit-finding-by-status", auditReportController.GetAuditFindingByStatus)
	reportsRoutes.GET("/reports/audit-finding-by-severity-level", auditReportController.GetAuditFindingBySeverityLevel)
	reportsRoutes.GET("/reports/audit-finding-by-responsible-unit", auditReportController.GetAuditFindingByResponsibleUnit)
	reportsRoutes.GET("/reports/audit-finding-by-month", auditReportController.GetAuditFindingByMonth)
	reportsRoutes.GET("/reports/audit-finding-by-auditor", auditReportController.GetAuditFindingByAuditor)

	//Audit Action plans
	reportsRoutes.GET("/reports/audit-action-plans", auditReportController.GetAllAuditActionPlans)
	reportsRoutes.GET("/reports/audit-action-plans/export", auditReportController.GetAuditActionPlanExport)
	reportsRoutes.GET("/reports/audit-action-plan-by-status", auditReportController.GetAuditActionPlanByStatus)
	reportsRoutes.GET("/reports/audit-action-plan-by-severity-level", auditReportController.GetAuditActionPlanBySeverityLevel)
	reportsRoutes.GET("/reports/audit-action-plan-by-responsible-user", auditReportController.GetAuditActionPlanByResponsibleUser)
	reportsRoutes.GET("/reports/audit-action-plan-by-responsible-unit", auditReportController.GetAuditActionPlanByResponsibleUnit)
	reportsRoutes.GET("/reports/audit-action-plan-by-on-time-vs-delayed", auditReportController.GetAuditActionPlanOnTimeVsDelayed)

	//Incidents
	reportsRoutes.GET("/reports/incidents", incidentReportController.GetAllIncidents)
	reportsRoutes.GET("/reports/incidents/export", incidentReportController.GetIncidentExport)
	reportsRoutes.GET("/reports/incident-by-category", incidentReportController.GetIncidentByCategory)
	reportsRoutes.GET("/reports/incident-by-status", incidentReportController.GetIncidentByStatus)
	reportsRoutes.GET("/reports/incident-by-severity-level", incidentReportController.GetIncidentBySeverityLevel)
	reportsRoutes.GET("/reports/incident-by-responsible-unit", incidentReportController.GetIncidentByResponsibleUnit)
	reportsRoutes.GET("/reports/incident-by-reported-user", incidentReportController.GetIncidentByReportedUser)
	reportsRoutes.GET("/reports/incident-by-reported-year", incidentReportController.GetIncidentByReportedYear)

	reportsRoutes.GET("/reports/incident-by-asset", incidentReportController.GetIncidentByAsset)
	reportsRoutes.GET("/reports/incident-by-sla", incidentReportController.GetIncidentBySLA)

	reportsRoutes.GET("/reports/incident-investigation-summary", incidentReportController.GetIncidentInvestigationSummary)
	reportsRoutes.GET("/reports/incident-rca", incidentReportController.GetIncidentRCA)
	reportsRoutes.GET("/reports/incident-closure", incidentReportController.GetIncidentClosure)
	reportsRoutes.GET("/reports/incident-by-risk", incidentReportController.GetOpenIncidentsWithRisk)
	reportsRoutes.GET("/reports/incident-pending-action-plan", incidentReportController.GetIncidentPendingActionPlan)

	reportsRoutes.GET("/reports/incident-action-plans", incidentReportController.GetAllIncidentActionPlans)
	reportsRoutes.GET("/reports/incident-action-plans/export", incidentReportController.GetIncidentActionPlanExport)
	reportsRoutes.GET("/reports/incident-action-plan-by-status", incidentReportController.GetIncidentActionPlanByStatus)
	reportsRoutes.GET("/reports/incident-action-plan-by-responsible-user", incidentReportController.GetIncidentActionPlanByResponsibleUser)
	reportsRoutes.GET("/reports/incident-action-plan-by-responsible-unit", incidentReportController.GetIncidentActionPlanByResponsibleUnit)

	//Meetings
	reportsRoutes.GET("/reports/meetings", meetingReportController.GetAllMeetings)
	reportsRoutes.GET("/reports/meetings/export", meetingReportController.GetMeetingExport)
	reportsRoutes.GET("/reports/meeting-by-type", meetingReportController.GetMeetingByType)
	reportsRoutes.GET("/reports/meeting-by-status", meetingReportController.GetMeetingByStatus)
	reportsRoutes.GET("/reports/meeting-by-responsible-unit", meetingReportController.GetMeetingByResponsibleUnit)
	reportsRoutes.GET("/reports/meeting-by-organizer", meetingReportController.GetMeetingByOrganizer)
	reportsRoutes.GET("/reports/meeting-minutes-by-status", meetingReportController.GetMeetingMinutesByStatus)

	reportsRoutes.GET("/reports/meeting-action-plan-by-status", meetingReportController.GetMeetingActionPlanByStatus)
	reportsRoutes.GET("/reports/meeting-action-plan-by-responsible-unit", meetingReportController.GetMeetingActionPlanByResponsibleUnit)
	reportsRoutes.GET("/reports/meeting-action-plan-by-responsible-user", meetingReportController.GetMeetingActionPlanByResponsibleUser)
	reportsRoutes.GET("/reports/meeting-action-plan-by-open-vs-closed", meetingReportController.GetMeetingActionPlanByOpenVsClosed)
	reportsRoutes.GET("/reports/meeting-mom-count", meetingReportController.GetMeetingMOMCount)
	reportsRoutes.GET("/reports/meeting-with-mom", meetingReportController.GetMeetingsWithMOM)
	reportsRoutes.GET("/reports/meeting-with-open-action-plan", meetingReportController.GetMeetingsWithOpenActionPlans)
	reportsRoutes.GET("/reports/meeting-with-overdue-action-plan", meetingReportController.GetMeetingsWithOverdueActionPlans)
	reportsRoutes.GET("/reports/meeting-by-year", meetingReportController.GetMeetingCountByYear)
	reportsRoutes.GET("/reports/meeting-by-user", meetingReportController.GetMeetingCountByUser)

	reportsRoutes.GET("/reports/meeting-action-plans", meetingReportController.GetAllMeetingActionPlans)
	reportsRoutes.GET("/reports/meeting-action-plans/export", meetingReportController.GetMeetingActionPlanExport)

}

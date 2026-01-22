package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	complianceRoutes := router.Group("/")

	complianceRoutes.Use(middlewares.AuthMiddleware(app.DB))

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	complianceCategoryController := app.ComplianceModule.ComplianceCategoryController
	complianceStatusController := app.ComplianceModule.ComplianceStatusController
	complianceReviewStatusController := app.ComplianceModule.ComplianceReviewStatusController
	complianceRegisterController := app.ComplianceModule.ComplianceRegisterController
	complianceAssociatedItemController := app.ComplianceModule.ComplianceAssociatedItemController
	complianceRequirementAssociatedItemController := app.ComplianceModule.ComplianceRequirementAssociatedItemController
	complianceRequirementCategoryController := app.ComplianceModule.ComplianceRequirementCategoryController
	complianceRequirementController := app.ComplianceModule.ComplianceRequirementController
	complianceActionPlanController := app.ComplianceModule.ComplianceActionPlanController
	complianceStatusUpdateController := app.ComplianceModule.ComplianceStatusUpdateController
	complianceRequirementStatusUpdateController := app.ComplianceModule.ComplianceRequirementStatusUpdateController
	complianceActionPlanStatusUpdateController := app.ComplianceModule.ComplianceActionPlanStatusUpdateController
	complianceReviewController := app.ComplianceModule.ComplianceReviewController
	complianceClauseController := app.ComplianceModule.CompliacneClauseController
	SLACategoryController := app.ComplianceModule.SLACategoryController
	SLARegisterController := app.ComplianceModule.SLARegisterController
	SLACheckQuestionController := app.ComplianceModule.SLACheckQuestionController
	SLACheckAnswerController := app.ComplianceModule.SLACheckAnswerController
	SLARenewController := app.ComplianceModule.SLARenewController
	slaActionPlanController := app.ComplianceModule.SLAActionPlanController
	slaActionPlanCommentController := app.ComplianceModule.SLAActionPlanCommentController
	slaActionPlanStatusUpdateController := app.ComplianceModule.SLAActionPlanStatusUpdateController
	slaAssociatedItemController := app.ComplianceModule.SLAAssociatedItemController
	ComplianceActionPlanCommentController := app.ComplianceModule.ComplianceActionPlanCommentController
	complianceAssessmentController := app.ComplianceModule.ComplianceAssessmentController
	complianceAssessmentRequirementController := app.ComplianceModule.ComplianceAssessmentRequirementController
	complianceFolderController := app.ComplianceModule.ComplianceFolderController
	complianceAssessmentStatusController := app.ComplianceModule.ComplianceAssessmentStatusController
	complianceAssessmentRequirementStatusController := app.ComplianceModule.ComplianceAssessmentRequirementStatusController
	complianceAssessmentRequirementCommentController := app.ComplianceModule.ComplianceAssessmentRequirementCommentController
	complianceAssessmentRequirementActionPlanController := app.ComplianceModule.ComplianceAssessmentRequirementActionPlanController
	complianceAssessmentRequirementActionPlanStatusUpdateController := app.ComplianceModule.ComplianceAssessmentRequirementActionPlanStatusUpdateController
	complianceAssessmentRequirementActionPlanCommentController := app.ComplianceModule.ComplianceAssessmentRequirementActionPlanCommentController
	complianceReportController := app.ComplianceModule.ComplianceReportController
	policyExceptionController := app.ComplianceModule.PolicyExceptionController
	complianceAssessmentFrameworkTypeController := app.ComplianceModule.ComplianceAssessmentFrameworkTypeController
	complianceAssessmentFrameworkController := app.ComplianceModule.ComplianceAssessmentFrameworkController
	complianceAssessmentFrameworkValueController := app.ComplianceModule.ComplianceAssessmentFrameworkValueController
	complianceAssessmentControlController := app.ComplianceModule.ComplianceAssessmentControlController
	complianceAssessmentControlCommentController := app.ComplianceModule.ComplianceAssessmentControlCommentController
	complianceAssessmentControlActionPlanController := app.ComplianceModule.ComplianceAssessmentControlActionPlanController
	complianceAssessmentControlActionPlanStatusUpdateController := app.ComplianceModule.ComplianceAssessmentControlActionPlanStatusUpdateController
	complianceAssessmentControlActionPlanCommentController := app.ComplianceModule.ComplianceAssessmentControlActionPlanCommentController
	complianceDetailDashboardController := app.ComplianceModule.ComplianceDetailDashboardController
	consultationRequestStatusController := app.ComplianceModule.ConsultationRequestStatusController
	LegislativeBodyController := app.ComplianceModule.LegislativeBodyController
	LegislationTypeController := app.ComplianceModule.LegislationTypeController

	// Compliance Category
	complianceRoutes.GET("/compliance-categories", complianceCategoryController.GetAllComplianceCategories)
	complianceRoutes.GET("/compliance-categories/:id", complianceCategoryController.GetComplianceCategoryDetails)
	complianceRoutes.POST("/compliance-categories", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceCategoryController.CreateComplianceCategory)
	complianceRoutes.PUT("/compliance-categories/:id", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceCategoryController.UpdateComplianceCategory)
	complianceRoutes.POST("/compliance-categories/delete", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceCategoryController.DeleteComplianceCategories)
	complianceRoutes.PUT("/compliance-categories/activate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceCategoryController.ActivateComplianceCategories)
	complianceRoutes.PUT("/compliance-categories/deactivate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceCategoryController.DeactivateComplianceCategories)

	// Compliance Status
	complianceRoutes.GET("/compliance-statuses", complianceStatusController.GetAllComplianceStatuses)
	complianceRoutes.GET("/compliance-statuses/:id", complianceStatusController.GetComplianceStatusDetails)

	//consultation-request-statuses
	complianceRoutes.GET("/consultation-request-statuses", consultationRequestStatusController.GetAllConsultationRequestStatuses)
	complianceRoutes.GET("/consultation-request-statuses/:id", consultationRequestStatusController.GetConsultationRequestStatusDetails)

	// Compliance Requirement Category
	complianceRoutes.GET("/compliance-requirement-categories", complianceRequirementCategoryController.GetAllComplianceRequirementCategories)
	complianceRoutes.GET("/compliance-requirement-categories/:id", complianceRequirementCategoryController.GetComplianceRequirementCategoryDetails)
	complianceRoutes.POST("/compliance-requirement-categories", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceRequirementCategoryController.CreateComplianceRequirementCategory)
	complianceRoutes.PUT("/compliance-requirement-categories/:id", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceRequirementCategoryController.UpdateComplianceRequirementCategory)
	complianceRoutes.POST("/compliance-requirement-categories/delete", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceRequirementCategoryController.DeleteComplianceRequirementCategories)
	complianceRoutes.PUT("/compliance-requirement-categories/activate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceRequirementCategoryController.ActivateComplianceRequirementCategories)
	complianceRoutes.PUT("/compliance-requirement-categories/deactivate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceRequirementCategoryController.DeactivateComplianceRequirementCategories)

	// SLA Category
	complianceRoutes.GET("/sla-categories/:id", SLACategoryController.GetSLACategoryDetails)
	complianceRoutes.POST("/sla-categories", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), SLACategoryController.CreateCategory)
	complianceRoutes.GET("/sla-categories", SLACategoryController.GetAllSLACategories)
	complianceRoutes.PUT("/sla-categories/:id", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), SLACategoryController.UpdateCategory)
	complianceRoutes.POST("/sla-categories/delete", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), SLACategoryController.DeleteSLACategories)
	complianceRoutes.PUT("/sla-categories/activate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), SLACategoryController.ActivateSLACategories)
	complianceRoutes.PUT("/sla-categories/deactivate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), SLACategoryController.DeactivateSLACategories)
	// Legislative Body
	complianceRoutes.GET("/legislative-bodies/:id", LegislativeBodyController.GetLegislativeBodyDetails)
	complianceRoutes.POST("/legislative-bodies", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), LegislativeBodyController.CreateLegislativeBody)
	complianceRoutes.GET("/legislative-bodies", LegislativeBodyController.GetAllLegislativeBodies)
	complianceRoutes.PUT("/legislative-bodies/:id", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), LegislativeBodyController.UpdateLegislativeBody)
	complianceRoutes.POST("/legislative-bodies/delete", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), LegislativeBodyController.DeleteLegislativeBodies)

	// Legislation Type routes
	complianceRoutes.GET("/legislation-types/:id", LegislationTypeController.GetLegislationTypeDetails)
	complianceRoutes.POST("/legislation-types", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), LegislationTypeController.CreateLegislationType)
	complianceRoutes.GET("/legislation-types", LegislationTypeController.GetAllLegislationTypes)
	complianceRoutes.PUT("/legislation-types/:id", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), LegislationTypeController.UpdateLegislationType)
	complianceRoutes.POST("/legislation-types/delete", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), LegislationTypeController.DeleteLegislationTypes)
	complianceRoutes.POST("/legislation-types/activate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), LegislationTypeController.ActivateLegislationTypes)
	complianceRoutes.POST("/legislation-types/deactivate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), LegislationTypeController.DeactivateLegislationTypes)

	// Compliance Review Statuses
	complianceRoutes.GET("/compliance-review-statuses", complianceReviewStatusController.GetAllComplianceReviewStatuses)
	complianceRoutes.GET("/compliance-review-statuses/:id", complianceReviewStatusController.GetComplianceReviewStatusDetails)

	//compliance folder
	complianceRoutes.GET("/compliance-folders", complianceFolderController.GetAllComplianceFolders)
	complianceRoutes.GET("/compliance-folders/:id", complianceFolderController.GetComplianceFolderDetails)
	complianceRoutes.POST("/compliance-folders", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceFolderController.CreateComplianceFolder)
	complianceRoutes.PUT("/compliance-folders/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceFolderController.UpdateComplianceFolder)
	complianceRoutes.POST("/compliance-folders/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceFolderController.DeleteComplianceFolders)
	complianceRoutes.PUT("/compliance-folders/activate", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceFolderController.ActivateComplianceFolders)
	complianceRoutes.PUT("/compliance-folders/deactivate", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceFolderController.DeactivateComplianceFolders)

	// Compliance Register
	unauthenticatedRoutes.GET("/compliance-register/:id/document-preview", complianceRegisterController.GetComplianceRegisterDocumentPreview)
	unauthenticatedRoutes.GET("/compliance-register/:id/document-download", complianceRegisterController.DownloadComplianceRegisterDocument)
	complianceRoutes.GET("/compliance-register/export", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRegisterController.GetComplianceRegisterExport)
	complianceRoutes.GET("/compliance-register", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRegisterController.GetAllComplianceRegisters)
	complianceRoutes.GET("/compliance-register/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRegisterController.GetComplianceRegisterDetails)
	complianceRoutes.POST("/compliance-register", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRegisterController.CreateComplianceRegister)
	complianceRoutes.PUT("/compliance-register/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRegisterController.UpdateComplianceRegister)
	complianceRoutes.POST("/compliance-register/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRegisterController.DeleteComplianceRegisters)
	complianceRoutes.POST("/compliance-register/force-delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRegisterController.DeleteForceComplianceRegisters)
	complianceRoutes.PUT("/compliance-register/archive", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRegisterController.ArchiveComplianceRegisters)
	complianceRoutes.PUT("/compliance-register/restore", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRegisterController.RestoreComplianceRegisters)
	complianceRoutes.GET("/compliance-register/:id/versions", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRegisterController.GetAllVersions)
	complianceRoutes.POST("/compliance-register/:id/versions", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRegisterController.CreateComplianceVersion)

	// Compliance Associated Items
	complianceRoutes.GET("/compliances/:complianceId/controls", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllComplianceControls)
	complianceRoutes.GET("/compliances/:complianceId/risks", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllComplianceRisks)
	complianceRoutes.GET("/compliances/:complianceId/incidents", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllComplianceIncidents)
	complianceRoutes.GET("/compliances/:complianceId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllComplianceStrategicObjectives)
	complianceRoutes.GET("/compliances/:complianceId/assets", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllComplianceAssets)
	complianceRoutes.GET("/compliances/:complianceId/processes", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllComplianceProcesses)
	complianceRoutes.GET("/compliances/:complianceId/kpis", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllComplianceKPI)
	complianceRoutes.GET("/compliances/:complianceId/associated-item-summary", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllAssociatedCounts)
	complianceRoutes.GET("/compliances/:complianceId/requirements-status-summary", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetComplianceRequirementCountByStatus)
	complianceRoutes.GET("/compliances/:complianceId/control-status-summary", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetComplianceControlCountByStatus)

	complianceRoutes.POST("/compliances/:complianceId/controls", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceAssociatedItemController.CreateComplianceControls)
	complianceRoutes.DELETE("/compliances/:complianceId/controls/:controlId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceAssociatedItemController.DeleteComplianceControl)

	complianceRoutes.POST("/compliances/:complianceId/risks", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceAssociatedItemController.CreateComplianceRisk)
	complianceRoutes.DELETE("/compliances/:complianceId/risks/:riskId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceAssociatedItemController.DeleteComplianceRisk)

	complianceRoutes.GET("/compliances/:complianceId/documents", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllComplianceDocuments)
	complianceRoutes.POST("/compliances/:complianceId/documents", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceAssociatedItemController.CreateComplianceDocument)
	complianceRoutes.DELETE("/compliances/:complianceId/documents/:documentId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceAssociatedItemController.DeleteComplianceDocument)

	// SLA Register
	unauthenticatedRoutes.GET("/sla-register/:id/document-download", SLARegisterController.DownloadSLARegisterDocument)
	unauthenticatedRoutes.GET("/sla-register/:id/document-preview", SLARegisterController.GetSLARegisterDocumentPreview)
	complianceRoutes.GET("/sla-register/export", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLARegisterController.GetSLARegisterExport)
	complianceRoutes.GET("/sla-register", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLARegisterController.GetAllSLARegisters)
	complianceRoutes.GET("/sla-register/:id", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLARegisterController.GetSLARegisterDetails)
	complianceRoutes.POST("/sla-register", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARegisterController.CreateSLARegister)
	complianceRoutes.PUT("/sla-register/:id", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARegisterController.UpdateSLARegister)
	complianceRoutes.POST("/sla-register/delete", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARegisterController.DeleteSLARegisters)
	complianceRoutes.PUT("/sla-register/archive", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARegisterController.ArchiveSLARegisters)
	complianceRoutes.PUT("/sla-register/restore", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARegisterController.RestoreSLARegisters)
	complianceRoutes.GET("/slas/:slaId/assessment-status-summary", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLARegisterController.AssessmentStatusSummary)

	// compliance clause
	complianceRoutes.GET("/compliances/:complianceId/clause", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.GetAllComplianceClauses)
	complianceRoutes.GET("/compliances/:complianceId/clause/:clauseId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.GetComplianceClauseDetails)
	complianceRoutes.POST("/compliances/:complianceId/clause", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.CreateComplianceClause)
	complianceRoutes.PUT("/compliances/:complianceId/clause/:clauseId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.UpdateComplianceClause)
	complianceRoutes.DELETE("/compliances/:complianceId/clause/:clauseId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.DeleteComplianceClause)
	complianceRoutes.POST("/compliances/:complianceId/clause/:clauseId/controls", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.CreateComplianceClauseControls)
	complianceRoutes.DELETE("/compliances/:complianceId/clause/:clauseId/controls/:controlId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.DeleteComplianceClauseControl)

	// compliance clause import
	complianceRoutes.GET("/compliances/:complianceId/clause/template", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.GetClausesTemplate)
	complianceRoutes.POST("/compliances/:complianceId/clause/template/validate", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.ValidateClauseTemplate)
	complianceRoutes.POST("/compliances/:complianceId/clause/import/validate", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.ValidateImportClause)
	complianceRoutes.POST("/compliances/:complianceId/clause/import", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.ImportClauses)

	complianceRoutes.GET("/compliances/:complianceId/clause/:clauseId/associated-requirements", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.GetAssociatedComplianceRequirement)
	complianceRoutes.GET("/compliances/:complianceId/clause/:clauseId/controls", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceClauseController.GetAssociatedControls)

	complianceRoutes.GET("/compliances/:complianceId/clause/:clauseId/risks", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceClauseController.GetAssociatedRisks)
	complianceRoutes.POST("/compliances/:complianceId/clause/:clauseId/risks", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.CreateComplianceClauseRisk)
	complianceRoutes.DELETE("/compliances/:complianceId/clause/:clauseId/risks/:riskId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.DeleteComplianceClauseRisk)

	complianceRoutes.GET("/compliances/:complianceId/clause/:clauseId/documents", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceClauseController.GetAssociatedDocuments)
	complianceRoutes.POST("/compliances/:complianceId/clause/:clauseId/documents", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.CreateComplianceClauseDocument)
	complianceRoutes.DELETE("/compliances/:complianceId/clause/:clauseId/documents/:documentId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.DeleteComplianceClauseDocument)

	complianceRoutes.GET("/compliances/:complianceId/clause/:clauseId/incidents", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceClauseController.GetAssociatedIncident)
	complianceRoutes.POST("/compliances/:complianceId/clause/:clauseId/incidents", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.CreateComplianceClauseIncident)
	complianceRoutes.DELETE("/compliances/:complianceId/clause/:clauseId/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.DeleteComplianceClauseIncident)

	complianceRoutes.GET("/compliances/:complianceId/clause/:clauseId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceClauseController.GetAssociatedStrategicObjective)
	complianceRoutes.POST("/compliances/:complianceId/clause/:clauseId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.CreateComplianceClauseStrategicObjective)
	complianceRoutes.DELETE("/compliances/:complianceId/clause/:clauseId/strategic-objectives/:strategicObjectiveId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.DeleteComplianceClausStrategicObjective)

	complianceRoutes.GET("/compliances/:complianceId/clause/:clauseId/assets", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceClauseController.GetAssociatedAssets)
	complianceRoutes.POST("/compliances/:complianceId/clause/:clauseId/assets", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.CreateComplianceClauseAssets)
	complianceRoutes.DELETE("/compliances/:complianceId/clause/:clauseId/assets/:assetId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.DeleteComplianceAssets)

	complianceRoutes.GET("/compliances/:complianceId/clause/:clauseId/processes", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceClauseController.GetAssociatedProcesses)
	complianceRoutes.POST("/compliances/:complianceId/clause/:clauseId/processes", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.CreateComplianceClausProcesses)
	complianceRoutes.DELETE("/compliances/:complianceId/clause/:clauseId/processes/:processId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.DeleteComplianceClausProcesses)

	complianceRoutes.GET("/compliances/:complianceId/clause/:clauseId/kpis", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceClauseController.GetAssociatedKPI)
	complianceRoutes.POST("/compliances/:complianceId/clause/:clauseId/kpis", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.CreateComplianceClausKPI)
	complianceRoutes.DELETE("/compliances/:complianceId/clause/:clauseId/kpis/:kpiId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.DeleteComplianceClausKPI)

	// SLA Renew
	complianceRoutes.GET("/slas/:slaId/renew-history", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLARenewController.GetSLARenewHistory)
	complianceRoutes.PUT("/slas/:slaId/renew", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARenewController.Renew)
	complianceRoutes.PUT("/slas/:slaId/revert", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARenewController.Revert)

	// Compliance Requirement
	complianceRoutes.GET("/compliance-requirements", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementController.GetAllComplianceRequirements)
	complianceRoutes.GET("/compliances/:complianceId/requirements/export", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementController.GetComplianceRequirementsExport)
	complianceRoutes.GET("/compliances/:complianceId/requirements", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementController.GetComplianceRequirementsByID)
	complianceRoutes.GET("/compliances/:complianceId/requirements/:requirementId", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementController.GetComplianceRequirementDetails)
	complianceRoutes.POST("/compliances/:complianceId/requirements", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementController.CreateComplianceRequirement)
	complianceRoutes.PUT("/compliances/:complianceId/requirements/:requirementId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementController.UpdateComplianceRequirement)
	complianceRoutes.POST("/compliances/:complianceId/requirements/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementController.DeleteComplianceRequirements)

	// Compliance Requirements Associated Items
	complianceRoutes.GET("/compliance-requirements/:requirementId/controls", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementAssociatedItemController.GetAllComplianceRequirementControls)
	complianceRoutes.GET("/compliance-requirements/:requirementId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementAssociatedItemController.GetAllComplianceRequirementStrategicObjectives)
	complianceRoutes.GET("/compliance-requirements/:requirementId/risks", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementAssociatedItemController.GetAllComplianceRequirementRisks)
	complianceRoutes.GET("/compliance-requirements/:requirementId/incidents", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementAssociatedItemController.GetAllComplianceRequirementIncidents)
	complianceRoutes.GET("/compliance-requirements/:requirementId/documents", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementAssociatedItemController.GetAllComplianceRequirementDocuments)

	// Compliance Action Plan
	complianceRoutes.GET("/compliance-action-plans/export", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceActionPlanController.GetComplianceActionPlanExport)
	complianceRoutes.GET("/compliance-action-plans", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceActionPlanController.GetAllComplianceActionPlans)
	complianceRoutes.GET("/compliances/:complianceId/action-plans", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceActionPlanController.GetAllComplianceActionPlans)
	complianceRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceActionPlanController.GetComplianceActionPlanDetails)
	complianceRoutes.POST("/compliances/:complianceId/action-plans", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceActionPlanController.CreateComplianceActionPlan)
	complianceRoutes.PUT("/compliances/:complianceId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceActionPlanController.UpdateComplianceActionPlan)
	complianceRoutes.POST("/compliances/:complianceId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceActionPlanController.DeleteComplianceActionPlans)

	// Compliance Status Update
	unauthenticatedRoutes.GET("/compliances/:complianceId/status-updates/:id/document-download", complianceStatusUpdateController.DownloadComplianceStatusUpdateDocument)
	unauthenticatedRoutes.GET("/compliances/:complianceId/status-updates/:id/document-preview", complianceStatusUpdateController.GetComplianceStatusUpdateDocumentPreview)
	complianceRoutes.GET("/compliances/:complianceId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceStatusUpdateController.GetAllComplianceStatusUpdates)
	complianceRoutes.GET("/compliances/:complianceId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceStatusUpdateController.GetComplianceStatusUpdateDetails)
	complianceRoutes.POST("/compliances/:complianceId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceStatusUpdateController.CreateComplianceStatusUpdate)
	complianceRoutes.PUT("/compliances/:complianceId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceStatusUpdateController.UpdateComplianceStatusUpdate)
	complianceRoutes.GET("/compliances/:complianceId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceStatusUpdateController.GetLatestComplianceStatusUpdate)

	// Compliance Requirement Status Update
	unauthenticatedRoutes.GET("/compliances/:complianceId/requirements/:requirementId/status-updates/:id/document-download", complianceRequirementStatusUpdateController.DownloadComplianceRequirementStatusUpdateDocument)
	unauthenticatedRoutes.GET("/compliances/:complianceId/requirements/:requirementId/status-updates/:id/document-preview", complianceRequirementStatusUpdateController.GetComplianceRequirementStatusUpdateDocumentPreview)
	complianceRoutes.GET("/compliances/:complianceId/requirements/:requirementId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementStatusUpdateController.GetAllComplianceRequirementStatusUpdates)
	complianceRoutes.POST("/compliances/:complianceId/requirements/:requirementId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementStatusUpdateController.CreateComplianceRequirementStatusUpdate)
	complianceRoutes.PUT("/compliances/:complianceId/requirements/:requirementId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementStatusUpdateController.UpdateComplianceRequirementStatusUpdate)
	complianceRoutes.GET("/compliances/:complianceId/requirements/:requirementId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementStatusUpdateController.GetLatestComplianceRequirementStatusUpdate)

	// Compliance Action Plan
	unauthenticatedRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/status-updates/:id/document-download", complianceActionPlanStatusUpdateController.DownloadComplianceActionPlanStatusUpdateDocument)
	unauthenticatedRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/status-updates/:id/document-preview", complianceActionPlanStatusUpdateController.GetComplianceActionPlanStatusUpdateDocumentPreview)
	complianceRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceActionPlanStatusUpdateController.GetAllComplianceActionPlanStatusUpdates)
	complianceRoutes.POST("/compliances/:complianceId/action-plans/:actionPlanId/status-updates", complianceActionPlanStatusUpdateController.CreateComplianceActionPlanStatusUpdate)
	complianceRoutes.PUT("/compliances/:complianceId/action-plans/:actionPlanId/status-updates/:id", complianceActionPlanStatusUpdateController.UpdateComplianceActionPlanStatusUpdate)
	complianceRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceActionPlanStatusUpdateController.GetLatestComplianceActionPlanStatusUpdate)
	complianceRoutes.PUT("/compliances/:complianceId/action-plans/:actionPlanId/approve", complianceActionPlanStatusUpdateController.ApproveComplianceActionPlan)
	complianceRoutes.PUT("/compliances/:complianceId/action-plans/:actionPlanId/revert", complianceActionPlanStatusUpdateController.RevertComplianceActionPlan)

	// Compliance Review
	complianceRoutes.GET("/compliances/:complianceId/reviews", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceReviewController.GetAllComplianceReviews)
	complianceRoutes.POST("/compliances/:complianceId/reviews", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceReviewController.CreateComplianceReview)
	complianceRoutes.PUT("/compliances/:complianceId/reviews/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceReviewController.UpdateComplianceReview)

	// SLA Check Question
	complianceRoutes.GET("/slas/:slaId/checks/:checkId/question-groups", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLACheckQuestionController.GetAllSLACheckQuestionGroups)
	complianceRoutes.GET("/slas/:slaId/checks/:checkId/question-groups/:groupId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLACheckQuestionController.GetSLACheckQuestionGroupDetails)
	complianceRoutes.POST("/slas/:slaId/checks/:checkId/question-groups", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.CreateSLACheckQuestionGroup)
	complianceRoutes.PUT("/slas/:slaId/checks/:checkId/question-groups/:groupId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.UpdateSLACheckQuestionGroup)
	complianceRoutes.DELETE("/slas/:slaId/checks/:checkId/question-groups/:groupId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.DeleteSLACheckQuestionGroup)

	complianceRoutes.POST("/slas/:slaId/checks/:checkId/choose-questions", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.ChooseSLACheckQuestions)
	complianceRoutes.GET("/slas/:slaId/checks/:checkId/questions/:id", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLACheckQuestionController.GetSLACheckQuestionDetails)

	complianceRoutes.GET("/slas/:slaId/checks/:checkId/question-groups/:groupId/questions", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLACheckQuestionController.GetSLACheckQuestionsByGroup)
	complianceRoutes.POST("/slas/:slaId/checks/:checkId/question-groups/:groupId/question", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.CreateSLACheckQuestion)
	complianceRoutes.PUT("/slas/:slaId/checks/:checkId/question-groups/:groupId/questions/:id", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.UpdateSLACheckQuestion)
	complianceRoutes.DELETE("/slas/:slaId/checks/:checkId/question-groups/:groupId/questions/:id", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.DeleteSLACheckQuestion)

	// SLA Answer
	complianceRoutes.POST("/slas/:slaId/questions/:questionId/answers", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckAnswerController.CreateSLACheckAnswer)
	complianceRoutes.PUT("/slas/:slaId/questions/:questionId/answers/:answerId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckAnswerController.UpdateSLACheckAnswer)
	complianceRoutes.DELETE("/slas/:slaId/questions/:questionId/answers/:answerId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckAnswerController.DeleteSLACheckAnswer)
	complianceRoutes.GET("/slas/:slaId/assessment-summary", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLACheckAnswerController.AssessmentSummary)

	// SLA Action plan
	complianceRoutes.GET("/slas/:slaId/action-plans", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaActionPlanController.GetAllSLAActionPlans)
	complianceRoutes.GET("/slas/:slaId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaActionPlanController.GetSLAActionPlanDetails)
	complianceRoutes.POST("/slas/:slaId/action-plans", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), slaActionPlanController.CreateSLAActionPlan)
	complianceRoutes.PUT("/slas/:slaId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), slaActionPlanController.UpdateSLAActionPlan)
	complianceRoutes.POST("/slas/:slaId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), slaActionPlanController.DeleteSLAActionPlans)
	complianceRoutes.GET("/slas/:slaId/action-plan-summary", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaActionPlanController.GetSLAActionPlanCountByStatus)

	// SLA Action Plan
	unauthenticatedRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/status-updates/:id/document-download", slaActionPlanStatusUpdateController.DownloadSLAActionPlanStatusUpdateDocument)
	unauthenticatedRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/status-updates/:id/document-preview", slaActionPlanStatusUpdateController.GetSLAActionPlanStatusUpdateDocumentPreview)
	complianceRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaActionPlanStatusUpdateController.GetAllSLAActionPlanStatusUpdates)
	complianceRoutes.POST("/slas/:slaId/action-plans/:actionPlanId/status-updates", slaActionPlanStatusUpdateController.CreateSLAActionPlanStatusUpdate)
	complianceRoutes.PUT("/slas/:slaId/action-plans/:actionPlanId/status-updates/:id", slaActionPlanStatusUpdateController.UpdateSLAActionPlanStatusUpdate)
	complianceRoutes.PUT("/slas/:slaId/action-plans/:actionPlanId/approve", slaActionPlanStatusUpdateController.ApproveSLAActionPlan)
	complianceRoutes.PUT("/slas/:slaId/action-plans/:actionPlanId/revert", slaActionPlanStatusUpdateController.RevertSLAActionPlan)
	complianceRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaActionPlanStatusUpdateController.GetLatestSLAActionPlanStatusUpdate)

	// SLA Action Plan Comments
	unauthenticatedRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/comments/:commentId/document-download", slaActionPlanCommentController.DownloadSLAActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/comments/:commentId/document-preview", slaActionPlanCommentController.GetSLAActionPlanCommentDocumentPreview)
	complianceRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaActionPlanCommentController.GetAllComments)
	complianceRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaActionPlanCommentController.GetCommentDetails)
	complianceRoutes.POST("/slas/:slaId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), slaActionPlanCommentController.CreateComment)
	complianceRoutes.PUT("/slas/:slaId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), slaActionPlanCommentController.UpdateComment)
	complianceRoutes.POST("/slas/:slaId/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), slaActionPlanCommentController.DeleteComment)

	complianceRoutes.GET("/slas/:slaId/processes", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaAssociatedItemController.GetAllSLAProcesses)
	complianceRoutes.GET("/slas/:slaId/documents", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaAssociatedItemController.GetAllSLADocuments)
	complianceRoutes.GET("/slas/:slaId/risks", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaAssociatedItemController.GetAllSLARisks)
	complianceRoutes.GET("/slas/:slaId/incidents", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaAssociatedItemController.GetAllSLAIncidents)
	complianceRoutes.GET("/slas/:slaId/associated-item-summary", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaAssociatedItemController.GetAllAssociatedCounts)

	// Compliance Action Plan Comments
	unauthenticatedRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/comments/:commentId/document-download", ComplianceActionPlanCommentController.DownloadComplianceActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/comments/:commentId/document-preview", ComplianceActionPlanCommentController.GetComplianceActionPlanCommentDocumentPreview)
	complianceRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), ComplianceActionPlanCommentController.GetAllComments)
	complianceRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), ComplianceActionPlanCommentController.GetCommentDetails)
	complianceRoutes.POST("/compliances/:complianceId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), ComplianceActionPlanCommentController.CreateComment)
	complianceRoutes.PUT("/compliances/:complianceId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), ComplianceActionPlanCommentController.UpdateComment)
	complianceRoutes.POST("/compliances/:complianceId/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), ComplianceActionPlanCommentController.DeleteComment)

	// Compliance Requirement Assessment
	complianceRoutes.GET("/compliance-assessments", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentController.GetAllComplianceAssessments)
	complianceRoutes.GET("/compliance-assessments/export", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentController.GetComplianceAssessmentsExport)
	complianceRoutes.GET("/compliance-assessments/:assessmentId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentController.GetComplianceAssessmentDetails)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/is-all-requirements-approved", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentController.GetApprovalStatus)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/is-all-controls-approved", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentController.GetControlsApprovalStatus)
	complianceRoutes.POST("/compliance-assessments", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentController.CreateComplianceAssessment)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentController.UpdateComplianceAssessment)
	complianceRoutes.POST("/compliance-assessments/force-delete", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentController.DeleteForceComplianceAssessments)
	complianceRoutes.POST("/compliance-assessments/delete", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentController.DeleteComplianceAssessments)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/completion-summary", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentController.GetComplianceAssessmentSummary)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/compliance-status-summary", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentController.GetComplianceStatusSummary)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/compliance-score-summary", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentController.GetComplianceScoreSummary)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/submit", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentController.SubmitComplianceAssessment)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/approve", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentController.ApproveComplianceAssessment)

	complianceRoutes.GET("/compliance-assessments/:assessmentId/clauses", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentRequirementController.GetComplianceAssessmentRequirementClauses)

	complianceRoutes.GET("/compliance-assessments/:assessmentId/requirements", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentRequirementController.GetComplianceAssessmentRequirements)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentRequirementController.GetComplianceAssessmentRequirementDetails)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/requirements/:requirementId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementController.UpdateComplianceAssessmentRequirement)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/requirements/:requirementId/mark-as-completed", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementController.MarkAsCompleted)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/requirements/:requirementId/mark-as-open", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementController.MarkAsOpen)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/requirements/:requirementId/approve", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementController.ApproveComplianceAssessmentRequirement)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/requirements/:requirementId/revert", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementController.RevertComplianceAssessmentRequirement)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/requirements/:requirementId/resubmit", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementController.ResubmitComplianceAssessmentRequirement)

	// Compliance Requirement Assessment Documents
	unauthenticatedRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/required-documents/:documentId/document-download", complianceAssessmentRequirementController.DownloadComplianceAssessmentRequirementDocument)
	unauthenticatedRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/required-documents/:documentId/document-preview", complianceAssessmentRequirementController.GetComplianceAssessmentRequirementDocumentPreview)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/requirements/:requirementId/required-documents/:documentId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementController.CreateDocuments)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/requirements/:requirementId/required-documents/:documentId/delete", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementController.DeleteDocuments)

	// Compliance Assessment Requirement Comment
	complianceRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/comments", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentRequirementCommentController.GetAllComments)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentRequirementCommentController.GetCommentDetails)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/requirements/:requirementId/comments", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementCommentController.CreateComment)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/requirements/:requirementId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementCommentController.UpdateComment)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/requirements/:requirementId/comments/delete", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementCommentController.DeleteComment)

	// Compliance Assessment Requirement Action Plan
	complianceRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentRequirementActionPlanController.GetAllComplianceAssessmentRequirementActionPlans)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentRequirementActionPlanController.GetComplianceAssessmentRequirementActionPlanDetails)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementActionPlanController.CreateComplianceAssessmentRequirementActionPlan)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementActionPlanController.UpdateComplianceAssessmentRequirementActionPlan)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementActionPlanController.DeleteComplianceAssessmentRequirementActionPlans)

	// Compliance Action Plan
	unauthenticatedRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/status-updates/:id/document-download", complianceAssessmentRequirementActionPlanStatusUpdateController.DownloadComplianceAssessmentRequirementActionPlanStatusUpdateDocument)
	unauthenticatedRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/status-updates/:id/document-preview", complianceAssessmentRequirementActionPlanStatusUpdateController.GetComplianceAssessmentRequirementActionPlanStatusUpdateDocumentPreview)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentRequirementActionPlanStatusUpdateController.GetAllComplianceAssessmentRequirementActionPlanStatusUpdates)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/status-updates", complianceAssessmentRequirementActionPlanStatusUpdateController.CreateComplianceAssessmentRequirementActionPlanStatusUpdate)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/status-updates/:id", complianceAssessmentRequirementActionPlanStatusUpdateController.UpdateComplianceAssessmentRequirementActionPlanStatusUpdate)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentRequirementActionPlanStatusUpdateController.GetLatestComplianceAssessmentRequirementActionPlanStatusUpdate)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/approve", complianceAssessmentRequirementActionPlanStatusUpdateController.ApproveComplianceAssessmentRequirementActionPlan)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/revert", complianceAssessmentRequirementActionPlanStatusUpdateController.RevertComplianceAssessmentRequirementActionPlan)

	complianceRoutes.GET("/compliance-assessment-statuses", complianceAssessmentStatusController.GetAllComplianceAssessmentStatuses)
	complianceRoutes.GET("/compliance-assessment-statuses/:id", complianceAssessmentStatusController.GetComplianceAssessmentStatusDetails)

	complianceRoutes.GET("/compliance-assessment-requirement-statuses", complianceAssessmentRequirementStatusController.GetAllComplianceAssessmentRequirementStatuses)
	complianceRoutes.GET("/compliance-assessment-requirement-statuses/:id", complianceAssessmentRequirementStatusController.GetComplianceAssessmentRequirementStatusDetails)

	// Compliance Assessment Requirment Action Plan Comments
	unauthenticatedRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/comments/:commentId/document-download", complianceAssessmentRequirementActionPlanCommentController.DownloadComplianceAssessmentRequirementActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/comments/:commentId/document-preview", complianceAssessmentRequirementActionPlanCommentController.GetComplianceAssessmentRequirementActionPlanCommentDocumentPreview)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentRequirementActionPlanCommentController.GetAllComments)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentRequirementActionPlanCommentController.GetCommentDetails)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementActionPlanCommentController.CreateComment)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementActionPlanCommentController.UpdateComment)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/requirements/:requirementId/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentRequirementActionPlanCommentController.DeleteComment)

	// Compliance Report
	complianceRoutes.GET("/compliance-reports", complianceReportController.GetAllComplianceReports)
	complianceRoutes.POST("/compliance-reports", complianceReportController.CreateComplianceReport)
	complianceRoutes.GET("/compliance-reports/:complianceReportId", complianceReportController.GetComplianceReportDetails)
	complianceRoutes.GET("/compliance-reports/:complianceReportId/export-pdf", complianceReportController.ExportComplianceReportToPDF)
	complianceRoutes.DELETE("/compliance-reports/:complianceReportId", complianceReportController.DeleteComplianceReport)
	complianceRoutes.GET("/compliance-reports/:complianceReportId/sections/:sectionId", complianceReportController.GetComplianceReportSectionDetails)
	complianceRoutes.PUT("/compliance-reports/:complianceReportId/sections/:sectionId", complianceReportController.UpdateComplianceReportSection)

	complianceRoutes.GET("/policy-exceptions", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), policyExceptionController.GetAllPolicyExceptions)
	complianceRoutes.GET("/policy-exceptions/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), policyExceptionController.GetPolicyExceptionDetails)
	complianceRoutes.POST("/policy-exceptions", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), policyExceptionController.CreatePolicyException)
	complianceRoutes.PUT("/policy-exceptions/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), policyExceptionController.UpdatePolicyException)
	complianceRoutes.POST("/policy-exceptions/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), policyExceptionController.DeletePolicyExceptions)
	complianceRoutes.POST("/policy-exceptions/:id/risks", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), policyExceptionController.MapRisksToPolicyException)
	complianceRoutes.DELETE("/policy-exceptions/:id/risks/:riskID", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), policyExceptionController.DeleteRisksFromPolicyException)

	// Compliance Assessment Framework Type
	complianceRoutes.GET("/compliance-assessment-framework-types", complianceAssessmentFrameworkTypeController.GetAllComplianceAssessmentFrameworkTypes)
	complianceRoutes.GET("/compliance-assessment-framework-types/:id", complianceAssessmentFrameworkTypeController.GetComplianceAssessmentFrameworkTypeDetails)

	// Compliance Assessment Framework
	complianceRoutes.GET("/compliance-assessment-frameworks", complianceAssessmentFrameworkController.GetAllComplianceAssessmentFrameworks)
	complianceRoutes.GET("/compliance-assessment-frameworks/:complianceAssessmentFrameworkId", complianceAssessmentFrameworkController.GetComplianceAssessmentFrameworkDetails)
	complianceRoutes.POST("/compliance-assessment-frameworks", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceAssessmentFrameworkController.CreateComplianceAssessmentFramework)
	complianceRoutes.PUT("/compliance-assessment-frameworks/:complianceAssessmentFrameworkId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceAssessmentFrameworkController.UpdateComplianceAssessmentFramework)
	complianceRoutes.POST("/compliance-assessment-frameworks/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceAssessmentFrameworkController.DeleteComplianceAssessmentFramework)

	// Compliance Assessment Framework Values
	complianceRoutes.GET("/compliance-assessment-frameworks/:complianceAssessmentFrameworkId/values", complianceAssessmentFrameworkValueController.GetAllComplianceAssessmentFrameworkValues)
	complianceRoutes.GET("/compliance-assessment-frameworks/:complianceAssessmentFrameworkId/values/:complianceAssessmentFrameworkValueId", complianceAssessmentFrameworkValueController.GetComplianceAssessmentFrameworkValueDetails)
	complianceRoutes.POST("/compliance-assessment-frameworks/:complianceAssessmentFrameworkId/values", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceAssessmentFrameworkValueController.CreateComplianceAssessmentFrameworkValue)
	complianceRoutes.PUT("/compliance-assessment-frameworks/:complianceAssessmentFrameworkId/values/:complianceAssessmentFrameworkValueId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceAssessmentFrameworkValueController.UpdateComplianceAssessmentFrameworkValue)
	complianceRoutes.POST("/compliance-assessment-frameworks/:complianceAssessmentFrameworkId/values/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceAssessmentFrameworkValueController.DeleteComplianceAssessmentFrameworkValues)

	complianceRoutes.GET("/compliances/:complianceId/compliance-requirements/template", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementController.GetRequirementsTemplate)
	complianceRoutes.POST("/compliances/:complianceId/compliance-requirements/template/validate", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementController.ValidateRequirementTemplate)
	complianceRoutes.POST("/compliances/:complianceId/compliance-requirements/import/validate", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementController.ValidateImportRequirement)
	complianceRoutes.POST("/compliances/:complianceId/compliance-requirements/import", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementController.ImportRequirements)
	// CLAUSES
	complianceRoutes.GET("/compliance-assessments/:assessmentId/control-clauses", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentControlController.GetComplianceAssessmentControlClauses)
	// CONTROLS LIST + DETAILS + UPDATE
	complianceRoutes.GET("/compliance-assessments/:assessmentId/controls", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentControlController.GetComplianceAssessmentControls)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentControlController.GetComplianceAssessmentControlDetails)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/controls/:controlId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlController.UpdateComplianceAssessmentControl)
	// MARK COMPLETED / OPEN
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/controls/:controlId/mark-as-completed", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlController.MarkAsCompleted)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/controls/:controlId/mark-as-open", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlController.MarkAsOpen)
	// WORKFLOW: APPROVE / REVERT / RESUBMIT
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/controls/:controlId/approve", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlController.ApproveComplianceAssessmentControl)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/controls/:controlId/revert", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlController.RevertComplianceAssessmentControl)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/controls/:controlId/resubmit", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlController.ResubmitComplianceAssessmentControl)

	// Compliance Assessment CONTROLS Comment
	complianceRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId/comments", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentControlCommentController.GetAllComments)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentControlCommentController.GetCommentDetails)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/controls/:controlId/comments", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlCommentController.CreateComment)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/controls/:controlId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlCommentController.UpdateComment)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/controls/:controlId/comments/delete", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlCommentController.DeleteComment)
	// Compliance Assessment Control Action Plan
	complianceRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId/action-plans", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentControlActionPlanController.GetAllComplianceAssessmentControlActionPlans)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentControlActionPlanController.GetComplianceAssessmentControlActionPlanDetails)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/controls/:controlId/action-plans", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlActionPlanController.CreateComplianceAssessmentControlActionPlan)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlActionPlanController.UpdateComplianceAssessmentControlActionPlan)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlActionPlanController.DeleteComplianceAssessmentControlActionPlans)

	unauthenticatedRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/status-updates/:id/document-download", complianceAssessmentControlActionPlanStatusUpdateController.DownloadComplianceAssessmentControlActionPlanStatusUpdateDocument)
	unauthenticatedRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/status-updates/:id/document-preview", complianceAssessmentControlActionPlanStatusUpdateController.GetComplianceAssessmentControlActionPlanStatusUpdateDocumentPreview)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentControlActionPlanStatusUpdateController.GetAllComplianceAssessmentControlActionPlanStatusUpdates)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/status-updates", complianceAssessmentControlActionPlanStatusUpdateController.CreateComplianceAssessmentControlActionPlanStatusUpdate)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/status-updates/:id", complianceAssessmentControlActionPlanStatusUpdateController.UpdateComplianceAssessmentControlActionPlanStatusUpdate)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentControlActionPlanStatusUpdateController.GetLatestComplianceAssessmentControlActionPlanStatusUpdate)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/approve", complianceAssessmentControlActionPlanStatusUpdateController.ApproveComplianceAssessmentControlActionPlan)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/revert", complianceAssessmentControlActionPlanStatusUpdateController.RevertComplianceAssessmentControlActionPlan)

	// Compliance Assessment Control Action Plan Comments
	unauthenticatedRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/comments/:commentId/document-download", complianceAssessmentControlActionPlanCommentController.DownloadComplianceAssessmentControlActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/comments/:commentId/document-preview", complianceAssessmentControlActionPlanCommentController.GetComplianceAssessmentControlActionPlanCommentDocumentPreview)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentControlActionPlanCommentController.GetAllComments)
	complianceRoutes.GET("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "view"), complianceAssessmentControlActionPlanCommentController.GetCommentDetails)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlActionPlanCommentController.CreateComment)
	complianceRoutes.PUT("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlActionPlanCommentController.UpdateComment)
	complianceRoutes.POST("/compliance-assessments/:assessmentId/controls/:controlId/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "compliance_assessment", "modify"), complianceAssessmentControlActionPlanCommentController.DeleteComment)

	// Compliance Clause Control Maping Import
	complianceRoutes.GET("/compliances/:complianceId/controls/template", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.GetControlMappingTemplate)
	complianceRoutes.POST("/compliances/:complianceId/controls/template/validate", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.ValidateControlMappingTemplate)
	complianceRoutes.POST("/compliances/:complianceId/controls/import/validate", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.ValidateImportControlMapping)
	complianceRoutes.POST("/compliances/:complianceId/controls/import", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceClauseController.ImportControlMapping)

	// Compliance Detail Dashboard
	complianceRoutes.GET("/compliances/:complianceId/dashboard/summary/compliances", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetComplianceDetailSummary)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/summary/clauses", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetComplianceClauseSummary)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/overview/top-5-clauses", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetTop5ClausePerformance)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/overview/bottom-5-clauses", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetBottom5ClausePerformance)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/overview/domains", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetComplianceDomainWiseOverview)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/summary/policy-organization-units", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetPolicyByOrganizationUnit)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/summary/control-effectiveness", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetControlEffectivenessSummary)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/summary/control-implementations", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetControlImplementationSummary)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/summary/controls", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetControlSummary)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/overview/control-effectiveness", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetControlEfficiencySummary)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/summary/domain-controls", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetControlCoverageByDomain)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/summary/policies", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetCompliancePolicySummary)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/summary/policy-acknowledgements", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetPolicyAcknowledgementByOrganizationUnit)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/overview/action-plans", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetActionPlanOverview)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/summary/action-plan-statuses", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetActionPlanStatusSummary)
	complianceRoutes.GET("/compliances/:complianceId/dashboard/summary/action-plan-organization-units", middlewares.RolePermissionMiddleware(db, "compliance_dashboard", "view"), complianceDetailDashboardController.GetActionPlansByOrganizationUnit)
}

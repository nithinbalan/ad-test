package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	auditRoutes := router.Group("/")

	auditRoutes.Use(middlewares.AuthMiddleware())

	auditStatusController := app.AuditModule.AuditStatusController
	auditorRiskRatingController := app.AuditModule.AuditorRiskRatingController
	auditUniverseController := app.AuditModule.AuditUniverseController
	auditableItemController := app.AuditModule.AuditableItemController
	auditTeamController := app.AuditModule.AuditTeamController
	auditCategoryController := app.AuditModule.AuditCategoryController
	auditPlanController := app.AuditModule.AuditPlanController
	auditPlanChecklistController := app.AuditModule.AuditPlanChecklistController
	auditPlanToolController := app.AuditModule.AuditPlanToolController
	auditPlanAuditableItemController := app.AuditModule.AuditPlanAuditableItemController
	auditActivityTypeController := app.AuditModule.AuditActivityTypeController
	auditScheduleStatusController := app.AuditModule.AuditScheduleStatusController
	auditPlanBudgetController := app.AuditModule.AuditPlanBudgetController
	auditPlanScheduleController := app.AuditModule.AuditPlanScheduleController

	//Audit Status
	auditRoutes.GET("/audit-statuses", auditStatusController.GetAllAuditStatuses)
	auditRoutes.GET("/audit-statuses/:id", auditStatusController.GetAuditStatusDetails)

	//Auditor Risk Rating
	auditRoutes.GET("/auditor-risk-ratings", auditorRiskRatingController.GetAllAuditorRiskRatings)
	auditRoutes.GET("/auditor-risk-ratings/:id", auditorRiskRatingController.GetAuditorRiskRatingDetails)

	//Audit Universe
	auditRoutes.POST("/audit-universe/mark-as-auditable-item", auditUniverseController.MarkAsAuditableItem)
	auditRoutes.POST("/audit-universe/unmark-as-auditable-item", auditUniverseController.UnmarkAsAuditableItem)

	//Auditable Item
	auditRoutes.GET("/auditable-items", auditableItemController.GetAllAuditableItems)
	auditRoutes.GET("/auditable-items/:id", auditableItemController.GetAuditableItemDetails)
	auditRoutes.POST("/auditable-items", auditableItemController.CreateAuditableItem)
	auditRoutes.PUT("/auditable-items/:id", auditableItemController.UpdateAuditableItem)
	auditRoutes.POST("/auditable-items/delete", auditableItemController.DeleteAuditableItems)

	//Audit Team
	auditRoutes.GET("/audit-teams", auditTeamController.GetAllAuditTeams)
	auditRoutes.GET("/audit-teams/:id", auditTeamController.GetAuditTeamDetails)
	auditRoutes.POST("/audit-teams", auditTeamController.CreateAuditTeam)
	auditRoutes.PUT("/audit-teams/:id", auditTeamController.UpdateAuditTeam)
	auditRoutes.POST("/audit-teams/delete", auditTeamController.DeleteAuditTeams)

	//Audit Category
	auditRoutes.GET("/audit-categories", auditCategoryController.GetAllAuditCategories)
	auditRoutes.GET("/audit-categories/:id", auditCategoryController.GetAuditCategoryDetails)
	auditRoutes.POST("/audit-categories", auditCategoryController.CreateAuditCategory)
	auditRoutes.PUT("/audit-categories/:id", auditCategoryController.UpdateAuditCategory)
	auditRoutes.POST("/audit-categories/delete", auditCategoryController.DeleteAuditCategories)
	auditRoutes.PUT("/audit-categories/activate", auditCategoryController.ActivateAuditCategories)
	auditRoutes.PUT("/audit-categories/deactivate", auditCategoryController.DeactivateAuditCategories)

	//AuditPlan
	auditRoutes.GET("/audit-plans", auditPlanController.GetAllAuditPlans)
	auditRoutes.GET("/audit-plans/:planId", auditPlanController.GetAuditPlanDetails)
	auditRoutes.POST("/audit-plans", auditPlanController.CreateAuditPlan)
	auditRoutes.PUT("/audit-plans/:planId", auditPlanController.UpdateAuditPlan)
	auditRoutes.POST("/audit-plans/delete", auditPlanController.DeleteAuditPlans)

	//AuditPlan Checklist
	auditRoutes.GET("/audit-plans/:planId/checklists", auditPlanChecklistController.GetAllAuditPlanChecklists)
	auditRoutes.POST("/audit-plans/:planId/checklists", auditPlanChecklistController.CreateAuditPlanChecklists)
	auditRoutes.PUT("/audit-plans/:planId/checklists/:id", auditPlanChecklistController.UpdateAuditPlanChecklist)
	auditRoutes.POST("/audit-plans/:planId/checklists/delete", auditPlanChecklistController.DeleteAuditPlanChecklists)

	//AuditPlan Tools
	auditRoutes.GET("/audit-plans/:planId/tools", auditPlanToolController.GetAllAuditPlanTools)
	auditRoutes.GET("/audit-plans/:planId/tools/:id", auditPlanToolController.GetAuditPlanToolDetails)
	auditRoutes.POST("/audit-plans/:planId/tools", auditPlanToolController.CreateAuditPlanTool)
	auditRoutes.PUT("/audit-plans/:planId/tools/:id", auditPlanToolController.UpdateAuditPlanTool)
	auditRoutes.POST("/audit-plans/:planId/tools/delete", auditPlanToolController.DeleteAuditPlanTools)

	//AuditPlan Auditable Item
	auditRoutes.GET("/audit-plans/:planId/auditable-items", auditPlanAuditableItemController.GetAllAuditPlanAuditableItems)
	auditRoutes.POST("/audit-plans/:planId/auditable-items", auditPlanAuditableItemController.CreateAuditPlanAuditableItem)
	auditRoutes.POST("/audit-plans/:planId/auditable-items/delete", auditPlanAuditableItemController.DeleteAuditPlanAuditableItems)

	//Audit Activity Type
	auditRoutes.GET("/audit-activity-types", auditActivityTypeController.GetAllAuditActivityTypes)
	auditRoutes.GET("/audit-activity-types/:id", auditActivityTypeController.GetAuditActivityTypeDetails)
	auditRoutes.POST("/audit-activity-types", auditActivityTypeController.CreateAuditActivityType)
	auditRoutes.PUT("/audit-activity-types/:id", auditActivityTypeController.UpdateAuditActivityType)
	auditRoutes.POST("/audit-activity-types/delete", auditActivityTypeController.DeleteAuditActivityTypes)
	auditRoutes.PUT("/audit-activity-types/activate", auditActivityTypeController.ActivateAuditActivityTypes)
	auditRoutes.PUT("/audit-activity-types/deactivate", auditActivityTypeController.DeactivateAuditActivityTypes)

	//Audit Schedule Statuses
	auditRoutes.GET("/audit-schedule-statuses", auditScheduleStatusController.GetAllAuditScheduleStatuses)
	auditRoutes.GET("/audit-schedule-statuses/:id", auditScheduleStatusController.GetAuditScheduleStatusDetails)

	//AuditPlan Tools
	auditRoutes.GET("/audit-plans/:planId/budgets", auditPlanBudgetController.GetAllAuditPlanBudgets)
	auditRoutes.GET("/audit-plans/:planId/budgets/:id", auditPlanBudgetController.GetAuditPlanBudgetDetails)
	auditRoutes.POST("/audit-plans/:planId/budgets", auditPlanBudgetController.CreateAuditPlanBudget)
	auditRoutes.PUT("/audit-plans/:planId/budgets/:id", auditPlanBudgetController.UpdateAuditPlanBudget)
	auditRoutes.POST("/audit-plans/:planId/budgets/delete", auditPlanBudgetController.DeleteAuditPlanBudgets)

	//Audit Schedule
	auditRoutes.GET("/audit-plans/:planId/schedules", auditPlanScheduleController.GetAllAuditPlanSchedules)
	auditRoutes.GET("/audit-plans/:planId/schedules/:id", auditPlanScheduleController.GetAuditPlanScheduleDetails)
	auditRoutes.POST("/audit-plans/:planId/schedules", auditPlanScheduleController.CreateAuditPlanSchedule)
	auditRoutes.PUT("/audit-plans/:planId/schedules/:id", auditPlanScheduleController.UpdateAuditPlanSchedule)
	auditRoutes.POST("/audit-plans/:planId/schedules/delete", auditPlanScheduleController.DeleteAuditPlanSchedules)
}

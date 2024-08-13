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
	auditRoutes.GET("/audit-plans/:id", auditPlanController.GetAuditPlanDetails)
	auditRoutes.POST("/audit-plans", auditPlanController.CreateAuditPlan)
	auditRoutes.PUT("/audit-plans/:id", auditPlanController.UpdateAuditPlan)
	auditRoutes.POST("/audit-plans/delete", auditPlanController.DeleteAuditPlans)

	//AuditPlan Checklist
	auditRoutes.GET("/audit-plans/:planId/checklists", auditPlanChecklistController.GetAllAuditPlanChecklists)
	auditRoutes.POST("/audit-plans/:planId/checklists", auditPlanChecklistController.CreateAuditPlanChecklists)
	auditRoutes.PUT("/audit-plans/:planId/checklists/:id", auditPlanChecklistController.UpdateAuditPlanChecklist)
	auditRoutes.POST("/audit-plans/:planId/checklists/delete", auditPlanChecklistController.DeleteAuditPlanChecklists)

}

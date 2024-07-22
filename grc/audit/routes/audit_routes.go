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
	auditRiskRatingController := app.AuditModule.AuditRiskRatingController

	//Audit Status
	auditRoutes.GET("/audit-statuses", auditStatusController.GetAllAuditStatuses)
	auditRoutes.GET("/audit-statuses/:id", auditStatusController.GetAuditStatusDetails)

	//Audit Risk Rating
	auditRoutes.GET("/audit-risk-ratings", auditRiskRatingController.GetAllAuditRiskRatings)
	auditRoutes.GET("/audit-risk-ratings/:id", auditRiskRatingController.GetAuditRiskRatingDetails)

}

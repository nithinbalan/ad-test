package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	tprmRoutes := router.Group("/")

	tprmRoutes.Use(middlewares.AuthMiddleware())

	db := app.DB

	thirdPartyTypeController := app.TprmModule.ThirdPartyTypeController

	tprmRoutes.GET("/third-party-types", thirdPartyTypeController.GetAllThirdPartyTypes)
	tprmRoutes.GET("/third-party-types/:id", thirdPartyTypeController.GetThirdPartyTypeDetails)
	tprmRoutes.POST("/third-party-types", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.CreateThirdPartyType)
	tprmRoutes.PUT("/third-party-types/:id", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.UpdateThirdPartyType)
	tprmRoutes.POST("/third-party-types/delete", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.DeleteThirdPartyTypes)
	tprmRoutes.PUT("/third-party-types/activate", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.ActivateThirdPartyTypes)
	tprmRoutes.PUT("/third-party-types/deactivate", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.DeactivateThirdPartyTypes)

}

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
	thirdPartyRegisterController := app.TprmModule.ThirdPartyRegisterController
	thirdPartyQuestionCategoryController := app.TprmModule.ThirdPartyQuestionCategoryController
	thirdPartyQuestionTypeController := app.TprmModule.ThirdPartyQuestionTypeController
	thirdPartyQuestionLibraryController := app.TprmModule.ThirdPartyQuestionLibraryController

	tprmRoutes.GET("/third-party-types", thirdPartyTypeController.GetAllThirdPartyTypes)
	tprmRoutes.GET("/third-party-types/:id", thirdPartyTypeController.GetThirdPartyTypeDetails)
	tprmRoutes.POST("/third-party-types", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.CreateThirdPartyType)
	tprmRoutes.PUT("/third-party-types/:id", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.UpdateThirdPartyType)
	tprmRoutes.POST("/third-party-types/delete", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.DeleteThirdPartyTypes)
	tprmRoutes.PUT("/third-party-types/activate", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.ActivateThirdPartyTypes)
	tprmRoutes.PUT("/third-party-types/deactivate", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.DeactivateThirdPartyTypes)

	// === Third Party Register ===
	tprmRoutes.GET("/third-party-register", thirdPartyRegisterController.GetAllThirdPartyRegister)
	tprmRoutes.GET("/third-party-register/export", thirdPartyRegisterController.GetThirdPartyRegisterExport)
	tprmRoutes.GET("/third-party-register/:id", thirdPartyRegisterController.GetThirdPartyRegisterDetails)
	tprmRoutes.POST("/third-party-register", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyRegisterController.CreateThirdPartyRegister)
	tprmRoutes.PUT("/third-party-register/:id", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyRegisterController.UpdateThirdPartyRegister)
	tprmRoutes.POST("/third-party-register/delete", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyRegisterController.DeleteThirdPartyRegister)

	tprmRoutes.GET("/third-party-question-categories", thirdPartyQuestionCategoryController.GetAllThirdPartyQuestionCategories)
	tprmRoutes.GET("/third-party-question-categories/:id", thirdPartyQuestionCategoryController.GetThirdPartyQuestionCategoryDetails)
	tprmRoutes.POST("/third-party-question-categories", thirdPartyQuestionCategoryController.CreateThirdPartyQuestionCategory)
	tprmRoutes.PUT("/third-party-question-categories/:id", thirdPartyQuestionCategoryController.UpdateThirdPartyQuestionCategory)
	tprmRoutes.POST("/third-party-question-categories/delete", thirdPartyQuestionCategoryController.DeleteThirdPartyQuestionCategories)
	tprmRoutes.PUT("/third-party-question-categories/activate", thirdPartyQuestionCategoryController.ActivateThirdPartyQuestionCategories)
	tprmRoutes.PUT("/third-party-question-categories/deactivate", thirdPartyQuestionCategoryController.DeactivateThirdPartyQuestionCategories)

	tprmRoutes.GET("/third-party-question-types", thirdPartyQuestionTypeController.GetAllThirdPartyQuestionTypes)
	tprmRoutes.GET("/third-party-question-types/:id", thirdPartyQuestionTypeController.GetThirdPartyQuestionTypeDetails)

	// === Third Party Question Library ===
	tprmRoutes.GET("/third-party-question-libraries", thirdPartyQuestionLibraryController.GetAllThirdPartyQuestionLibraries)
	tprmRoutes.GET("/third-party-question-libraries/:id", thirdPartyQuestionLibraryController.GetThirdPartyQuestionLibraryDetails)
	tprmRoutes.POST("/third-party-question-libraries", thirdPartyQuestionLibraryController.CreateThirdPartyQuestionLibrary)
	tprmRoutes.PUT("/third-party-question-libraries/:id", thirdPartyQuestionLibraryController.UpdateThirdPartyQuestionLibrary)
	tprmRoutes.POST("/third-party-question-libraries/delete", thirdPartyQuestionLibraryController.DeleteThirdPartyQuestionLibraries)

}

package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	documentRoutes := router.Group("/")

	documentRoutes.Use(middlewares.AuthMiddleware())

	documentCategoryController := app.DocumentModule.DocumentCategoryController

	documentRoutes.GET("/document-categories", documentCategoryController.GetAllDocumentCategories)
	documentRoutes.GET("/document-categories/:id", documentCategoryController.GetDocumentCategoryDetails)
	documentRoutes.POST("/document-categories", documentCategoryController.CreateDocumentCategory)
	documentRoutes.PUT("/document-categories/:id", documentCategoryController.UpdateDocumentCategory)
	documentRoutes.POST("/document-categories/delete", documentCategoryController.DeleteDocumentCategories)
	documentRoutes.PUT("/document-categories/activate", documentCategoryController.ActivateDocumentCategories)
	documentRoutes.PUT("/document-categories/deactivate", documentCategoryController.DeactivateDocumentCategories)

}

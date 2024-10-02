package controllers

import (
	"jethur/grc/documents/services"
)

type DocumentWorkflowController struct {
	service services.DocumentWorkflowServiceInterface
}

func NewDocumentWorkflowController(service services.DocumentWorkflowServiceInterface) *DocumentWorkflowController {
	return &DocumentWorkflowController{
		service: service,
	}
}

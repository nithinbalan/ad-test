package controllers

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditUniverseController struct {
	service services.AuditUniverseServiceInterface
	db      *gorm.DB
}

func NewAuditUniverseController(service services.AuditUniverseServiceInterface) *AuditUniverseController {
	return &AuditUniverseController{
		service: service,
	}
}

func (auc *AuditUniverseController) MarkAsAuditableItem(c *gin.Context) {
	var auditUniverseFormDTO dtos.AuditUniverseFormDTO

	if err := c.ShouldBindJSON(&auditUniverseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(auc.db, c, err.Error())
		return
	}

	err = auc.service.Create(currentUserID, &auditUniverseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(auc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

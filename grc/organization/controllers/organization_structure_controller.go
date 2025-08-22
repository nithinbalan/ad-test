package controllers

import (
	"net/http"

	"jethur/grc/organization/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrganizationStructureController struct {
	service services.OrganizationStructureServiceInterface
	db      *gorm.DB
}

func NewOrganizationStructureController(service services.OrganizationStructureServiceInterface) *OrganizationStructureController {
	return &OrganizationStructureController{
		service: service,
	}
}

func (rhc *OrganizationStructureController) GetOrganizationStructure(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")

	organizationStructure, err := rhc.service.GetOrganizationStructure(sort, sortBy)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rhc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, organizationStructure)
}

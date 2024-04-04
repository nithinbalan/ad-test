package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceStatusUpdateController struct {
	service services.ComplianceStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewComplianceStatusUpdateController(service services.ComplianceStatusUpdateServiceInterface) *ComplianceStatusUpdateController {
	return &ComplianceStatusUpdateController{
		service: service,
	}
}

func (csuc *ComplianceStatusUpdateController) GetAllComplianceStatusUpdates(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	complianceStatuses, err := csuc.service.GetAll(uint(id), sort, sortBy, keyword)

	if err != nil {
		utility.SendInternalErrResp(csuc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, complianceStatuses)
}

func (csuc *ComplianceStatusUpdateController) CreateComplianceStatusUpdate(c *gin.Context) {

	var complianceStatusUpdateFormDTO complianceDTOs.ComplianceStatusUpdateFormDTO

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&complianceStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(csuc.db, c, err.Error())
		return
	}

	createdID, err := csuc.service.Create(uint(id), currentUserID, &complianceStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(csuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

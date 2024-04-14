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

type ComplianceActionPlanStatusUpdateController struct {
	service services.ComplianceActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewComplianceActionPlanStatusUpdateController(service services.ComplianceActionPlanStatusUpdateServiceInterface) *ComplianceActionPlanStatusUpdateController {
	return &ComplianceActionPlanStatusUpdateController{
		service: service,
	}
}

func (casuc *ComplianceActionPlanStatusUpdateController) GetAllComplianceActionPlanStatusUpdates(c *gin.Context) {

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	complianceStatuses, err := casuc.service.GetAll(uint(complianceID), uint(actionPlanID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(casuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, complianceStatuses)
}

func (casuc *ComplianceActionPlanStatusUpdateController) CreateComplianceActionPlanStatusUpdate(c *gin.Context) {

	var complianceStatusUpdateFormDTO complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
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
		utility.SendInternalErrResp(casuc.db, c, err.Error())
		return
	}

	createdID, err := casuc.service.Create(uint(complianceID), uint(actionPlanID), currentUserID, &complianceStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(casuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (casuc *ComplianceActionPlanStatusUpdateController) UpdateComplianceActionPlanStatusUpdate(c *gin.Context) {
	var complianceStatusUpdateFormDTO complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
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
		utility.SendInternalErrResp(casuc.db, c, err.Error())
		return
	}

	err = casuc.service.Update(uint(complianceID), uint(actionPlanID), uint(id), currentUserID, &complianceStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(casuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

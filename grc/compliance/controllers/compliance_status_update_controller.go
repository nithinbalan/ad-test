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
	complianceIdStr := c.Param("complianceId")

	complianceId, err := strconv.ParseUint(complianceIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	complianceStatuses, err := csuc.service.GetAll(uint(complianceId), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(csuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, complianceStatuses)
}

func (csuc *ComplianceStatusUpdateController) CreateComplianceStatusUpdate(c *gin.Context) {
	var complianceStatusUpdateFormDTO complianceDTOs.ComplianceStatusUpdateFormDTO

	complianceIdStr := c.Param("complianceId")

	complianceId, err := strconv.ParseUint(complianceIdStr, 10, 32)
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

	createdID, err := csuc.service.Create(uint(complianceId), currentUserID, &complianceStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
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

func (csuc *ComplianceStatusUpdateController) GetComplianceStatusUpdateDetails(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceStatusUpdate, err := csuc.service.GetFullDetailsByID(uint(complianceID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(csuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, complianceStatusUpdate)
}

func (csuc *ComplianceStatusUpdateController) UpdateComplianceStatusUpdate(c *gin.Context) {
	var complianceStatusUpdateFormDTO complianceDTOs.ComplianceStatusUpdateFormDTO

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
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
		utility.SendInternalErrResp(csuc.db, c, err.Error())
		return
	}

	err = csuc.service.Update(uint(complianceID), uint(id), currentUserID, &complianceStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(csuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (csuc *ComplianceStatusUpdateController) GetLatestComplianceStatusUpdate(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceStatusUpdate, err := csuc.service.GetLatestComplianceStatusUpdate(uint(complianceID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(csuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, complianceStatusUpdate)
}

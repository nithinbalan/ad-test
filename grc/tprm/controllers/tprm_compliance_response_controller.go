package controllers

import (
	"fmt"
	"jethur/grc/tprm/dtos"
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TPRMComplianceResponseController struct {
	service services.TPRMComplianceResponseServiceInterface
	db      *gorm.DB
}

func NewTPRMComplianceResponseController(service services.TPRMComplianceResponseServiceInterface) *TPRMComplianceResponseController {
	return &TPRMComplianceResponseController{
		service: service,
	}
}

func (tcrc *TPRMComplianceResponseController) GetAllTPRMComplianceResponses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	fmt.Println("1")
	impactCategories, err := tcrc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tcrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, impactCategories)
}

func (tcrc *TPRMComplianceResponseController) GetTPRMComplianceResponseDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceResponse, err := tcrc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, complianceResponse)
}

func (tcrc *TPRMComplianceResponseController) CreateTPRMComplianceResponse(c *gin.Context) {
	var complianceResponseFormDTO *dtos.TPRMComplianceResponseFormDTO

	if err := c.ShouldBindJSON(&complianceResponseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tcrc.db, c, err.Error())
		return
	}
	createdTPRMComplianceResponseID, err := tcrc.service.Create(currentUserID, complianceResponseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tcrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdTPRMComplianceResponseID)
}

func (tcrc *TPRMComplianceResponseController) UpdateTPRMComplianceResponse(c *gin.Context) {
	complianceResponseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var complianceResponseFormDTO *dtos.TPRMComplianceResponseFormDTO
	if err := c.ShouldBindJSON(&complianceResponseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tcrc.db, c, err.Error())
		return
	}

	err = tcrc.service.Update(uint(complianceResponseID), currentUserID, complianceResponseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tcrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tcrc *TPRMComplianceResponseController) DeleteTPRMComplianceResponses(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tcrc.db, c, err.Error())
		return
	}

	failedIDs, err := tcrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tcrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (tcrc *TPRMComplianceResponseController) ActivateTPRMComplianceResponses(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tcrc.db, c, err.Error())
		return
	}

	err = tcrc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tcrc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (tcrc *TPRMComplianceResponseController) DeactivateTPRMComplianceResponses(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tcrc.db, c, err.Error())
		return
	}

	err = tcrc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tcrc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

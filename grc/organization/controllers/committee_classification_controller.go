package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommitteeClassificationController struct {
	service services.CommitteeClassificationServiceInterface
	db      *gorm.DB
}

func NewCommitteeClassificationController(service services.CommitteeClassificationServiceInterface) *CommitteeClassificationController {
	return &CommitteeClassificationController{
		service: service,
	}
}

func (ccc *CommitteeClassificationController) GetAllCommitteeClassifications(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	queryParams := dtos.CommitteeClassificationListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
	}

	committeeClassifications, err := ccc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, committeeClassifications)
}

func (ccc *CommitteeClassificationController) GetCommitteeClassificationDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	committeeClassification, err := ccc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, committeeClassification)
}

func (ccc *CommitteeClassificationController) CreateCommitteeClassification(c *gin.Context) {
	var ccDTO dtos.CommitteeClassificationFormDTO

	// Bind request body
	if err := c.ShouldBindJSON(&ccDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get current user ID
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	// Call service
	id, err := ccc.service.Create(currentUserID, &ccDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(ccc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (ccc *CommitteeClassificationController) UpdateCommitteeClassification(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var ccDTO dtos.CommitteeClassificationFormDTO

	// Bind request body
	if err := c.ShouldBindJSON(&ccDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get current user ID
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	// Call service
	err = ccc.service.Update(uint(id), currentUserID, &ccDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(ccc.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (ccc *CommitteeClassificationController) DeleteCommitteeClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	// Bind request body
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get current user ID
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	// Call service
	failedIds, err := ccc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	if len(failedIds) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Some committee classifications could not be deleted",
			"failed_ids": failedIds,
		})
		return
	}

	c.Status(http.StatusOK)
}

func (ccc *CommitteeClassificationController) ActivateCommitteeClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	// Bind request body
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get current user ID
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	// Call service
	err = ccc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ccc *CommitteeClassificationController) DeactivateCommitteeClassifications(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	// Bind request body
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get current user ID
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	// Call service
	err = ccc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

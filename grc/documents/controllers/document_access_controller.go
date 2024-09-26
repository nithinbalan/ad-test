package controllers

import (
	"fmt"
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentAccessController struct {
	service services.DocumentAccessServiceInterface
	db      *gorm.DB
}

func NewDocumentAccessController(service services.DocumentAccessServiceInterface) *DocumentAccessController {
	return &DocumentAccessController{
		service: service,
	}
}

func (dac *DocumentAccessController) GetSharedUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(dac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(dac.db, c, "Invalid page_size parameter")
		return
	}

	sharedUsers, err := dac.service.GetSharedUsers(uint(id), uint(page), uint(pageSize))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, sharedUsers)
}

func (dac *DocumentAccessController) GetReviewUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reviewUsers, err := dac.service.GetReviewUsers(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, reviewUsers)
}

func (dac *DocumentAccessController) GetSharedUnits(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sharedUnits, err := dac.service.GetSharedUnits(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, sharedUnits)
}

func (dac *DocumentAccessController) CreateSharedUsers(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	var requestData struct {
		UserIDs []uint `json:"user_ids"`
	}

	var userIDs []uint
	if err := c.ShouldBindJSON(&requestData); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	userIDs = requestData.UserIDs

	if len(userIDs) == 0 {
		utility.SendInvalidErrResp(c, fmt.Errorf("the user IDs array cannot be empty"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dac.db, c, err.Error())
		return
	}

	err = dac.service.CreateSharedUsers(uint(documentID), userIDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(dac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Shared users created successfully"})
}

func (dac *DocumentAccessController) CreateSharedUnits(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	var requestData struct {
		UnitIDs []uint `json:"unit_ids"`
	}

	var unitIDs []uint
	if err := c.ShouldBindJSON(&requestData); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	unitIDs = requestData.UnitIDs

	if len(unitIDs) == 0 {
		utility.SendInvalidErrResp(c, fmt.Errorf("the unit IDs array cannot be empty"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dac.db, c, err.Error())
		return
	}

	err = dac.service.CreateSharedUnits(uint(documentID), unitIDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(dac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Shared unit created successfully"})
}

func (dac *DocumentAccessController) DeleteSharedUsers(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	var requestData struct {
		UserIDs []uint `json:"user_ids"`
	}

	var userIDs []uint
	if err := c.ShouldBindJSON(&requestData); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	userIDs = requestData.UserIDs

	if err := dac.service.DeleteSharedUsers(uint(documentID), userIDs); err != nil {
		utility.SendInternalErrResp(dac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shared users deleted successfully"})
}

func (dac *DocumentAccessController) DeleteSharedUnits(c *gin.Context) {
	documentIDStr := c.Param("id")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	var requestData struct {
		UnitIDs []uint `json:"unit_ids"`
	}

	var unitIDs []uint
	if err := c.ShouldBindJSON(&requestData); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	unitIDs = requestData.UnitIDs

	if err := dac.service.DeleteSharedUnits(uint(documentID), unitIDs); err != nil {
		utility.SendInternalErrResp(dac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shared unit deleted successfully"})
}

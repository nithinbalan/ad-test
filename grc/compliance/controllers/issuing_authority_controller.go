package controllers

import (
	"jethur/grc/compliance/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IssuingAuthorityController struct {
	service services.IssuingAuthorityServiceInterface
	db      *gorm.DB
}

func NewIssuingAuthorityController(service services.IssuingAuthorityServiceInterface) *IssuingAuthorityController {
	return &IssuingAuthorityController{
		service: service,
	}
}

func (iac *IssuingAuthorityController) GetAllIssuingAuthorities(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	authorities, err := iac.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, authorities)
}

func (iac *IssuingAuthorityController) GetIssuingAuthorityDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	authority, err := iac.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, authority)
}

func (iac *IssuingAuthorityController) CreateIssuingAuthority(c *gin.Context) {
	var authority *models.IssuingAuthority

	if err := c.ShouldBindJSON(&authority); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	createdAuthority, err := iac.service.Create(currentUserID, authority)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdAuthority)
}

func (iac *IssuingAuthorityController) UpdateIssuingAuthority(c *gin.Context) {
	var updatedAuthority *models.IssuingAuthority

	authorityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedAuthority); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	updatedAuthority, err = iac.service.Update(uint(authorityID), currentUserID, updatedAuthority)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedAuthority)
}

func (iac *IssuingAuthorityController) DeleteIssuingAuthorities(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	failedIDs, err := iac.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (iac *IssuingAuthorityController) ActivateIssuingAuthorities(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	if err := iac.service.Activate(request.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (iac *IssuingAuthorityController) DeactivateIssuingAuthorities(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	if err := iac.service.Deactivate(request.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(iac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

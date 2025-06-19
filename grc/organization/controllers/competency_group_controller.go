package controllers

import (
	"jethur/grc/models"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompetencyGroupController struct {
	service services.CompetencyGroupServiceInterface
	db      *gorm.DB
}

func NewCompetencyGroupController(service services.CompetencyGroupServiceInterface) *CompetencyGroupController {
	return &CompetencyGroupController{
		service: service,
	}
}

func (cgc *CompetencyGroupController) GetAllCompetencyGroups(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	groups, err := cgc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(cgc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, groups)
}

func (cgc *CompetencyGroupController) GetCompetencyGroupDetails(c *gin.Context) {
	idStr := c.Param("competencyGroupId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	group, err := cgc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, group)
}

func (cgc *CompetencyGroupController) CreateCompetencyGroup(c *gin.Context) {
	var group *models.CompetencyGroup

	if err := c.ShouldBindJSON(&group); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cgc.db, c, err.Error())
		return
	}

	created, err := cgc.service.Create(currentUserID, group)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, created)
}

func (cgc *CompetencyGroupController) UpdateCompetencyGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("competencyGroupId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updated *models.CompetencyGroup
	if err := c.ShouldBindJSON(&updated); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cgc.db, c, err.Error())
		return
	}

	updated, err = cgc.service.Update(uint(id), currentUserID, updated)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updated)
}

func (cgc *CompetencyGroupController) DeleteCompetencyGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cgc.db, c, err.Error())
		return
	}

	failedIDs, err := cgc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cgc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (cgc *CompetencyGroupController) ActivateCompetencyGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cgc.db, c, err.Error())
		return
	}

	err = cgc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cgc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (cgc *CompetencyGroupController) DeactivateCompetencyGroups(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cgc.db, c, err.Error())
		return
	}

	err = cgc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cgc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

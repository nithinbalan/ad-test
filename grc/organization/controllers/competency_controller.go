package controllers

import (
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompetencyController struct {
	service services.CompetencyServiceInterface
	db      *gorm.DB
}

func NewCompetencyController(service services.CompetencyServiceInterface) *CompetencyController {
	return &CompetencyController{
		service: service,
	}
}

func (cc *CompetencyController) GetAllCompetencies(c *gin.Context) {
	competencyGroupId, err := strconv.Atoi(c.Param("competencyGroupId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	competencies, err := cc.service.GetAll(uint(competencyGroupId), sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, competencies)
}

func (cc *CompetencyController) GetCompetencyDetails(c *gin.Context) {
	competencyGroupId, err := strconv.Atoi(c.Param("competencyGroupId"))
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

	competency, err := cc.service.GetByID(uint(competencyGroupId), uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, competency)
}

func (cc *CompetencyController) CreateCompetency(c *gin.Context) {
	competencyGroupId, err := strconv.Atoi(c.Param("competencyGroupId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var competency *organizationDTOs.CompetencyFormDTO

	if err := c.ShouldBindJSON(&competency); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	createdCompetency, err := cc.service.Create(uint(competencyGroupId), currentUserID, competency)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCompetency)
}

func (cc *CompetencyController) UpdateCompetency(c *gin.Context) {
	competencyGroupId, err := strconv.Atoi(c.Param("competencyGroupId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	competencyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedCompetency *models.Competency
	if err := c.ShouldBindJSON(&updatedCompetency); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	updatedCompetency, err = cc.service.Update(uint(competencyGroupId), uint(competencyID), currentUserID, updatedCompetency)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCompetency)
}

func (cc *CompetencyController) DeleteCompetencies(c *gin.Context) {
	competencyGroupId, err := strconv.Atoi(c.Param("competencyGroupId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	failedIDs, err := cc.service.Delete(uint(competencyGroupId), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (cc *CompetencyController) ActivateCompetencies(c *gin.Context) {
	competencyGroupId, err := strconv.Atoi(c.Param("competencyGroupId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	err = cc.service.Activate(uint(competencyGroupId), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (cc *CompetencyController) DeactivateCompetencies(c *gin.Context) {
	competencyGroupId, err := strconv.Atoi(c.Param("competencyGroupId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	err = cc.service.Deactivate(uint(competencyGroupId), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

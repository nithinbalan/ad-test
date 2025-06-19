package controllers

import (
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
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
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	allowedParams := []string{
		"competency_group_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := dtos.CompetencyListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	result, err := cc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (cc *CompetencyController) GetCompetencyDetails(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	competency, err := cc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, competency)
}

func (cc *CompetencyController) CreateCompetency(c *gin.Context) {

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

	createdCompetency, err := cc.service.Create(currentUserID, competency)
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

	updatedCompetency, err = cc.service.Update(uint(competencyID), currentUserID, updatedCompetency)
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

	failedIDs, err := cc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (cc *CompetencyController) ActivateCompetencies(c *gin.Context) {

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

	err = cc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (cc *CompetencyController) DeactivateCompetencies(c *gin.Context) {

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

	err = cc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

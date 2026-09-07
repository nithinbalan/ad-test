package controllers

import (
	"errors"
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSecurityGateController struct {
	service services.ProjectSecurityGateServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateController(service services.ProjectSecurityGateServiceInterface) *ProjectSecurityGateController {
	return &ProjectSecurityGateController{
		service: service,
	}
}

func (psgc *ProjectSecurityGateController) GetAllProjectSecurityGates(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "project_title_en")
	keyword := c.Query("keyword")

	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	validSortFields := map[string]string{
		"project_title_en":             "project_security_gates.project_title_en",
		"project_title_ar":             "project_security_gates.project_title_ar",
		"stage_title_en":               "stages.title_en",
		"stage_title_ar":               "stages.title_ar",
		"business_owner_first_name_en": "business_owner.first_name_en",
		"business_owner_first_name_ar": "business_owner.first_name_ar",
		"readiness_status_title_en":    "readiness_status.title_en",
		"readiness_status_title_ar":    "readiness_status.title_ar",
		"created_at":                   "project_security_gates.created_at",
		"updated_at":                   "project_security_gates.updated_at",
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := dtos.ProjectSecurityGateListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	projectSecurityGates, err := psgc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, projectSecurityGates)
}

func (psgc *ProjectSecurityGateController) GetProjectSecurityGateSummary(c *gin.Context) {
	summary, err := psgc.service.GetSummary()
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (psgc *ProjectSecurityGateController) GetProjectSecurityGateDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectSecurityGate, err := psgc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, projectSecurityGate)
}

func (psgc *ProjectSecurityGateController) CreateProjectSecurityGate(c *gin.Context) {
	var form dtos.ProjectSecurityGateFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	id, err := psgc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (psgc *ProjectSecurityGateController) UpdateProjectSecurityGate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ProjectSecurityGateFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	err = psgc.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psgc *ProjectSecurityGateController) DeleteProjectSecurityGates(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	failedIDs, err := psgc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

package controllers

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PolicyExceptionController struct {
	service services.PolicyExceptionServiceInterface
	db      *gorm.DB
}

func NewPolicyExceptionController(service services.PolicyExceptionServiceInterface) *PolicyExceptionController {
	return &PolicyExceptionController{
		service: service,
	}
}

func (pec *PolicyExceptionController) GetAllPolicyExceptions(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "document_title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	allowedParams := []string{
		"document_ids",
		"document_version_ids",
		"document_version_clause_ids",
		"document_requirement_group_ids",
		"requested_by_ids",
		"organization_unit_ids",
		"workflow_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(pec.db, c, "Invalid page parameter")
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(pec.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"reference_number":   "policy_exceptions.reference_number",
		"document_title_en":  "documents.title_en",
		"documents.title_en": "documents.title_en",
		"document_title_ar":  "documents.title_ar",
		"documents.title_ar": "documents.title_ar",
		"start_date":         "policy_exceptions.start_date",
		"end_date":           "policy_exceptions.end_date",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := dtos.PolicyExceptionListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	list, err := pec.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(pec.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}

func (pec *PolicyExceptionController) GetPolicyExceptionDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	data, err := pec.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, data)
}

func (pec *PolicyExceptionController) CreatePolicyException(c *gin.Context) {
	var dto dtos.PolicyExceptionFormDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pec.db, c, err.Error())
		return
	}

	id, err := pec.service.Create(currentUserID, &dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (pec *PolicyExceptionController) UpdatePolicyException(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var dto dtos.PolicyExceptionFormDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	userID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pec.db, c, err.Error())
		return
	}
	err = pec.service.Update(uint(id), &dto, userID)
	if err != nil {
		utility.SendInternalErrResp(pec.db, c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (pec *PolicyExceptionController) DeletePolicyExceptions(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pec.db, c, err.Error())
		return
	}

	failedIDs, err := pec.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pec.db, c, err.Error())
		return
	}

	response := commonDTOs.FailedIDsResponse{
		FailedIds: failedIDs,
	}

	c.JSON(http.StatusOK, response)
}
func (pec *PolicyExceptionController) MapRisksToPolicyException(c *gin.Context) {
	idStr := c.Param("id")
	policyExceptionID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var dto dtos.MapPolicyExceptionRisksDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pec.db, c, err.Error())
		return
	}

	if err := pec.service.MapRisks(uint(policyExceptionID), currentUserID, &dto); err != nil {
		utility.SendInternalErrResp(pec.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (pec *PolicyExceptionController) DeleteRisksFromPolicyException(c *gin.Context) {
	policyExceptionIDStr := c.Param("id")
	riskIDStr := c.Param("riskID")

	policyExceptionID, err := strconv.ParseUint(policyExceptionIDStr, 10, 64)
	if err != nil {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid policy exception ID"))
		return
	}

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid risk ID"))
		return
	}

	if err := pec.service.DeletePolicyExceptionRisk(uint(policyExceptionID), uint(riskID)); err != nil {
		utility.SendInternalErrResp(pec.db, c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

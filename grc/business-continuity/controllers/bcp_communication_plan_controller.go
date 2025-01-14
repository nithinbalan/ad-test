package controllers

import (
	"errors"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BcpCommunicationPlanController struct {
	service services.BcpCommunicationPlanServiceInterface
	db      *gorm.DB
}

func NewBcpCommunicationPlanController(service services.BcpCommunicationPlanServiceInterface) *BcpCommunicationPlanController {
	return &BcpCommunicationPlanController{
		service: service,
	}
}

func (bcpc *BcpCommunicationPlanController) GetAllBcpCommunicationPlans(c *gin.Context) {

	idStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"responsible_user_ids",
		"bcp_phase_ids",
	}
	validSortFields := map[string]string{
		"title":                          "bcp_communication_plans.title",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"responsible_user_last_name_en":  "responsible_user.last_name_en",
		"responsible_user_last_name_ar":  "responsible_user.last_name_ar",
		"responsible_user_email":         "responsible_user.email",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(bcpc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(bcpc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := businessContinuityDTOs.BcpCommunicationPlanListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	bcpCommunicationPlans, err := bcpc.service.GetAll(uint(bcpId), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(bcpc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, bcpCommunicationPlans)
}

func (bcpc *BcpCommunicationPlanController) GetBcpCommunicationPlanDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	bcpCommunicationPlan, err := bcpc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(bcpc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, bcpCommunicationPlan)
}

func (bcpc *BcpCommunicationPlanController) CreateBcpCommunicationPlan(c *gin.Context) {

	idStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var bcpCommunicationPlanFormDTO businessContinuityDTOs.BcpCommunicationPlanFormDTO

	if err := c.ShouldBindJSON(&bcpCommunicationPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bcpc.db, c, err.Error())
		return
	}

	id, err := bcpc.service.Create(uint(bcpId), currentUserID, &bcpCommunicationPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(bcpc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (bcpc *BcpCommunicationPlanController) UpdateBcpCommunicationPlan(c *gin.Context) {
	var bcpCommunicationPlanFormDTO businessContinuityDTOs.BcpCommunicationPlanFormDTO

	bcpCommunicationPlanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&bcpCommunicationPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bcpc.db, c, err.Error())
		return
	}

	err = bcpc.service.Update(uint(bcpCommunicationPlanID), currentUserID, &bcpCommunicationPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(bcpc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (bcpc *BcpCommunicationPlanController) DeleteBcpCommunicationPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bcpc.db, c, err.Error())
		return
	}

	failedIDs, err := bcpc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(bcpc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

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

type BcpRecoveryPlanController struct {
	service services.BcpRecoveryPlanServiceInterface
	db      *gorm.DB
}

func NewBcpRecoveryPlanController(service services.BcpRecoveryPlanServiceInterface) *BcpRecoveryPlanController {
	return &BcpRecoveryPlanController{
		service: service,
	}
}

func (brpc *BcpRecoveryPlanController) GetAllBcpRecoveryPlans(c *gin.Context) {

	idStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"severity_level_ids",
		"bcp_phase_ids",
		"responsible_user_ids",
	}
	validSortFields := map[string]string{
		"title":                          "bcp_recovery_plans.title",
		"severity_level_en":              "severity_levels.title_en",
		"severity_level_ar":              "severity_levels.title_ar",
		"severity_level_color_code":      "severity_levels.color_code",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"responsible_user_last_name_en":  "responsible_user.last_name_en",
		"responsible_user_last_name_ar":  "responsible_user.last_name_ar",
		"responsible_user_email":         "responsible_user.email",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := businessContinuityDTOs.BcpRecoveryPlanListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}
	bcpRecoveryPlans, err := brpc.service.GetAll(uint(bcpId), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, bcpRecoveryPlans)
}

func (brpc *BcpRecoveryPlanController) GetBcpRecoveryPlanDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	bcpRecoveryPlan, err := brpc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(brpc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, bcpRecoveryPlan)
}

func (brpc *BcpRecoveryPlanController) CreateBcpRecoveryPlan(c *gin.Context) {

	idStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var bcpRecoveryPlanFormDTO businessContinuityDTOs.BcpRecoveryPlanFormDTO

	if err := c.ShouldBindJSON(&bcpRecoveryPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, err.Error())
		return
	}

	id, err := brpc.service.Create(uint(bcpId), currentUserID, &bcpRecoveryPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(brpc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (brpc *BcpRecoveryPlanController) UpdateBcpRecoveryPlan(c *gin.Context) {
	var bcpRecoveryPlanFormDTO businessContinuityDTOs.BcpRecoveryPlanFormDTO

	bcpRecoveryPlanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&bcpRecoveryPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, err.Error())
		return
	}

	err = brpc.service.Update(uint(bcpRecoveryPlanID), currentUserID, &bcpRecoveryPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(brpc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (brpc *BcpRecoveryPlanController) DeleteBcpRecoveryPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, err.Error())
		return
	}

	failedIDs, err := brpc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(brpc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

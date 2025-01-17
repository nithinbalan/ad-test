package controllers

import (
	"errors"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BcpBiaMatrixController struct {
	service services.BcpBiaMatrixServiceInterface
	db      *gorm.DB
}

func NewBcpBiaMatrixController(service services.BcpBiaMatrixServiceInterface) *BcpBiaMatrixController {
	return &BcpBiaMatrixController{
		service: service,
	}
}

func (bcpc *BcpBiaMatrixController) GetAllBcpKeyProcesses(c *gin.Context) {

	idStr := c.Param("bcpId")

	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "process_title_en")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"ids",
	}
	validSortFields := map[string]string{
		"process_title_en": "processes.title_en",
		"process_title_ar": "processes.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := businessContinuityDTOs.BcpKeyProcessListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortField,
		Keyword: keyword,
		Filters: filters,
	}
	bcpKeyProcesses, err := bcpc.service.GetAllKeyProcess(uint(bcpId), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(bcpc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, bcpKeyProcesses)
}

func (bbmc *BcpBiaMatrixController) GetBcpBiaMatrixDetails(c *gin.Context) {
	idStr := c.Param("bcpId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	bcpBiaMatrix, err := bbmc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(bbmc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, bcpBiaMatrix)
}

func (bbmc *BcpBiaMatrixController) UpdateBcpBiaMatrix(c *gin.Context) {
	var bcpProcessFormDTO businessContinuityDTOs.BcpProcessFormDTO

	businessContinuityPlanID, err := strconv.Atoi(c.Param("bcpId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&bcpProcessFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bbmc.db, c, err.Error())
		return
	}

	err = bbmc.service.Update(uint(businessContinuityPlanID), currentUserID, &bcpProcessFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(bbmc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

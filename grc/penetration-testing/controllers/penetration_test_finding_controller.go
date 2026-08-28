package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestFindingController struct {
	service services.PenetrationTestFindingServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestFindingController(service services.PenetrationTestFindingServiceInterface) *PenetrationTestFindingController {
	return &PenetrationTestFindingController{
		service: service,
	}
}

func (ptfc *PenetrationTestFindingController) GetAllPenetrationTestFindings(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"reference_number":        "reference_number",
		"title":                   "title",
		"severity_level_id":       "severity_level_id",
		"target_remediation_date": "target_remediation_date",
		"created_at":              "created_at",
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page parameter"))
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page_size parameter"))
		return
	}
	if page < 1 || pageSize < 1 {
		utility.SendInvalidErrResp(c, errors.New("page and page_size must be greater than zero"))
		return
	}
	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	penetrationTestID, err := utility.ParseQueryParamToInt(c, "penetration_test_id", 0)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid penetration_test_id parameter"))
		return
	}
	if penetrationTestID < 0 {
		utility.SendInvalidErrResp(c, errors.New("penetration_test_id must be greater than zero"))
		return
	}

	queryParams := penetrationTestDTOs.PenetrationTestFindingListQueryParamsDTO{
		Sort:              sort,
		SortBy:            sortField,
		Keyword:           keyword,
		Page:              page,
		PageSize:          pageSize,
		PenetrationTestID: uint(penetrationTestID),
	}
	findings, err := ptfc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(ptfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, findings)
}

func (ptfc *PenetrationTestFindingController) GetPenetrationTestFindingDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	finding, err := ptfc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptfc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, finding)
}

func (ptfc *PenetrationTestFindingController) CreatePenetrationTestFinding(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestFindingFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptfc.db, c, err.Error())
		return
	}

	id, err := ptfc.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptfc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (ptfc *PenetrationTestFindingController) UpdatePenetrationTestFinding(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestFindingFormDTO

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptfc.db, c, err.Error())
		return
	}

	if err := ptfc.service.Update(uint(id), currentUserID, &formDTO); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptfc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ptfc *PenetrationTestFindingController) DeletePenetrationTestFindings(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptfc.db, c, err.Error())
		return
	}

	failedIDs, err := ptfc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ptfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

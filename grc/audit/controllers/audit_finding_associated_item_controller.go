package controllers

import (
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditFindingAssociatedItemController struct {
	service services.AuditFindingAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewAuditFindingAssociatedItemController(service services.AuditFindingAssociatedItemServiceInterface) *AuditFindingAssociatedItemController {
	return &AuditFindingAssociatedItemController{
		service: service,
	}
}

func (afaic *AuditFindingAssociatedItemController) GetAllAuditFindingDocuments(c *gin.Context) {

	findingIdStr := c.Param("findingId")

	findingId, err := strconv.ParseUint(findingIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	valueStr := c.Query("value")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(afaic.db, c, "Invalid value parameter")
			return
		}
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Value:    value,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := afaic.service.GetAllAuditFindingDocuments(uint(findingId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}
func (afaic *AuditFindingAssociatedItemController) CreateAuditFindingDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	findingId, err := strconv.ParseUint(c.Param("findingId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	err = afaic.service.CreateAuditFindingDocument(uint(findingId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (afaic *AuditFindingAssociatedItemController) DeleteAuditFindingDocument(c *gin.Context) {
	findingId, _ := strconv.ParseUint(c.Param("findingId"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := afaic.service.DeleteAuditFindingDocument(uint(findingId), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from AuditFinding successfully"})
}

func (afaic *AuditFindingAssociatedItemController) GetAllAuditFindingControls(c *gin.Context) {

	findingIdStr := c.Param("findingId")

	findingId, err := strconv.ParseUint(findingIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	valueStr := c.Query("value")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(afaic.db, c, "Invalid value parameter")
			return
		}
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Value:    value,
		Page:     page,
		PageSize: pageSize,
	}
	controls, err := afaic.service.GetAllAuditFindingControls(uint(findingId), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controls)
}
func (afaic *AuditFindingAssociatedItemController) CreateAuditFindingControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	findingId, err := strconv.ParseUint(c.Param("findingId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	err = afaic.service.CreateAuditFindingControl(uint(findingId), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (afaic *AuditFindingAssociatedItemController) DeleteAuditFindingControl(c *gin.Context) {
	findingId, _ := strconv.ParseUint(c.Param("findingId"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	err := afaic.service.DeleteAuditFindingControl(uint(findingId), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(afaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from AuditFinding successfully"})
}

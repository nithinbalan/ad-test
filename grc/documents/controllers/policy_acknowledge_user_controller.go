package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"

	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PolicyAcknowledgeUserController struct {
	service services.DocumentAcknowledgeUserServiceInterface
	db      *gorm.DB
}

func NewPolicyAcknowledgeUserController(service services.DocumentAcknowledgeUserServiceInterface) *PolicyAcknowledgeUserController {
	return &PolicyAcknowledgeUserController{
		service: service,
	}
}

func (pauc *PolicyAcknowledgeUserController) GetAllPolicyAcknowledgeUsers(c *gin.Context) {
	idStr := c.Param("id")
	documentId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	versionIdStr := c.Param("versionId")
	documentVersionId, err := strconv.ParseUint(versionIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "first_name_en")
	keyword := c.Query("keyword")

	allowedSortFields := map[string]bool{
		"first_name_en":              true,
		"first_name_ar":              true,
		"last_name_en":               true,
		"last_name_ar":               true,
		"email":                      true,
		"organization_unit_title_en": true,
		"organization_unit_title_ar": true,
	}

	if !allowedSortFields[sortBy] {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(pauc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(pauc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := documentDTOs.DocumentAcknowledgeListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	users, err := pauc.service.GetAll(uint(documentId), uint(documentVersionId), &queryParams)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (pauc *PolicyAcknowledgeUserController) CreatePolicyAcknowledgeUser(c *gin.Context) {

	var documentAcknowledgeFormDTO documentDTOs.DocumentAcknowledgeFormDTO
	if err := c.ShouldBindJSON(&documentAcknowledgeFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	documentIdStr := c.Param("id")

	documentId, err := strconv.ParseUint(documentIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	versionIdStr := c.Param("versionId")
	documentVersionId, err := strconv.ParseUint(versionIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pauc.db, c, err.Error())
		return
	}

	id, err := pauc.service.Create(uint(documentId), uint(documentVersionId), currentUserID, documentAcknowledgeFormDTO.Comment)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(pauc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (pauc *PolicyAcknowledgeUserController) DeletePolicyAcknowledgeUser(c *gin.Context) {
	documentIdStr := c.Param("id")

	documentId, err := strconv.ParseUint(documentIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	versionIdStr := c.Param("versionId")
	documentVersionId, err := strconv.ParseUint(versionIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pauc.db, c, err.Error())
		return
	}

	failedID, err := pauc.service.Delete(uint(documentId), uint(documentVersionId), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(pauc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}

func (pauc *PolicyAcknowledgeUserController) GetPolicyAcknowledgmentSummary(c *gin.Context) {
	documentIdStr := c.Param("id")
	documentId, err := strconv.ParseUint(documentIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	versionIdStr := c.Param("versionId")
	documentVersionId, err := strconv.ParseUint(versionIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := pauc.service.GetDocumentAcknowledgmentSummary(uint(documentId), uint(documentVersionId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pauc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}

package controllers

import (
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentClauseController struct {
	service services.DocumentClauseServiceInterface
	db      *gorm.DB
}

func NewDocumentClauseController(service services.DocumentClauseServiceInterface) *DocumentClauseController {
	return &DocumentClauseController{
		service: service,
	}
}

func (dcc *DocumentClauseController) GetAllDocumentClauses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := documentDTOs.DocumentClauseListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}

	documentClauses, err := dcc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, documentClauses)
}

func (dcc *DocumentClauseController) GetDocumentClauseDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Implementation goes here
	c.JSON(http.StatusOK, gin.H{"message": "not implemented yet", "id": id})
}

func (dcc *DocumentClauseController) CreateDocumentClause(c *gin.Context) {
	var documentClauseFormDTO documentDTOs.DocumentClauseFormDTO
	if err := c.ShouldBindJSON(&documentClauseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	idStr := c.Param("id")
	documentId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	id, err := dcc.service.Create(uint(documentId), currentUserID, &documentClauseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, id)
}

func (dcc *DocumentClauseController) UpdateDocumentClause(c *gin.Context) {
	var documentClauseFormDTO documentDTOs.DocumentClauseFormDTO
	if err := c.ShouldBindJSON(&documentClauseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	idStr := c.Param("id")
	documentId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	clauseIdStr := c.Param("clauseId")
	clauseId, err := strconv.ParseUint(clauseIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	err = dcc.service.Update(uint(documentId), uint(clauseId), currentUserID, &documentClauseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (dcc *DocumentClauseController) DeleteDocumentClauses(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	failedIDs, err := dcc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

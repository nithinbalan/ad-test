package controllers

import (
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PolicyClauseController struct {
	service services.PolicyClauseServiceInterface
	db      *gorm.DB
}

func NewPolicyClauseController(service services.PolicyClauseServiceInterface) *PolicyClauseController {
	return &PolicyClauseController{
		service: service,
	}
}

func (dcc *PolicyClauseController) GetAllPolicyClauses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := policyDTOs.PolicyClauseListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}

	idStr := c.Param("id")
	documentId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	policyClauses, err := dcc.service.GetAll(&queryParams, uint(documentId))
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, policyClauses)
}

func (dcc *PolicyClauseController) GetPolicyClauseDetails(c *gin.Context) {
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

	clause, err := dcc.service.GetFullDetailsByID(uint(documentId), uint(clauseId))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, clause)
}

func (dcc *PolicyClauseController) CreatePolicyClause(c *gin.Context) {
	var policyClauseFormDTO policyDTOs.PolicyClauseFormDTO
	if err := c.ShouldBindJSON(&policyClauseFormDTO); err != nil {
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

	id, err := dcc.service.Create(uint(documentId), currentUserID, &policyClauseFormDTO)
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

func (dcc *PolicyClauseController) UpdatePolicyClause(c *gin.Context) {
	var policyClauseFormDTO policyDTOs.PolicyClauseFormDTO
	if err := c.ShouldBindJSON(&policyClauseFormDTO); err != nil {
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

	err = dcc.service.Update(uint(documentId), uint(clauseId), currentUserID, &policyClauseFormDTO)
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

func (dcc *PolicyClauseController) DeletePolicyClauses(c *gin.Context) {
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

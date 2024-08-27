package controllers

import (
	auditFindingDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditFindingRootCauseAnalysisController struct {
	service services.AuditFindingRootCauseAnalysisServiceInterface
	db      *gorm.DB
}

func NewAuditFindingRootCauseAnalysisController(service services.AuditFindingRootCauseAnalysisServiceInterface) *AuditFindingRootCauseAnalysisController {
	return &AuditFindingRootCauseAnalysisController{
		service: service,
	}
}

func (afrcac *AuditFindingRootCauseAnalysisController) GetAllAuditFindingRootCauseAnalyses(c *gin.Context) {
	auditIDStr := c.Param("auditId")
	findingIDStr := c.Param("id")

	auditID, err := strconv.ParseUint(auditIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := auditFindingDTOs.AuditFindingRootCauseAnalysisListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}
	auditFindingRootCauseAnalyses, err := afrcac.service.GetAll(uint(auditID), uint(findingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afrcac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditFindingRootCauseAnalyses)
}

func (afrcac *AuditFindingRootCauseAnalysisController) CreateAuditFindingRootCauseAnalysis(c *gin.Context) {
	var AuditFindingRootCauseAnalysisFormDTO auditFindingDTOs.AuditFindingRootCauseAnalysisFormDTO

	auditIDStr := c.Param("auditId")
	findingIDStr := c.Param("id")

	auditID, err := strconv.ParseUint(auditIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&AuditFindingRootCauseAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrcac.db, c, err.Error())
		return
	}

	id, err := afrcac.service.Create(uint(auditID), uint(findingID), currentUserID, &AuditFindingRootCauseAnalysisFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afrcac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (afrcac *AuditFindingRootCauseAnalysisController) UpdateAuditFindingRootCauseAnalysis(c *gin.Context) {
	var AuditFindingRootCauseAnalysisFormDTO auditFindingDTOs.AuditFindingRootCauseAnalysisFormDTO

	auditIDStr := c.Param("auditId")
	findingIDStr := c.Param("id")

	auditID, err := strconv.ParseUint(auditIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("analysisId")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&AuditFindingRootCauseAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrcac.db, c, err.Error())
		return
	}

	err = afrcac.service.Update(uint(auditID), uint(findingID), uint(id), currentUserID, &AuditFindingRootCauseAnalysisFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afrcac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (afrcac *AuditFindingRootCauseAnalysisController) DeleteAuditFindingRootCauseAnalyses(c *gin.Context) {
	auditIDStr := c.Param("auditId")
	findingIDStr := c.Param("id")

	auditID, err := strconv.ParseUint(auditIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("analysisId")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrcac.db, c, err.Error())
		return
	}

	failedID, err := afrcac.service.Delete(uint(auditID), uint(findingID), uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afrcac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}

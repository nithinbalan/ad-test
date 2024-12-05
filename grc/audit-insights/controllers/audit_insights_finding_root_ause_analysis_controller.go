package controllers

import (
	"jethur/grc/audit-insights/services"
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditInsightsFindingRootCauseAnalysisController struct {
	service services.AuditInsightsFindingRootCauseAnalysisServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsFindingRootCauseAnalysisController(service services.AuditInsightsFindingRootCauseAnalysisServiceInterface) *AuditInsightsFindingRootCauseAnalysisController {
	return &AuditInsightsFindingRootCauseAnalysisController{
		service: service,
	}
}

func (aifrcac *AuditInsightsFindingRootCauseAnalysisController) GetAllAuditInsightsFindingRootCauseAnalyses(c *gin.Context) {

	findingIDStr := c.Param("findingId")

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := auditDTOs.AuditFindingRootCauseAnalysisListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}
	auditFindingRootCauseAnalyses, err := aifrcac.service.GetAll(uint(findingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aifrcac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditFindingRootCauseAnalyses)
}

func (aifrcac *AuditInsightsFindingRootCauseAnalysisController) CreateAuditInsightsFindingRootCauseAnalysis(c *gin.Context) {
	var AuditInsightsFindingRootCauseAnalysisFormDTO auditDTOs.AuditFindingRootCauseAnalysisFormDTO

	findingIDStr := c.Param("findingId")

	findingID, err := strconv.ParseUint(findingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&AuditInsightsFindingRootCauseAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aifrcac.db, c, err.Error())
		return
	}

	id, err := aifrcac.service.Create(uint(findingID), currentUserID, &AuditInsightsFindingRootCauseAnalysisFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aifrcac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (aifrcac *AuditInsightsFindingRootCauseAnalysisController) UpdateAuditInsightsFindingRootCauseAnalysis(c *gin.Context) {
	var AuditInsightsFindingRootCauseAnalysisFormDTO auditDTOs.AuditFindingRootCauseAnalysisFormDTO

	findingIDStr := c.Param("findingId")

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

	if err := c.ShouldBindJSON(&AuditInsightsFindingRootCauseAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aifrcac.db, c, err.Error())
		return
	}

	err = aifrcac.service.Update(uint(findingID), uint(id), currentUserID, &AuditInsightsFindingRootCauseAnalysisFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aifrcac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (aifrcac *AuditInsightsFindingRootCauseAnalysisController) DeleteAuditInsightsFindingRootCauseAnalyses(c *gin.Context) {
	findingIDStr := c.Param("findingId")

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
		utility.SendInternalErrResp(aifrcac.db, c, err.Error())
		return
	}

	failedID, err := aifrcac.service.Delete(uint(findingID), uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aifrcac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}

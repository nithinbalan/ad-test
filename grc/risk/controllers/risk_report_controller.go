package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	"jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskReportController struct {
	service services.RiskReportServiceInterface
	db      *gorm.DB
}

func NewRiskReportController(service services.RiskReportServiceInterface) *RiskReportController {
	return &RiskReportController{
		service: service,
	}
}

func (rrc *RiskReportController) GetAllRiskReports(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"title":      "risk_reports.title",
		"start_date": "risk_reports.start_date",
		"end_date":   "risk_reports.end_date",
		"created_at": "risk_reports.created_at",
	}

	filters := utility.ParseQueryFilters([]string{}, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := dtos.RiskReportListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	reports, err := rrc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (rrc *RiskReportController) GetRiskReportDetails(c *gin.Context) {

	riskReportID, err := strconv.Atoi(c.Param("riskReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := rrc.service.GetRiskReportDetails(uint(riskReportID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, details)
}

func (rrc *RiskReportController) CreateRiskReport(c *gin.Context) {
	var request dtos.RiskReportFormDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	riskReportID, err := rrc.service.Create(currentUserID, &request)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":        "Risk report generated successfully",
		"risk_report_id": riskReportID,
	})
}

func (rrc *RiskReportController) GetRiskReportSection(c *gin.Context) {

	riskReportID, err := strconv.Atoi(c.Param("riskReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := rrc.service.GetSectionByID(uint(riskReportID), uint(sectionID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, details)
}

func (rrc *RiskReportController) UpdateRiskReportSection(c *gin.Context) {
	riskReportID, err := strconv.Atoi(c.Param("riskReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.RiskReportSectionFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.UpdateSection(uint(riskReportID), uint(sectionID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrc *RiskReportController) DeleteRiskReportSection(c *gin.Context) {
	riskReportID, err := strconv.Atoi(c.Param("riskReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.DeleteRiskReportSection(uint(riskReportID), uint(sectionID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrc *RiskReportController) DeleteRiskReport(c *gin.Context) {
	riskReportID, err := strconv.Atoi(c.Param("riskReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.DeleteRiskReport(uint(riskReportID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (rrc *RiskReportController) GetRiskReportDocumentPreview(c *gin.Context) {
	idStr := c.Param("riskReportId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := rrc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rrc.db, c, err.Error())
		}
		return
	}
	c.File(documentPath)
}

func (rrc *RiskReportController) DownloadRiskReportDocument(c *gin.Context) {
	idStr := c.Param("riskReportId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := rrc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rrc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(documentPath)

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(documentPath)
}

func (rc *RiskReportController) CreateRiskReportDocuments(c *gin.Context) {
	var riskReportDocumentFormDTO dtos.RiskReportDocumentFormDTO

	riskReportIDStr := c.Param("riskReportId")
	riskReportID, err := strconv.ParseUint(riskReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&riskReportDocumentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	createdID, err := rc.service.UploadRiskReportDocuments(uint(riskReportID), currentUserID, &riskReportDocumentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": createdID})
}

func (rrc *RiskReportController) DeleteRiskReportDocument(c *gin.Context) {
	riskReportIDStr := c.Param("riskReportId")
	riskReportID, err := strconv.ParseUint(riskReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentIDStr := c.Param("documentId")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
			return
		}
	}

	err = rrc.service.DeleteRiskReportDocument(uint(riskReportID), uint(documentID), currentUserID)
	if err != nil {

		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

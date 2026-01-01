package controllers

import (
	"errors"
	"fmt"
	"jethur/grc/tprm/dtos"
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThirdPartyAssessmentController struct {
	service services.ThirdPartyAssessmentServiceInterface
	db      *gorm.DB
}

func NewThirdPartyAssessmentController(service services.ThirdPartyAssessmentServiceInterface) *ThirdPartyAssessmentController {
	return &ThirdPartyAssessmentController{
		service: service,
	}
}

func (tprc *ThirdPartyAssessmentController) GetAllThirdPartyAssessment(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(tprc.db, c, "Invalid export parameter")
			return
		}
	}
	allowedParams := []string{
		"third_party_register_ids", "question_set_ids", "third_party_assessment_status_ids", "risk_level_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

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

	validSortFields := map[string]string{
		"third_party_register_title": "third_party_register.title",
		"title":                      "third_party_assessments.title",
		"reference_number":           "third_party_assessments.reference_number",
		"due_date":                   "third_party_assessments.due_date",
		"risk_level_title_en":        "risk_levels.title_en",
		"risk_level_title_ar":        "risk_levels.title_ar",
		"question_set_title_en":      "question_sets.title_en",
		"question_set_title_ar":      "question_sets.title_ar",
		"assessment_status_title_en": "third_party_assessment_statuses.title_en",
		"assessment_status_title_ar": "third_party_assessment_statuses.title_ar",
		"created_user_first_name_en": "created_user.first_name_en",
		"created_user_last_name_en":  "created_user.last_name_en",
		"created_at":                 "third_party_assessments.created_at",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid sort_by parameter"))
		return
	}

	query := dtos.ThirdPartyAssessmentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
	}

	thirdParties, err := tprc.service.GetAll(&query)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(tprc.db, c, err.Error())
			return
		}

		excelFile, err := tprc.service.ExportThirdPartyAssessments(thirdParties.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(tprc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Question Library.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, thirdParties)
}

func (tprc *ThirdPartyAssessmentController) GetThirdPartyAssessmentExport(c *gin.Context) {
	c.Set("export", true)

	tprc.GetAllThirdPartyAssessment(c)
}

func (tprc *ThirdPartyAssessmentController) GetThirdPartyAssessmentDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	result, err := tprc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (tprc *ThirdPartyAssessmentController) CreateThirdPartyAssessment(c *gin.Context) {
	var form dtos.ThirdPartyAssessmentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	id, err := tprc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (tprc *ThirdPartyAssessmentController) UpdateThirdPartyAssessment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ThirdPartyAssessmentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	err = tprc.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tprc *ThirdPartyAssessmentController) DeleteThirdPartyAssessment(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	failedIDs, err := tprc.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (tprc *ThirdPartyAssessmentController) ShareLink(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	err = tprc.service.ShareLink(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tprc *ThirdPartyAssessmentController) GetTPRMAssessmentExport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assessmentDetails, err := tprc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	excelFile, err := tprc.service.ExportTPRMAssessment(assessmentDetails, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}
	defer excelFile.Close()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", `attachment; filename="Assessment.xlsx"`)
	c.Header("Content-Type", "application/octet-stream")
	c.File(excelFile.Name())

}

func (tprc *ThirdPartyAssessmentController) ValidateTPRMAssessmentTemplate(c *gin.Context) {

	fileHeader, err := c.FormFile("file")
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, "Failed to read uploaded file")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, "Failed to open file")
		return
	}
	defer file.Close()

	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	validationResult, err := tprc.service.ValidateTPRMAssessmentTemplate(uint(id), file)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, validationResult)
}

func (tprc *ThirdPartyAssessmentController) ValidateImportTPRMAssessment(c *gin.Context) {

	var input []map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	validationResult, err := tprc.service.ValidateImportTPRMAssessment(uint(id), input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, validationResult)
}

func (tprc *ThirdPartyAssessmentController) ImportTPRMAssessment(c *gin.Context) {
	var input []map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	createdIDs, err := tprc.service.ImportTPRMAssessment(uint(id), input, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"created_ids": createdIDs})
}

func (tprc *ThirdPartyAssessmentController) MarkAsSend(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, _ := utility.GetCurrentUserID(c)

	err = tprc.service.MarkAsSend(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (tprc *ThirdPartyAssessmentController) GetThirdPartyAssessmentAnswerSummary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	response, err := tprc.service.GetAssessmentAnswerSummary(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}
	// c.Status(http.StatusOK)
	c.JSON(http.StatusOK, response)
}

func (tprc *ThirdPartyAssessmentController) GetTPRMAssessmentReportExport(c *gin.Context) {
	assessmentID, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	fileName, err := tprc.service.ExportTPRMAssessmentReport(uint(assessmentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return

		default:
			utility.SendInternalErrResp(tprc.db, c, fmt.Sprintf("Failed to export PDF: %s", err.Error()))
			return
		}
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/pdf")

	c.File(fileName)
	c.Writer.Flush()

	go func() {
		time.Sleep(2 * time.Second)
		err := os.Remove(fileName)
		if err != nil {
			fmt.Printf("Failed to delete file %s: %v\n", fileName, err)
		} else {
			fmt.Printf("Successfully deleted file %s\n", fileName)
		}
	}()
}

func (tprc *ThirdPartyAssessmentController) UpdateAnswer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("assessmentID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionId, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ThirdPartyAssessmentAnswerFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	err = tprc.service.UpdateAnswer(uint(id), uint(questionId), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tprc *ThirdPartyAssessmentController) GetAssessmentSectionDetails(c *gin.Context) {
	idStr := c.Param("assessmentID")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionIdStr := c.Param("sectionId")

	sectionId, err := strconv.ParseUint(sectionIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assessmentPlanSection, err := tprc.service.GetAssessmentSectionDetails(uint(id), uint(sectionId))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, assessmentPlanSection)
}

func (tprc *ThirdPartyAssessmentController) GetAssessmentSectionQuestionDetails(c *gin.Context) {
	idStr := c.Param("assessmentID")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionIdStr := c.Param("sectionId")

	sectionId, err := strconv.ParseUint(sectionIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionIdStr := c.Param("questionId")

	questionId, err := strconv.ParseUint(questionIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	question, err := tprc.service.GetAssessmentSectionQuestionDetails(uint(id), uint(sectionId), uint(questionId))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, question)
}

func (tprc *ThirdPartyAssessmentController) GetAllAnswerDocuments(c *gin.Context) {

	assessmentIDStr := c.Param("assessmentID")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionIDStr := c.Param("questionId")

	questionID, err := strconv.ParseUint(questionIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	answerIDStr := c.Param("answerId")

	answerID, err := strconv.ParseUint(answerIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documents, err := tprc.service.GetAllAnswerDocuments(uint(assessmentID), uint(questionID), uint(answerID))

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, documents)
}
func (tprc *ThirdPartyAssessmentController) CreateAnswerDocument(c *gin.Context) {
	var input dtos.AnswerDocumentFormDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	assessmentID, err := strconv.ParseUint(c.Param("assessmentID"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	questionIDStr := c.Param("questionId")

	questionID, err := strconv.ParseUint(questionIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	answerIDStr := c.Param("answerId")

	answerID, err := strconv.ParseUint(answerIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	err = tprc.service.CreateAnswerDocument(uint(assessmentID), uint(questionID), uint(answerID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (tprc *ThirdPartyAssessmentController) DeleteAnswerDocument(c *gin.Context) {
	assessmentID, _ := strconv.ParseUint(c.Param("assessmentID"), 10, 64)
	questionID, _ := strconv.ParseUint(c.Param("questionId"), 10, 64)
	answerID, _ := strconv.ParseUint(c.Param("answerId"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := tprc.service.DeleteAnswerDocument(uint(assessmentID), uint(questionID), uint(answerID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(tprc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document Deleted successfully"})
}

func (tprc *ThirdPartyAssessmentController) GetAnswerDocumentPreview(c *gin.Context) {
	assessmentID, _ := strconv.ParseUint(c.Param("assessmentID"), 10, 64)
	questionID, _ := strconv.ParseUint(c.Param("questionId"), 10, 64)
	answerID, _ := strconv.ParseUint(c.Param("answerId"), 10, 64)
	token := c.Query("token")

	fmt.Printf("Parameters received - AssessmentID: %s, QuestionID: %s, AnswerID: %s, Token: %s\n", assessmentID, questionID, answerID, token)

	documentPath, err := tprc.service.GetDocumentPath(uint(assessmentID), uint(questionID), uint(answerID), token)
	if err != nil {
		fmt.Printf("Error in service call for document path: %v\n", err)
		utility.SendNotFoundErrResp(c)
		return
	}

	c.File(documentPath)
}

func (tprc *ThirdPartyAssessmentController) DownloadAnswerDocument(c *gin.Context) {
	assessmentID, _ := strconv.ParseUint(c.Param("assessmentID"), 10, 64)
	questionID, _ := strconv.ParseUint(c.Param("questionId"), 10, 64)
	answerID, _ := strconv.ParseUint(c.Param("answerId"), 10, 64)
	token := c.Query("token")

	documentPath, err := tprc.service.GetDocumentPath(uint(assessmentID), uint(questionID), uint(answerID), token)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	fileName := filepath.Base(documentPath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}

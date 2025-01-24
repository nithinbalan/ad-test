package controllers

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestAndExerciseTestResultController struct {
	service services.TestAndExerciseTestResultServiceInterface
	db      *gorm.DB
}

func NewTestAndExerciseTestResultController(service services.TestAndExerciseTestResultServiceInterface) *TestAndExerciseTestResultController {
	return &TestAndExerciseTestResultController{
		service: service,
	}
}

func (tetrc *TestAndExerciseTestResultController) GetTestAndExerciseTestResultDetails(c *gin.Context) {
	testAndExerciseIdStr := c.Param("testAndExerciseId")

	testAndExerciseId, err := strconv.ParseUint(testAndExerciseIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	communicationActionPlan, err := tetrc.service.GetFullDetailsByID(uint(testAndExerciseId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tetrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, communicationActionPlan)
}

func (tetrc *TestAndExerciseTestResultController) CreateTestAndExerciseTestResult(c *gin.Context) {

	idStr := c.Param("testAndExerciseId")

	testAndExerciseId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var communicationActionPlanFormDTO businessContinuityDTOs.TestAndExerciseTestResultFormDTO

	if err := c.ShouldBindJSON(&communicationActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tetrc.db, c, err.Error())
		return
	}

	id, err := tetrc.service.Create(uint(testAndExerciseId), currentUserID, &communicationActionPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tetrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (tetrc *TestAndExerciseTestResultController) UpdateTestAndExerciseTestResult(c *gin.Context) {
	testAndExerciseIdStr := c.Param("testAndExerciseId")

	testAndExerciseId, err := strconv.ParseUint(testAndExerciseIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var communicationActionPlanFormDTO businessContinuityDTOs.TestAndExerciseTestResultFormDTO

	if err := c.ShouldBindJSON(&communicationActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tetrc.db, c, err.Error())
		return
	}

	err = tetrc.service.Update(uint(testAndExerciseId), currentUserID, &communicationActionPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tetrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tetrc *TestAndExerciseTestResultController) DeleteTestAndExerciseTestResults(c *gin.Context) {
	testAndExerciseIdStr := c.Param("testAndExerciseId")

	testAndExerciseId, err := strconv.ParseUint(testAndExerciseIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tetrc.db, c, err.Error())
		return
	}

	failedIDs, err := tetrc.service.Delete(uint(testAndExerciseId), currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tetrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (tetrc *TestAndExerciseTestResultController) GetTestAndExerciseTestResultDocumentPreview(c *gin.Context) {

	testAndExerciseIDStr := c.Param("testAndExerciseId")

	testAndExerciseID, err := strconv.ParseUint(testAndExerciseIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	testResultDocumentPath, err := tetrc.service.GetDocumentPath(uint(testAndExerciseID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(tetrc.db, c, err.Error())
		}
		return
	}
	c.File(testResultDocumentPath)
}

func (tetrc *TestAndExerciseTestResultController) DownloadTestAndExerciseTestResultDocument(c *gin.Context) {
	testAndExerciseIDStr := c.Param("testAndExerciseId")

	testAndExerciseID, err := strconv.ParseUint(testAndExerciseIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	testResultDocumentPath, err := tetrc.service.GetDocumentPath(uint(testAndExerciseID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(tetrc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(testResultDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(testResultDocumentPath)
}

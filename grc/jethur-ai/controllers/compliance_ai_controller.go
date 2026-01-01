package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/jethur-ai/services"
	"jethur/utility"
	"net/http"
	"strconv"
)

type ComplianceAIController struct {
	service services.ComplianceAIServiceInterface
	db      *gorm.DB
}

func NewComplianceAIController(service services.ComplianceAIServiceInterface) *ComplianceAIController {
	return &ComplianceAIController{
		service: service,
	}
}

func (caic *ComplianceAIController) GetComplianceAutofill(c *gin.Context) {
	var req dtos.ComplianceAutofillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	result, err := caic.service.GetComplianceAutofill(req)
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}


func (caic *ComplianceAIController) ImportComplianceRegister(c *gin.Context) {

    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
        return
    }

    result, err := caic.service.ImportComplianceRegister(file)
    if err != nil {
        utility.SendInternalErrResp(caic.db, c, err.Error())
        return
    }

    c.JSON(http.StatusOK, result)
}



func (caic *ComplianceAIController) GetComplianceSummary(c *gin.Context) {
	complianceID, err := strconv.ParseUint(c.Param("complianceId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid compliance ID"})
		return
	}

	result, err := caic.service.GetComplianceSummary(uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (caic *ComplianceAIController) GetComplianceClauseGuidance(c *gin.Context) {
	clauseID, err := strconv.ParseUint(c.Param("clauseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid clause ID"})
		return
	}

	result, err := caic.service.GetComplianceClauseGuidance(uint(clauseID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (caic *ComplianceAIController) GetComplianceStrategicObjectivesJob(c *gin.Context) {

	complianceID, err := strconv.ParseUint(c.Param("complianceId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid compliance ID"})
		return
	}

	result, err := caic.service.GetComplianceStrategicObjectivesJob(uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (caic *ComplianceAIController) GetComplianceSuggestKPIJob(c *gin.Context) {

	complianceID, err := strconv.ParseUint(c.Param("complianceId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid compliance ID"})
		return
	}

	result, err := caic.service.GetComplianceSuggestKPIJob(uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (caic *ComplianceAIController) GetComplianceSuggestProcessJob(c *gin.Context) {

	complianceID, err := strconv.ParseUint(c.Param("complianceId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid compliance ID"})
		return
	}

	result, err := caic.service.GetComplianceSuggestProcessJob(uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (caic *ComplianceAIController) GetComplianceSuggestRisksJob(c *gin.Context) {

	complianceID, err := strconv.ParseUint(
		c.Param("complianceId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid compliance ID"})
		return
	}

	result, err := caic.service.GetComplianceSuggestRisksJob(uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (caic *ComplianceAIController) GetComplianceSuggestControlsJob(c *gin.Context) {

	complianceID, err := strconv.ParseUint(
		c.Param("complianceId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid compliance ID"})
		return
	}

	result, err := caic.service.GetComplianceSuggestControlsJob(uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (caic *ComplianceAIController) GetComplianceSuggestAssociatesJob(c *gin.Context) {

	complianceID, err := strconv.ParseUint(c.Param("complianceId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid compliance ID"})
		return
	}

	result, err := caic.service.GetComplianceSuggestAssociatesJob(uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(caic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ajc *ComplianceAIController) GetComplianceJobStatus(c *gin.Context) {
	jobID := c.Param("jobId")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	result, err := ajc.service.GetComplianceJobStatus(jobID)
	if err != nil {
		utility.SendInternalErrResp(ajc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

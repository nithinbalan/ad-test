package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceClauseController struct {
	service services.ComplianceClauseServiceInterface
	db      *gorm.DB
}

func NewComplianceClauseController(service services.ComplianceClauseServiceInterface) *ComplianceClauseController {
	return &ComplianceClauseController{service: service}
}

func (ccc *ComplianceClauseController) GetAllComplianceClauses(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	clauses, err := ccc.service.GetAll(uint(complianceID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clauses)
}

func (ccc *ComplianceClauseController) GetComplianceClauseDetails(c *gin.Context) {
	idStr := c.Param("clauseId")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	clause, err := ccc.service.GetByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, clause)
}

func (ccc *ComplianceClauseController) CreateComplianceClause(c *gin.Context) {
	var complianceClauseFormDTO complianceDTOs.ComplianceClauseFormDTO

	if err := c.ShouldBindJSON(&complianceClauseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(ccc.db, c, e.Error())
			return
		}
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	id, err := ccc.service.Create(uint(complianceID), currentUserID, &complianceClauseFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ccc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (ccc *ComplianceClauseController) UpdateComplianceClause(c *gin.Context) {
	var complianceClauseFormDTO complianceDTOs.ComplianceClauseFormDTO

	if err := c.ShouldBindJSON(&complianceClauseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("clauseId")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	err = ccc.service.Update(uint(complianceID), uint(id), currentUserID, &complianceClauseFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ccc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ccc *ComplianceClauseController) DeleteComplianceClause(c *gin.Context) {
	idStr := c.Param("clauseId")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = ccc.service.Delete(uint(id))
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compliance Clause deleted successfully"})
}

func (ccc *ComplianceClauseController) GetAssociatedComplianceRequirement(c *gin.Context) {
	complianceID, err := strconv.Atoi(c.Param("complianceId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	clauseID, err := strconv.Atoi(c.Param("clauseId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Fetch associated compliance requirements
	requirements, err := ccc.service.GetAssociatedComplianceRequirements(uint(complianceID), uint(clauseID))
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, requirements)
}

func (cc *ComplianceClauseController) GetAssociatedControls(c *gin.Context) {
	complianceID, err := strconv.Atoi(c.Param("complianceId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	clauseID, err := strconv.Atoi(c.Param("clauseId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Get associated controls
	controls, err := cc.service.GetAssociatedControls(uint(complianceID), uint(clauseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, controls)
}

func (ccc *ComplianceClauseController) CreateComplianceClauseControls(c *gin.Context) {
	var input complianceDTOs.ComplianceClauseControlMapDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	complianceIDStr := c.Param("complianceId")
	clauseIDStr := c.Param("clauseId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid complianceId"))
		return
	}
	clauseID, err := strconv.ParseUint(clauseIDStr, 10, 64)
	if err != nil {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid clauseId"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	err = ccc.service.CreateComplianceClauseControls(uint(complianceID), uint(clauseID), currentUserID, &input)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (ccc *ComplianceClauseController) DeleteComplianceClauseControl(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	clauseIDStr := c.Param("clauseId")
	controlIDStr := c.Param("controlId")
	_, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid complianceId"))
		return
	}
	clauseID, err := strconv.ParseUint(clauseIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlID, err := strconv.ParseUint(controlIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = ccc.service.DeleteComplianceClauseControl(uint(clauseID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compliance Clause control mapping deleted successfully"})
}

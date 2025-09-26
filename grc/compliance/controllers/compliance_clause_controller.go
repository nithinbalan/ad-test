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
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	// Get associated controls
	controls, err := cc.service.GetAssociatedControls(uint(complianceID), uint(clauseID), &queryParams)
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

func (cc *ComplianceClauseController) GetAssociatedRisks(c *gin.Context) {
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
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := cc.service.GetAssociatedRisks(uint(complianceID), uint(clauseID), &queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, risks)
}

func (ccc *ComplianceClauseController) CreateComplianceClauseRisk(c *gin.Context) {
	var input commonDTOs.AssociatedRiskMapDTO

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

	err = ccc.service.CreateComplianceClauseRisk(uint(complianceID), uint(clauseID), currentUserID, &input)
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

	c.Status(http.StatusCreated)
}

func (ccc *ComplianceClauseController) DeleteComplianceClauseRisk(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	clauseIDStr := c.Param("clauseId")
	riskIDStr := c.Param("riskId")
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

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = ccc.service.DeleteComplianceClauseRisk(uint(clauseID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compliance Clause risk mapping deleted successfully"})
}

func (cc *ComplianceClauseController) GetAssociatedDocuments(c *gin.Context) {
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
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	documents, err := cc.service.GetAssociatedDocuments(uint(complianceID), uint(clauseID), &queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, documents)
}
func (ccc *ComplianceClauseController) CreateComplianceClauseDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO

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

	err = ccc.service.CreateComplianceClauseDocument(uint(complianceID), uint(clauseID), currentUserID, &input)
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

	c.Status(http.StatusCreated)
}

func (ccc *ComplianceClauseController) DeleteComplianceClauseDocument(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	clauseIDStr := c.Param("clauseId")
	documentIDStr := c.Param("documentId")
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

	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = ccc.service.DeleteComplianceClauseDocument(uint(clauseID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compliance Clause document mapping deleted successfully"})
}

func (cc *ComplianceClauseController) GetAssociatedIncident(c *gin.Context) {
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
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := cc.service.GetAssociatedIncident(uint(complianceID), uint(clauseID), &queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, risks)
}

func (ccc *ComplianceClauseController) CreateComplianceClauseIncident(c *gin.Context) {
	var input commonDTOs.AssociatedIncidentMapDTO

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

	err = ccc.service.CreateComplianceClauseIncident(uint(complianceID), uint(clauseID), currentUserID, &input)
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

	c.Status(http.StatusCreated)
}

func (ccc *ComplianceClauseController) DeleteComplianceClauseIncident(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	clauseIDStr := c.Param("clauseId")
	incidentIDStr := c.Param("incidentId")
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

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = ccc.service.DeleteComplianceClauseIncident(uint(clauseID), uint(incidentID))
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compliance Clause incident mapping deleted successfully"})
}

func (cc *ComplianceClauseController) GetAssociatedStrategicObjective(c *gin.Context) {
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
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := cc.service.GetAssociatedStrategicObjective(uint(complianceID), uint(clauseID), &queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, risks)
}

func (ccc *ComplianceClauseController) CreateComplianceClauseStrategicObjective(c *gin.Context) {
	var input commonDTOs.AssociatedStrategicObjectiveMapDTO

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

	err = ccc.service.CreateComplianceClauseStrategicObjective(uint(complianceID), uint(clauseID), currentUserID, &input)
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

	c.Status(http.StatusCreated)
}

func (ccc *ComplianceClauseController) DeleteComplianceClausStrategicObjective(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	clauseIDStr := c.Param("clauseId")
	strategicObjIDStr := c.Param("strategicObjectiveId")
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

	strategicObjID, err := strconv.ParseUint(strategicObjIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = ccc.service.DeleteComplianceClausStrategicObjective(uint(clauseID), uint(strategicObjID))
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compliance Clause strategic objective mapping deleted successfully"})
}

func (cc *ComplianceClauseController) GetAssociatedAssets(c *gin.Context) {
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
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := cc.service.GetAssociatedAssets(uint(complianceID), uint(clauseID), &queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, risks)
}

func (ccc *ComplianceClauseController) CreateComplianceClauseAssets(c *gin.Context) {
	var input commonDTOs.AssociatedAssetMapDTO

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

	err = ccc.service.CreateComplianceClauseAssets(uint(complianceID), uint(clauseID), currentUserID, &input)
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

	c.Status(http.StatusCreated)
}

func (ccc *ComplianceClauseController) DeleteComplianceAssets(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	clauseIDStr := c.Param("clauseId")
	assetIDStr := c.Param("assetId")
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

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = ccc.service.DeleteComplianceAssets(uint(clauseID), uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compliance Clause  asset mapping deleted successfully"})
}

func (cc *ComplianceClauseController) GetAssociatedProcesses(c *gin.Context) {
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
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := cc.service.GetAssociatedProcesses(uint(complianceID), uint(clauseID), &queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, risks)
}

func (ccc *ComplianceClauseController) CreateComplianceClausProcesses(c *gin.Context) {
	var input commonDTOs.AssociatedProcessMapDTO

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

	err = ccc.service.CreateComplianceClausProcesses(uint(complianceID), uint(clauseID), currentUserID, &input)
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

	c.Status(http.StatusCreated)
}

func (ccc *ComplianceClauseController) DeleteComplianceClausProcesses(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	clauseIDStr := c.Param("clauseId")
	processIDStr := c.Param("processId")
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

	processID, err := strconv.ParseUint(processIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = ccc.service.DeleteComplianceClausProcesses(uint(clauseID), uint(processID))
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compliance Clause process  mapping deleted successfully"})
}

func (cc *ComplianceClauseController) GetAssociatedKPI(c *gin.Context) {
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
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	risks, err := cc.service.GetAssociatedKPI(uint(complianceID), uint(clauseID), &queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, risks)
}

func (ccc *ComplianceClauseController) CreateComplianceClausKPI(c *gin.Context) {
	var input commonDTOs.AssociatedKPIMapDTO

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

	err = ccc.service.CreateComplianceClausKPI(uint(complianceID), uint(clauseID), currentUserID, &input)
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

	c.Status(http.StatusCreated)
}

func (ccc *ComplianceClauseController) DeleteComplianceClausKPI(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")
	clauseIDStr := c.Param("clauseId")
	kpiStr := c.Param("kpiId")
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

	kpiID, err := strconv.ParseUint(kpiStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = ccc.service.DeleteComplianceClausKPI(uint(clauseID), uint(kpiID))
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compliance Clause KPI  mapping deleted successfully"})
}

func (ccc *ComplianceClauseController) GetClausesTemplate(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	_, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	file, err := ccc.service.GenerateClauseTemplate()
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, "Failed to generate template")
		return
	}
	defer file.Close()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", `attachment; filename="ComplianceClauseTemplate.xlsx"`)
	c.Header("Content-Type", "application/octet-stream")
	c.File(file.Name())
}

func (ccc *ComplianceClauseController) ValidateClauseTemplate(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, "Failed to read uploaded file")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, "Failed to open file")
		return
	}
	defer file.Close()

	validationResult, err := ccc.service.ValidateClauseTemplate(uint(complianceID), file)
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

	c.JSON(http.StatusOK, validationResult)
}

func (ccc *ComplianceClauseController) ValidateImportClause(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var input []map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	validationResult, err := ccc.service.ValidateImportClause(uint(complianceID), input)
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

	c.JSON(http.StatusOK, validationResult)
}


func (crc *ComplianceClauseController) ImportClauses(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var input []map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	createdIDs, err := crc.service.ImportClauses(uint(complianceID), input, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"created_ids": createdIDs})
}

func (ccc *ComplianceClauseController) GetControlMappingTemplate(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	file, err := ccc.service.GenerateComplianceClauseControlTemplate(uint(complianceID))
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, "Failed to generate control mapping template")
		return
	}
	defer file.Close()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", `attachment; filename="ControlMappingTemplate.xlsx"`)
	c.Header("Content-Type", "application/octet-stream")
	c.File(file.Name())

}

func (ccc *ComplianceClauseController) ValidateControlMappingTemplate(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, "Failed to read uploaded file")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, "Failed to open file")
		return
	}
	defer file.Close()

	validationResult, err := ccc.service.ValidateControlMappingTemplate(uint(complianceID), file)
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

	c.JSON(http.StatusOK, validationResult)
}

func (ccc *ComplianceClauseController) ValidateImportControlMapping(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var input []map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	validationResult, err := ccc.service.ValidateImportControlMapping(uint(complianceID), input)
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

	c.JSON(http.StatusOK, validationResult)
}

func (ccc *ComplianceClauseController) ImportControlMapping(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var input []map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ccc.db, c, err.Error())
		return
	}

	createdIDs, err := ccc.service.ImportControlMapping(uint(complianceID), input, currentUserID)
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

	c.JSON(http.StatusOK, map[string]interface{}{
		"created_ids": createdIDs,
	})
}
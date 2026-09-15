package controllers

import (
	"errors"
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSecurityGateController struct {
	service services.ProjectSecurityGateServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateController(service services.ProjectSecurityGateServiceInterface) *ProjectSecurityGateController {
	return &ProjectSecurityGateController{
		service: service,
	}
}

func (psgc *ProjectSecurityGateController) GetAllProjectSecurityGates(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "project_title_en")
	keyword := c.Query("keyword")

	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	validSortFields := map[string]string{
		"project_title_en":             "project_security_gates.project_title_en",
		"project_title_ar":             "project_security_gates.project_title_ar",
		"stage_title_en":               "stages.title_en",
		"stage_title_ar":               "stages.title_ar",
		"business_owner_first_name_en": "business_owner.first_name_en",
		"business_owner_first_name_ar": "business_owner.first_name_ar",
		"readiness_status_title_en":    "readiness_status.title_en",
		"readiness_status_title_ar":    "readiness_status.title_ar",
		"created_at":                   "project_security_gates.created_at",
		"updated_at":                   "project_security_gates.updated_at",
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, "Invalid page_size parameter")
		return
	}

	exportFromContext, exists := c.Get("export")
	var export bool
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid 'export' parameter"))
			return
		}
	}

	queryParams := dtos.ProjectSecurityGateListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Export:   export,
	}

	projectSecurityGates, err := psgc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(psgc.db, c, err.Error())
			return
		}

		excelFile, err := psgc.service.ExportProjectSecurityGates(projectSecurityGates.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(psgc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Project Security Gates.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, projectSecurityGates)
}

func (psgc *ProjectSecurityGateController) GetProjectSecurityGatesExport(c *gin.Context) {
	c.Set("export", true)
	psgc.GetAllProjectSecurityGates(c)
}

func (psgc *ProjectSecurityGateController) GetProjectSecurityGatesByReadinessStatus(c *gin.Context) {
	readinessStatusIDStr := c.Param("readinessStatusId")
	readinessStatusID, err := strconv.ParseUint(readinessStatusIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid 'readinessStatusId' parameter"))
		return
	}

	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "project_title_en")
	keyword := c.Query("keyword")

	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	validSortFields := map[string]string{
		"project_title_en":             "project_security_gates.project_title_en",
		"project_title_ar":             "project_security_gates.project_title_ar",
		"stage_title_en":               "stages.title_en",
		"stage_title_ar":               "stages.title_ar",
		"business_owner_first_name_en": "business_owner.first_name_en",
		"business_owner_first_name_ar": "business_owner.first_name_ar",
		"readiness_status_title_en":    "readiness_status.title_en",
		"readiness_status_title_ar":    "readiness_status.title_ar",
		"created_at":                   "project_security_gates.created_at",
		"updated_at":                   "project_security_gates.updated_at",
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, "Invalid page_size parameter")
		return
	}

	exportFromContext, exists := c.Get("export")
	var export bool
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid 'export' parameter"))
			return
		}
	}

	queryParams := dtos.ProjectSecurityGateListQueryParamsDTO{
		Sort:              sort,
		SortBy:            dbField,
		Keyword:           keyword,
		Page:              page,
		PageSize:          pageSize,
		Export:            export,
		ReadinessStatusID: uint(readinessStatusID),
	}

	projectSecurityGates, err := psgc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(psgc.db, c, err.Error())
			return
		}

		excelFile, err := psgc.service.ExportProjectSecurityGates(projectSecurityGates.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(psgc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		fileName := psgc.service.GetReadinessStatusExportFilename(uint(readinessStatusID))
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="`+fileName+`"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}

	c.JSON(http.StatusOK, projectSecurityGates)
}

func (psgc *ProjectSecurityGateController) GetProjectSecurityGatesByReadinessStatusExport(c *gin.Context) {
	c.Set("export", true)
	psgc.GetProjectSecurityGatesByReadinessStatus(c)
}

func (psgc *ProjectSecurityGateController) GetProjectSecurityGatesByStage(c *gin.Context) {
	stageIDStr := c.Param("stageId")
	stageID, err := strconv.ParseUint(stageIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid 'stageId' parameter"))
		return
	}

	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "project_title_en")
	keyword := c.Query("keyword")

	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	validSortFields := map[string]string{
		"project_title_en":             "project_security_gates.project_title_en",
		"project_title_ar":             "project_security_gates.project_title_ar",
		"stage_title_en":               "stages.title_en",
		"stage_title_ar":               "stages.title_ar",
		"business_owner_first_name_en": "business_owner.first_name_en",
		"business_owner_first_name_ar": "business_owner.first_name_ar",
		"readiness_status_title_en":    "readiness_status.title_en",
		"readiness_status_title_ar":    "readiness_status.title_ar",
		"created_at":                   "project_security_gates.created_at",
		"updated_at":                   "project_security_gates.updated_at",
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, "Invalid page_size parameter")
		return
	}

	exportFromContext, exists := c.Get("export")
	var export bool
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid 'export' parameter"))
			return
		}
	}

	queryParams := dtos.ProjectSecurityGateListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Export:   export,
		StageID:  uint(stageID),
	}

	projectSecurityGates, err := psgc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(psgc.db, c, err.Error())
			return
		}

		excelFile, err := psgc.service.ExportProjectSecurityGates(projectSecurityGates.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(psgc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		fileName := psgc.service.GetStageExportFilename(uint(stageID))
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="`+fileName+`"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}

	c.JSON(http.StatusOK, projectSecurityGates)
}

func (psgc *ProjectSecurityGateController) GetProjectSecurityGatesByStageExport(c *gin.Context) {
	c.Set("export", true)
	psgc.GetProjectSecurityGatesByStage(c)
}

func (psgc *ProjectSecurityGateController) GetProjectSecurityGatesByClassification(c *gin.Context) {
	classificationIDStr := c.Param("classificationId")
	classificationID, err := strconv.ParseUint(classificationIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid 'classificationId' parameter"))
		return
	}

	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "project_title_en")
	keyword := c.Query("keyword")

	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	validSortFields := map[string]string{
		"project_title_en":             "project_security_gates.project_title_en",
		"project_title_ar":             "project_security_gates.project_title_ar",
		"stage_title_en":               "stages.title_en",
		"stage_title_ar":               "stages.title_ar",
		"business_owner_first_name_en": "business_owner.first_name_en",
		"business_owner_first_name_ar": "business_owner.first_name_ar",
		"readiness_status_title_en":    "readiness_status.title_en",
		"readiness_status_title_ar":    "readiness_status.title_ar",
		"created_at":                   "project_security_gates.created_at",
		"updated_at":                   "project_security_gates.updated_at",
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, "Invalid page_size parameter")
		return
	}

	exportFromContext, exists := c.Get("export")
	var export bool
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInvalidErrResp(c, errors.New("invalid 'export' parameter"))
			return
		}
	}

	queryParams := dtos.ProjectSecurityGateListQueryParamsDTO{
		Sort:             sort,
		SortBy:           dbField,
		Keyword:          keyword,
		Page:             page,
		PageSize:         pageSize,
		Export:           export,
		ClassificationID: uint(classificationID),
	}

	projectSecurityGates, err := psgc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(psgc.db, c, err.Error())
			return
		}

		excelFile, err := psgc.service.ExportProjectSecurityGates(projectSecurityGates.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(psgc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		fileName := psgc.service.GetClassificationExportFilename(uint(classificationID))
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="`+fileName+`"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())
		return
	}

	c.JSON(http.StatusOK, projectSecurityGates)
}

func (psgc *ProjectSecurityGateController) GetProjectSecurityGatesByClassificationExport(c *gin.Context) {
	c.Set("export", true)
	psgc.GetProjectSecurityGatesByClassification(c)
}

func (psgc *ProjectSecurityGateController) GetProjectSecurityGateSummary(c *gin.Context) {
	summary, err := psgc.service.GetSummary()
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (psgc *ProjectSecurityGateController) GetProjectSecurityGateDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	projectSecurityGate, err := psgc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, projectSecurityGate)
}

func (psgc *ProjectSecurityGateController) CreateProjectSecurityGate(c *gin.Context) {
	var form dtos.ProjectSecurityGateFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	id, err := psgc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (psgc *ProjectSecurityGateController) UpdateProjectSecurityGate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ProjectSecurityGateFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	err = psgc.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psgc *ProjectSecurityGateController) GoLiveProjectSecurityGate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	err = psgc.service.GoLive(uint(id), currentUserID)
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
			utility.SendInternalErrResp(psgc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psgc *ProjectSecurityGateController) DeleteProjectSecurityGates(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	failedIDs, err := psgc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(psgc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

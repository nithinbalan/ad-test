package controllers

import (
	"errors"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemRegisterController struct {
	service services.SystemRegisterServiceInterface
	db      *gorm.DB
}

func NewSystemRegisterController(service services.SystemRegisterServiceInterface) *SystemRegisterController {
	return &SystemRegisterController{
		service: service,
	}
}

func (src *SystemRegisterController) GetAllSystemRegisters(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "reference_number")
	keyword := c.Query("keyword")

	// Validate sort_order to prevent SQL injection
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(src.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"business_owner_unit_ids",
		"it_responsible_unit_ids",
		"process_operation_team_ids",
		"process_recovery_team_ids",
		"recovery_manager_ids",
	}

	validSortFields := map[string]string{
		"reference_number":                 "system_registers.reference_number",
		"title_en":                         "system_registers.title_en",
		"title_ar":                         "system_registers.title_ar",
		"number_of_people_necessary":       "system_registers.number_of_people_necessary",
		"third_party_dependencies":         "system_registers.third_party_dependencies",
		"outsourcing_partners":             "system_registers.outsourcing_partners",
		"customer_user":                    "system_registers.customer_user",
		"business_owner_unit_en":           "business_owner_unit.title_en",
		"business_owner_unit_ar":           "business_owner_unit.title_ar",
		"it_responsible_unit_en":           "it_responsible_unit.title_en",
		"it_responsible_unit_ar":           "it_responsible_unit.title_ar",
		"process_operation_team_title_en":  "process_operation_team.title_en",
		"process_operation_team_title_ar":  "process_operation_team.title_ar",
		"process_recovery_team_title_en":   "process_recovery_team.title_en",
		"process_recovery_team_title_ar":   "process_recovery_team.title_ar",
		"recovery_manager_en":              "recovery_manager.first_name_en",
		"recovery_manager_ar":              "recovery_manager.first_name_ar",
		"recovery_manager_contact_details": "system_registers.recovery_manager_contact_details",
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

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	params := dtos.SystemRegisterListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
		Export:   export,
	}

	systemRegisters, err := src.service.GetAll(&params)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(src.db, c, err.Error())
			return
		}

		excelFile, err := src.service.ExportSystemRegister(systemRegisters.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(src.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="SystemRegisters.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, systemRegisters)
}

func (src *SystemRegisterController) GetSystemRegisterDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	systemRegister, err := src.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, systemRegister)
}

func (src *SystemRegisterController) CreateSystemRegister(c *gin.Context) {
	var systemRegisterForm dtos.SystemRegisterFormDTO
	if err := c.ShouldBindJSON(&systemRegisterForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	systemRegisterID, err := src.service.Create(&systemRegisterForm, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(src.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": systemRegisterID})
}

func (src *SystemRegisterController) UpdateSystemRegister(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var systemRegisterForm dtos.SystemRegisterFormDTO
	if err := c.ShouldBindJSON(&systemRegisterForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	err = src.service.Update(uint(id), &systemRegisterForm, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(src.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (src *SystemRegisterController) DeleteSystemRegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	failedIDs, err := src.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(src.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (src *SystemRegisterController) GetSystemRegistersExport(c *gin.Context) {
	c.Set("export", true)
	src.GetAllSystemRegisters(c)
}

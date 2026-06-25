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
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	allowedParams := []string{
		"service_importance_ids",
		"responsible_it_team_ids",
		"business_owner_ids",
		"process_operation_team_ids",
		"process_recovery_team_ids",
		"impact_type_ids",
	}

	validSortFields := map[string]string{
		"title_en":                         "system_registers.title_en",
		"title_ar":                         "system_registers.title_ar",
		"service_description":              "system_registers.service_description",
		"modules_components":               "system_registers.modules_components",
		"third_party_dependencies":         "system_registers.third_party_dependencies",
		"service_importance_title_en":      "service_importance.title_en",
		"service_importance_title_ar":      "service_importance.title_ar",
		"responsible_it_team_title_en":     "responsible_it_team.title_en",
		"responsible_it_team_title_ar":     "responsible_it_team.title_ar",
		"outsourcing_partner":              "system_registers.outsourcing_partner",
		"operation_location_title_en":      "operation_location.title_en",
		"operation_location_title_ar":      "operation_location.title_ar",
		"business_owner_title_en":          "business_owner.title_en",
		"business_owner_title_ar":          "business_owner.title_ar",
		"customers":                        "system_registers.customers",
		"disruption_impact_level_title_en": "disruption_impact_level.title_en",
		"disruption_impact_level_title_ar": "disruption_impact_level.title_ar",
		"rto":                              "system_registers.rto",
		"rto_unit_title_en":                "rto_unit.title_en",
		"rto_unit_title_ar":                "rto_unit.title_ar",
		"mtpd":                             "system_registers.mtpd",
		"mtpd_unit_title_en":               "mtpd_unit.title_en",
		"mtpd_unit_title_ar":               "mtpd_unit.title_ar",
		"rpo":                              "system_registers.rpo",
		"rpo_unit_title_en":                "rpo_unit.title_en",
		"rpo_unit_title_ar":                "rpo_unit.title_ar",
		"it_processes":                     "system_registers.it_processes",
		"it_service_owner_title_en":        "it_service_owner.title_en",
		"it_service_owner_title_ar":        "it_service_owner.title_ar",
		"process_operation_team_title_en":  "process_operation_team.title_en",
		"process_operation_team_title_ar":  "process_operation_team.title_ar",
		"process_recovery_team_title_en":   "process_recovery_team.title_en",
		"process_recovery_team_title_ar":   "process_recovery_team.title_ar",
		"recovery_manager_title_en":        "recovery_manager.title_en",
		"recovery_manager_title_ar":        "recovery_manager.title_ar",
		"recovery_manager_phone":           "system_registers.recovery_manager_phone",
		"recovery_manager_email":           "system_registers.recovery_manager_email",
		"recovery_location":                "system_registers.recovery_location",
		"vendor_rto":                       "system_registers.vendor_rto",
		"vendor_rto_unit_title_en":         "vendor_rto_unit.title_en",
		"vendor_rto_unit_title_ar":         "vendor_rto_unit.title_ar",
		"created_at":                       "system_registers.created_at",
		"updated_at":                       "system_registers.updated_at",
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
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	exportFromContext, exists := c.Get("export")
	var export bool
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, _ = strconv.ParseBool(c.DefaultQuery("export", "false"))
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

func (src *SystemRegisterController) GetSystemRegistersExport(c *gin.Context) {
	c.Set("export", true)
	src.GetAllSystemRegisters(c)
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

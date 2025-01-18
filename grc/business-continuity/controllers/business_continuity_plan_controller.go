package controllers

import (
	"errors"
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

type BusinessContinuityPlanController struct {
	service services.BusinessContinuityPlanServiceInterface
	db      *gorm.DB
}

func NewBusinessContinuityPlanController(service services.BusinessContinuityPlanServiceInterface) *BusinessContinuityPlanController {
	return &BusinessContinuityPlanController{
		service: service,
	}
}

func (bfc *BusinessContinuityPlanController) GetAllBusinessContinuityPlans(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "business_function_title")
	keyword := c.Query("keyword")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"business_function_ids",
	}
	validSortFields := map[string]string{
		"id":                      "id",
		"reference_number":        "business_continuity_plans.reference_number",
		"business_function_title": "business_functions.title",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := businessContinuityDTOs.BusinessContinuityPlanListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
		Export:   export,
	}
	businessContinuityPlans, err := bfc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, err.Error())
		return
	}

	if export {

		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(bfc.db, c, err.Error())
			return
		}

		excelFile, err := bfc.service.ExportBusinessContinuityPlan(businessContinuityPlans.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(bfc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="BusinessContinuityPlan.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, businessContinuityPlans)
}

func (bfc *BusinessContinuityPlanController) GetBusinessContinuityPlanDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	businessContinuityPlan, err := bfc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(bfc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, businessContinuityPlan)
}

func (bfc *BusinessContinuityPlanController) CreateBusinessContinuityPlan(c *gin.Context) {
	var businessContinuityPlanFormDTO businessContinuityDTOs.BusinessContinuityPlanFormDTO

	if err := c.ShouldBindJSON(&businessContinuityPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, err.Error())
		return
	}

	id, err := bfc.service.Create(currentUserID, &businessContinuityPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(bfc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (bfc *BusinessContinuityPlanController) UpdateBusinessContinuityPlan(c *gin.Context) {
	var businessContinuityPlanFormDTO businessContinuityDTOs.BusinessContinuityPlanFormDTO

	businessContinuityPlanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&businessContinuityPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, err.Error())
		return
	}

	err = bfc.service.Update(uint(businessContinuityPlanID), currentUserID, &businessContinuityPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(bfc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (bfc *BusinessContinuityPlanController) DeleteBusinessContinuityPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, err.Error())
		return
	}

	failedIDs, err := bfc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(bfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (bfc *BusinessContinuityPlanController) GetAllBusinessContinuityPlanExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	bfc.GetAllBusinessContinuityPlans(c)
}

func (crc *BusinessContinuityPlanController) GetBusinessContinuityPlanDocumentPreview(c *gin.Context) {
	idStr := c.Param("id")
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

	bcpDocumentPath, err := crc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(crc.db, c, err.Error())
		}
		return
	}
	c.File(bcpDocumentPath)
}

func (crc *BusinessContinuityPlanController) DownloadBusinessContinuityPlanDocument(c *gin.Context) {
	idStr := c.Param("id")
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

	bcpDocumentPath, err := crc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(crc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(bcpDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(bcpDocumentPath)
}

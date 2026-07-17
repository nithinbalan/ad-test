package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/violation-report/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ViolationReportAssociatedItemController struct {
	service services.ViolationReportAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewViolationReportAssociatedItemController(service services.ViolationReportAssociatedItemServiceInterface) *ViolationReportAssociatedItemController {
	return &ViolationReportAssociatedItemController{
		service: service,
	}
}

func (iaic *ViolationReportAssociatedItemController) GetAllViolationReportControls(c *gin.Context) {

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	violationReportControls, err := iaic.service.GetAllControls(uint(violationReportID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, violationReportControls)
}

func (iaic *ViolationReportAssociatedItemController) GetAllViolationReportComplianceRegisters(c *gin.Context) {

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	violationReportComplianceRequirements, err := iaic.service.GetAllComplianceRegisters(uint(violationReportID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, violationReportComplianceRequirements)
}

func (iaic *ViolationReportAssociatedItemController) GetAllViolationReportRisks(c *gin.Context) {

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	valueStr := c.Query("value")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page_size parameter")
		return
	}
	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(iaic.db, c, "Invalid value parameter")
			return
		}
	}
	allowedParams := []string{
		"current_likelihood_ids",
		"current_impact_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())
	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Value:    value,
		Filters:  filters,
	}
	violationReportRisks, err := iaic.service.GetAllRisks(uint(violationReportID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, violationReportRisks)
}

func (ic *ViolationReportAssociatedItemController) GetAllViolationReportDocuments(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(ic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(ic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	violationReportDocuments, err := ic.service.GetAllDocuments(uint(violationReportID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, violationReportDocuments)
}

func (iaic *ViolationReportAssociatedItemController) GetAllViolationReportAssets(c *gin.Context) {

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := commonDTOs.AssociatedItemListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	violationReportAssets, err := iaic.service.GetAllAssets(uint(violationReportID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, violationReportAssets)
}

func (iaic *ViolationReportAssociatedItemController) CreateViolationReportRisk(c *gin.Context) {
	var input commonDTOs.AssociatedRiskMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	violationReportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	err = iaic.service.CreateViolationReportRisk(uint(violationReportID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (iaic *ViolationReportAssociatedItemController) DeleteViolationReportRisk(c *gin.Context) {
	violationReportID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 64)
	err := iaic.service.DeleteViolationReportRisk(uint(violationReportID), uint(riskID))
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Risk unmapped from violation report successfully"})
}

func (iaic *ViolationReportAssociatedItemController) CreateViolationReportDocument(c *gin.Context) {
	var input commonDTOs.AssociatedDocumentMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	violationReportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	err = iaic.service.CreateViolationReportDocument(uint(violationReportID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (iaic *ViolationReportAssociatedItemController) DeleteViolationReportDocument(c *gin.Context) {
	violationReportID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	documentID, _ := strconv.ParseUint(c.Param("documentId"), 10, 64)
	err := iaic.service.DeleteViolationReportDocument(uint(violationReportID), uint(documentID))
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document unmapped from violation report successfully"})
}

func (iaic *ViolationReportAssociatedItemController) CreateViolationReportControl(c *gin.Context) {
	var input commonDTOs.AssociatedControlMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	violationReportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	err = iaic.service.CreateViolationReportControl(uint(violationReportID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (iaic *ViolationReportAssociatedItemController) DeleteViolationReportControl(c *gin.Context) {
	violationReportID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	controlID, _ := strconv.ParseUint(c.Param("controlId"), 10, 64)
	err := iaic.service.DeleteViolationReportControl(uint(violationReportID), uint(controlID))
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Control unmapped from violation report successfully"})
}

func (iaic *ViolationReportAssociatedItemController) CreateViolationReportAsset(c *gin.Context) {
	var input commonDTOs.AssociatedAssetMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	violationReportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	err = iaic.service.CreateViolationReportAsset(uint(violationReportID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusCreated)
}

func (iaic *ViolationReportAssociatedItemController) DeleteViolationReportAsset(c *gin.Context) {

	violationReportIDStr := c.Param("id")
	assetIDStr := c.Param("assetId")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	err = iaic.service.DeleteViolationReportAsset(uint(violationReportID), uint(assetID))
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset unmapped from violation report successfully"})
}

func (iaic *ViolationReportAssociatedItemController) CreateViolationReportComplianceRegister(c *gin.Context) {
	var input commonDTOs.AssociatedComplianceRegisterMapDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utility.SendInvalidErrResp(c, utility.NewValidationError(err, "invalid_input", "request_body"))
		return
	}

	violationReportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	err = iaic.service.CreateViolationReportComplianceRegister(uint(violationReportID), currentUserID, &input)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iaic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Compliance registers mapped to violation report successfully"})
}

func (iaic *ViolationReportAssociatedItemController) DeleteViolationReportComplianceRegister(c *gin.Context) {
	violationReportID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	complianceRegisterID, _ := strconv.ParseUint(c.Param("complianceRegisterId"), 10, 64)
	err := iaic.service.DeleteViolationReportComplianceRegister(uint(violationReportID), uint(complianceRegisterID))
	if err != nil {
		utility.SendInternalErrResp(iaic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compliance register unmapped from violation report successfully"})
}

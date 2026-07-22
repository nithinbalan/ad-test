package controllers

import (
	"errors"
	awarenessDTOs "jethur/grc/awareness-reporting/dtos"
	"jethur/grc/awareness-reporting/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AwarenessCampaignController struct {
	service services.AwarenessCampaignServiceInterface
	db      *gorm.DB
}

func NewAwarenessCampaignController(service services.AwarenessCampaignServiceInterface) *AwarenessCampaignController {
	return &AwarenessCampaignController{
		service: service,
	}
}

func (acc *AwarenessCampaignController) GetAllAwarenessCampaigns(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "campaign_title")
	keyword := c.Query("keyword")

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(acc.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"awareness_campaign_type_ids",
		"awareness_priority_ids",
		"awareness_campaign_status_ids",
		"organization_unit_ids",
		"owner_ids",
	}

	validSortFields := map[string]string{
		"campaign_title":                     "awareness_campaigns.campaign_title",
		"reference_number":                   "awareness_campaigns.reference_number",
		"objective":                          "awareness_campaigns.objective",
		"target_audience":                    "awareness_campaigns.target_audience",
		"start_date":                         "awareness_campaigns.start_date",
		"end_date":                           "awareness_campaigns.end_date",
		"awareness_campaign_type_title_en":   "awareness_campaign_types.title_en",
		"awareness_campaign_type_title_ar":   "awareness_campaign_types.title_ar",
		"awareness_priority_title_en":        "awareness_priorities.title_en",
		"awareness_priority_title_ar":        "awareness_priorities.title_ar",
		"awareness_campaign_status_title_en": "awareness_campaign_statuses.title_en",
		"awareness_campaign_status_title_ar": "awareness_campaign_statuses.title_ar",
		"organization_unit_title_en":         "organization_units.title_en",
		"organization_unit_title_ar":         "organization_units.title_ar",
		"owner_first_name_en":                "owner_user.first_name_en",
		"owner_first_name_ar":                "owner_user.first_name_ar",
		"owner_last_name_en":                 "owner_user.last_name_en",
		"owner_last_name_ar":                 "owner_user.last_name_ar",
		"owner_email":                        "owner_user.email",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, "Invalid page_size parameter")
		return
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := awarenessDTOs.AwarenessCampaignListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
		Export:   export,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	campaigns, err := acc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(acc.db, c, err.Error())
			return
		}

		excelFile, err := acc.service.ExportAwarenessCampaign(campaigns.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(acc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="AwarenessCampaigns.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, campaigns)
}

func (acc *AwarenessCampaignController) GetAwarenessCampaignDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	campaign, err := acc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, campaign)
}

func (acc *AwarenessCampaignController) CreateAwarenessCampaign(c *gin.Context) {
	var formDTO awarenessDTOs.AwarenessCampaignFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	id, err := acc.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (acc *AwarenessCampaignController) UpdateAwarenessCampaign(c *gin.Context) {
	var formDTO awarenessDTOs.AwarenessCampaignFormDTO

	campaignID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	err = acc.service.Update(uint(campaignID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (acc *AwarenessCampaignController) CompleteAwarenessCampaign(c *gin.Context) {
	campaignID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	if err := acc.service.CompleteCampaign(uint(campaignID), currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Awareness campaign completed successfully"})
}

func (acc *AwarenessCampaignController) CancelAwarenessCampaign(c *gin.Context) {
	campaignID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	if err := acc.service.CancelCampaign(uint(campaignID), currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(acc.db, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Awareness campaign cancelled successfully"})
}

func (acc *AwarenessCampaignController) DeleteAwarenessCampaigns(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	failedIDs, err := acc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(acc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (acc *AwarenessCampaignController) GetAwarenessCampaignDocumentPreview(c *gin.Context) {
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

	documentPath, err := acc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(acc.db, c, err.Error())
		}
		return
	}

	c.File(documentPath)
}

func (acc *AwarenessCampaignController) DownloadAwarenessCampaignDocument(c *gin.Context) {
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

	documentPath, err := acc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(acc.db, c, err.Error())
		}
		return
	}

	filename := filepath.Base(documentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}

func (acc *AwarenessCampaignController) GetAwarenessCampaignExport(c *gin.Context) {
	c.Set("export", true)

	acc.GetAllAwarenessCampaigns(c)
}

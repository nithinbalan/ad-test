package controllers

import (
	"errors"
	"jethur/grc/common/dtos"
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InternalCommitteeController struct {
	service services.InternalCommitteeServiceInterface
	db      *gorm.DB
}

func NewInternalCommitteeController(service services.InternalCommitteeServiceInterface) *InternalCommitteeController {
	return &InternalCommitteeController{
		service: service,
	}
}

func (icc *InternalCommitteeController) GetAllInternalCommittees(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	exportFromContext, exists := c.Get("export")
	var export bool
	var err error

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(icc.db, c, "Invalid export parameter")
			return
		}
	}

	// Define allowed parameter names for filters
	allowedParams := []string{
		"committee_type_ids",
		"committee_term_ids",
		"committee_classification_ids",
		"committee_status_ids",
		"requesting_entity_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse date filters
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"reference_number":           "internal_committees.reference_number",
		"requesting_entity_title_en": "requesting_entity.title_en",
		"requesting_entity_title_ar": "requesting_entity.title_ar",
		"committee_name_en":          "latest_version.committee_name_en",
		"committee_name_ar":          "latest_version.committee_name_ar",
		"committee_type_title_en":    "committee_types.title_en",
		"committee_type_title_ar":    "committee_types.title_ar",
		"committee_term_title_en":    "committee_terms.title_en",
		"committee_term_title_ar":    "committee_terms.title_ar",
		"committee_status_title_en":  "internal_committee_statuses.title_en",
		"committee_status_title_ar":  "internal_committee_statuses.title_ar",
		"created_at":                 "internal_committees.created_at",
		"updated_at":                 "internal_committees.updated_at",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := organizationDTOs.InternalCommitteeListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		Export:      export,
	}

	internalCommittees, err := icc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	// If export is true, generate and return the Excel file
	if export {
		// TODO: Implement export functionality
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Export not implemented yet"})
		return
	}

	// Otherwise, return the data as JSON
	c.JSON(http.StatusOK, internalCommittees)
}

func (icc *InternalCommitteeController) GetInternalCommitteeDetails(c *gin.Context) {
	idStr := c.Param("committeeId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	internalCommittee, err := icc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utility.SendNotFoundErrResp(c)
			return
		}
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, internalCommittee)
}

func (icc *InternalCommitteeController) CreateInternalCommittee(c *gin.Context) {
	var internalCommitteeFormDTO organizationDTOs.InternalCommitteeFormDTO

	if err := c.ShouldBindJSON(&internalCommitteeFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	committeeID, err := icc.service.Create(currentUserID, &internalCommitteeFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: committeeID})
}

func (icc *InternalCommitteeController) UpdateInternalCommittee(c *gin.Context) {
	idStr := c.Param("committeeId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var internalCommitteeFormDTO organizationDTOs.InternalCommitteeFormDTO

	if err := c.ShouldBindJSON(&internalCommitteeFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	err = icc.service.Update(uint(id), currentUserID, &internalCommitteeFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (icc *InternalCommitteeController) ReformInternalCommittee(c *gin.Context) {
	var reformDTO organizationDTOs.ReformInternalCommitteeDTO

	if err := c.ShouldBindJSON(&reformDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	versionID, err := icc.service.ReformCommittee(currentUserID, &reformDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: versionID})
}

func (icc *InternalCommitteeController) DeleteInternalCommittee(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get current user
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	// Call service
	failedIDs, err := icc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	// Response (same as ControlController)
	c.JSON(http.StatusOK, gin.H{
		"failed_ids": failedIDs,
	})
}

// Participant Methods

func (icc *InternalCommitteeController) CreateInternalCommitteeParticipants(c *gin.Context) {
	nomineeDeptIDStr := c.Param("nomineeDepartmentId")

	nomineeDeptID, err := strconv.ParseUint(nomineeDeptIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var participantFormDTO organizationDTOs.InternalCommitteeParticipantFormDTO

	if err := c.ShouldBindJSON(&participantFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	createdIDs, err := icc.service.CreateParticipants(currentUserID, uint(nomineeDeptID), &participantFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, dtos.IDResponse{ID: createdIDs[0]})
}

func (icc *InternalCommitteeController) DeleteInternalCommitteeParticipant(c *gin.Context) {
	nomineeDeptIDStr := c.Param("nomineeDepartmentId")

	nomineeDeptID, err := strconv.ParseUint(nomineeDeptIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	participantIDStr := c.Param("participantId")

	participantID, err := strconv.ParseUint(participantIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	err = icc.service.DeleteParticipant(uint(nomineeDeptID), uint(participantID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Participants deleted successfully",
	})
}

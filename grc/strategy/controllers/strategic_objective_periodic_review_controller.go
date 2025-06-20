package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	strategicObjectiveDTOs "jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategicObjectivePeriodicReviewController struct {
	service services.StrategicObjectivePeriodicReviewServiceInterface
	db      *gorm.DB
}

func NewStrategicObjectivePeriodicReviewController(service services.StrategicObjectivePeriodicReviewServiceInterface) *StrategicObjectivePeriodicReviewController {
	return &StrategicObjectivePeriodicReviewController{
		service: service,
	}
}

// Get all strategicObjective periodic reviews with filtering, sorting, and pagination
func (soprc *StrategicObjectivePeriodicReviewController) GetAllStrategicObjectivePeriodicReviews(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "review_date")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"review_status_ids",
	}

	validSortFields := map[string]string{
		"review_date": "strategic_objective_reviews.review_date",
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

	params := strategicObjectiveDTOs.StrategicObjectivePeriodicReviewQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	strategicObjectiveIdStr := c.Param("id")
	strategicObjectiveId, err := strconv.ParseUint(strategicObjectiveIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicObjectiveReviews, err := soprc.service.GetAll(&params, uint(strategicObjectiveId))
	if err != nil {
		utility.SendInternalErrResp(soprc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, strategicObjectiveReviews)
}

// Get details of a single strategicObjective periodic review by ID
func (soprc *StrategicObjectivePeriodicReviewController) GetStrategicObjectivePeriodicReviewDetails(c *gin.Context) {
	idStr := c.Param("reviewId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := soprc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(soprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, details)
}

// Create a new strategicObjective periodic review
func (soprc *StrategicObjectivePeriodicReviewController) CreateStrategicObjectivePeriodicReview(c *gin.Context) {
	var formDTO strategicObjectiveDTOs.StrategicObjectivePeriodicReviewFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	strategicObjectiveIdStr := c.Param("id")
	strategicObjectiveId, err := strconv.ParseUint(strategicObjectiveIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soprc.db, c, err.Error())
		return
	}

	id, err := soprc.service.Create(currentUserID, uint(strategicObjectiveId), &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		case *utility.ResourceConflictError:
			utility.SendConflictErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(soprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

// Update an existing strategicObjective periodic review
func (soprc *StrategicObjectivePeriodicReviewController) UpdateStrategicObjectivePeriodicReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("reviewId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicObjectiveIdStr := c.Param("id")
	strategicObjectiveId, err := strconv.ParseUint(strategicObjectiveIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO strategicObjectiveDTOs.StrategicObjectivePeriodicReviewFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soprc.db, c, err.Error())
		return
	}

	err = soprc.service.Update(uint(id), uint(strategicObjectiveId), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ResourceConflictError:
			utility.SendConflictErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(soprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

// Delete strategicObjective periodic reviews by IDs
func (soprc *StrategicObjectivePeriodicReviewController) DeleteStrategicObjectivePeriodicReview(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soprc.db, c, err.Error())
		return
	}

	failedIDs, err := soprc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(soprc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

// Download a strategicObjective periodic review file
func (soprc *StrategicObjectivePeriodicReviewController) DownloadStrategicObjectivePeriodicReviewDocument(c *gin.Context) {
	idStr := c.Param("reviewId")
	reviewID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicObjectivePath, err := soprc.service.GetDocumentPath(uint(reviewID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(soprc.db, c, err.Error())
		}
		return
	}

	filename := filepath.Base(strategicObjectivePath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(strategicObjectivePath)
}

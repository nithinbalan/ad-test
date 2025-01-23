package controllers

import (
	"errors"
	"net/http"
	"strconv"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BcpTeamController struct {
	service services.BcpTeamServiceInterface
	db      *gorm.DB
}

func NewBcpTeamController(service services.BcpTeamServiceInterface) *BcpTeamController {
	return &BcpTeamController{
		service: service,
	}
}

func (brc *BcpTeamController) GetAllBcpTeams(c *gin.Context) {
	bcpIdStr := c.Param("bcpId")
	bcpId, err := strconv.ParseUint(bcpIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "bcp_Teams.responsibility")
	keyword := c.Query("keyword")

	allowedParams := []string{"responsible_user_id", "backup_user_id", "responsible_unit_id"}
	validSortFields := map[string]string{
		"responsibility":                 "bcp_Teams.responsibility",
		"responsible_user_first_name":    "responsible_user.first_name_en",
		"responsible_user_last_name":     "responsible_user.last_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"responsible_user_last_name_ar":  "responsible_user.last_name_ar",
		"responsible_unit_title_en":      "organization_units.title_en",
		"responsible_unit_title_ar":      "organization_units.title_ar",
		"backup_user_first_name":         "backup_user.first_name_en",
		"backup_user_last_name":          "backup_user.last_name_en",
		"backup_user_first_name_ar":      "backup_user.first_name_ar",
		"backup_user_last_name_ar":       "backup_user.last_name_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	queryParams := businessContinuityDTOs.BcpTeamsListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortField,
		Keyword: keyword,
		Filters: filters,
	}

	teams, err := brc.service.GetAll(uint(bcpId), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(brc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, teams)
}

func (brc *BcpTeamController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	team, err := brc.service.GetByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(brc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, team)
}

func (btc *BcpTeamController) GetBcpTeamDetails(c *gin.Context) {
	bcpIdStr := c.Param("bcpId")
	bcpId, err := strconv.ParseUint(bcpIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	bcpTeam, err := btc.service.GetFullDetailsByID(uint(bcpId), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(btc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, bcpTeam)
}

func (brc *BcpTeamController) CreateBcpTeam(c *gin.Context) {
	idStr := c.Param("bcpId")
	bcpId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form businessContinuityDTOs.BcpTeamsFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brc.db, c, err.Error())
		return
	}

	id, err := brc.service.Create(uint(bcpId), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(brc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (brc *BcpTeamController) UpdateBcpTeam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form businessContinuityDTOs.BcpTeamsFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brc.db, c, err.Error())
		return
	}

	err = brc.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(brc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (brc *BcpTeamController) DeleteBcpTeams(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	bcpIdStr := c.Param("bcpId")
	bcpId, err := strconv.ParseUint(bcpIdStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(brc.db, c, err.Error())
		return
	}

	failedIDs, err := brc.service.Delete(uint(bcpId), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(brc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

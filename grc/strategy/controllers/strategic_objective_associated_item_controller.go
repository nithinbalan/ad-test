package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategicObjectiveAssociatedItemController struct {
	service services.StrategicObjectiveAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewStrategicObjectiveAssociatedItemController(service services.StrategicObjectiveAssociatedItemServiceInterface) *StrategicObjectiveAssociatedItemController {
	return &StrategicObjectiveAssociatedItemController{
		service: service,
	}
}

func (soai *StrategicObjectiveAssociatedItemController) GetAllClauses(c *gin.Context) {
	strategyIDStr := c.Param("id")

	strategicObjectiveID, err := strconv.ParseUint(strategyIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := strategyDTOs.StrategicObjectiveAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}
	documents, err := soai.service.GetAllClauses(uint(strategicObjectiveID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, documents)
}
func (soai *StrategicObjectiveAssociatedItemController) UpdateClauses(c *gin.Context) {
	var ControlClauseFormDTO commonDTOs.AssociatedClauseFormDTO

	strategicObjectiveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&ControlClauseFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soai.db, c, err.Error())
		return
	}

	err = soai.service.UpdateClauses(uint(strategicObjectiveID), currentUserID, &ControlClauseFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soai.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Clauses against the stategic objective updated successfully"})
}

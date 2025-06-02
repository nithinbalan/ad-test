package controllers

import (
	logDTOs "jethur/grc/logs/dtos"
	"jethur/grc/logs/repositories"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskSchedularLogController struct {
	repo repositories.TaskSchedularLogRepositoryInterface
	db   *gorm.DB
}

func NewTaskSchedularLogController(db *gorm.DB, repo repositories.TaskSchedularLogRepositoryInterface) *TaskSchedularLogController {
	return &TaskSchedularLogController{
		repo: repo,
		db:   db,
	}
}

func (tslc *TaskSchedularLogController) GetAllTaskSchedularLogs(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(tslc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(tslc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := logDTOs.TaskSchedularLogListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}
	logs, err := tslc.repo.GetAll(&queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tslc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, logs)
}

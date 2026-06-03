package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Status       string         `json:"status"`
	ResponseCode string         `json:"responseCode"`
	Message      string         `json:"message"`
	Data         interface{}    `json:"data"`
	Pagination   PaginationMeta `json:"pagination"`
}

func BuildPagination(page, limit int, totalRows int64) PaginationMeta {
	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(limit))),
	}
}

func GetPaginationParams(c *gin.Context) (page int, limit int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 && limit != -1 {
		limit = 10
	}

	return
}

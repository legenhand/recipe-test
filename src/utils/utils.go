package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/legenhand/recipe-test/src/config"
	"gorm.io/gorm"
	"strconv"
)

type PaginatedResponse struct {
	Page     int         `json:"page"`
	Limit    int         `json:"limit"`
	NextPage string      `json:"next_page,omitempty"`
	Count    int64       `json:"count"`
	Data     interface{} `json:"data"`
}

func GetPaginatedResults(c *gin.Context, query *gorm.DB, results interface{}) PaginatedResponse {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	search := c.DefaultQuery("search", "")

	countQuery := query

	if search != "" {
		countQuery = countQuery.Where("name ILIKE ?", "%"+search+"%")
	}

	var count int64
	countQuery.Model(results).Count(&count)

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	query = query.Preload("Unit")
	query.Find(results)

	var nextPage string
	if (page * limit) < int(count) {
		nextPage = fmt.Sprintf("%s%s?page=%d&limit=%d", config.Cfg.BaseUrl, c.Request.URL.Path, page+1, limit)
	}

	return PaginatedResponse{
		Page:     page,
		Limit:    limit,
		NextPage: nextPage,
		Count:    count,
		Data:     results,
	}
}

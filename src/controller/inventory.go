package controller

import (
	"github.com/legenhand/recipe-test/src/db"
	"github.com/legenhand/recipe-test/src/model"
	"github.com/legenhand/recipe-test/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetInventory(c *gin.Context) {
	var inventories []model.Inventory
	query := db.DB

	res := utils.GetPaginatedResults(c, query, &inventories)

	c.JSON(http.StatusOK, res)
}

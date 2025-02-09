package controller

import (
	"fmt"
	"github.com/legenhand/recipe-test/src/db"
	"github.com/legenhand/recipe-test/src/model"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Ingredient struct {
	InventoryID   uuid.UUID `json:"inventory_id"`
	InventoryName string    `json:"inventory_name"`
	Quantity      float64   `json:"quantity"`
	Unit          string    `json:"unit"`
}

type RecipeRequest struct {
	Name        string       `json:"name"`
	Ingredients []Ingredient `json:"ingredients"`
}

func GenerateSKU() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("SKU-%d", rand.Intn(100000))
}

func CreateRecipe(c *gin.Context) {
	var req RecipeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var totalCOGS float64
	var ingredients []model.RecipeIngredient

	for _, ingredient := range req.Ingredients {
		var inventory model.Inventory

		if ingredient.InventoryID != uuid.Nil {
			if err := db.DB.Preload("Unit").Where("id = ?", ingredient.InventoryID).First(&inventory).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory ID"})
				return
			}
		} else if ingredient.InventoryName != "" {
			if err := db.DB.Preload("Unit").Where("name = ?", ingredient.InventoryName).First(&inventory).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory name"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Either inventory_id or inventory_name must be provided"})
			return
		}

		var unit model.Unit
		if err := db.DB.Where("code = ?", ingredient.Unit).First(&unit).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid unit name"})
			return
		}

		var conversionFactor float64 = 1.0 // Default to 1 if no conversion is found

		if inventory.UnitID != unit.ID {
			var conversion model.UnitConversion
			if err := db.DB.Where("base_id = ? AND target_id = ?", inventory.UnitID, unit.ID).First(&conversion).Error; err == nil {
				conversionFactor = conversion.Value
			} else {
				fmt.Println("No conversion found, using default factor of 1.0")
			}
		}

		// Calculate  (Cost of Goods Sold)
		costPerUnit := inventory.Price / inventory.Qty
		totalCOGS += (ingredient.Quantity * costPerUnit) / conversionFactor

		fmt.Println("Inventory:", inventory.Name)
		fmt.Println("Price per unit:", inventory.Price)
		fmt.Println("Quantity:", inventory.Qty)
		fmt.Println("Cost per unit:", costPerUnit)
		fmt.Println("conversionFactor:", conversionFactor)
		fmt.Println("Calculated COGS:", (ingredient.Quantity/conversionFactor)*costPerUnit)

		recipeIngredient := model.RecipeIngredient{
			ID:          uuid.New(),
			RecipeID:    uuid.New(),
			InventoryID: inventory.ID,
			UnitID:      unit.ID,
			Quantity:    ingredient.Quantity,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		ingredients = append(ingredients, recipeIngredient)
	}

	recipe := model.Recipe{
		ID:        uuid.New(),
		Name:      req.Name,
		COGS:      totalCOGS,
		SKU:       GenerateSKU(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db.DB.Create(&recipe)

	for i := range ingredients {
		ingredients[i].RecipeID = recipe.ID
		db.DB.Create(&ingredients[i])
	}

	var recipeWithIngredients model.Recipe
	if err := db.DB.Preload("Ingredients.Unit").First(&recipeWithIngredients, recipe.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load ingredients"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe created", "recipe": recipeWithIngredients})
}

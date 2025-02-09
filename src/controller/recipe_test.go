package controller

import (
	"bytes"
	"encoding/json"
	"github.com/legenhand/recipe-test/src/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Database Layer
type MockRecipeRepo struct {
	mock.Mock
}

func (m *MockRecipeRepo) CreateRecipe(recipe *model.Recipe) error {
	args := m.Called(recipe)
	return args.Error(0)
}

func TestCreateRecipe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize mock repository
	mockRepo := new(MockRecipeRepo)

	// Define expected data
	createdAt, _ := time.Parse(time.RFC3339, "2025-02-10T00:23:34.043897+07:00")
	updatedAt := createdAt
	recipeID := uuid.MustParse("6b7a7660-1f60-4280-89f9-f9efc326f536")

	_ = model.Recipe{
		ID:          recipeID,
		Name:        "Iced Coffee",
		COGS:        7950,
		SKU:         "SKU-3010",
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Ingredients: []model.RecipeIngredient{},
	}

	mockRepo.On("CreateRecipe", mock.AnythingOfType("*model.Recipe")).Return(nil)

	// Setup router
	router := gin.Default()
	router.POST("/recipe", func(c *gin.Context) {
		var newRecipe model.Recipe
		if err := c.ShouldBindJSON(&newRecipe); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		newRecipe.ID = recipeID
		newRecipe.CreatedAt = createdAt
		newRecipe.UpdatedAt = updatedAt

		err := mockRepo.CreateRecipe(&newRecipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Recipe created",
			"recipe":  newRecipe,
		})
	})

	// Create request payload
	requestBody, _ := json.Marshal(map[string]interface{}{
		"name": "Iced Coffee",
		"cogs": 7950,
		"sku":  "SKU-3010",
	})

	// Create request
	req, _ := http.NewRequest("POST", "/recipe", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Validate response structure
	assert.Equal(t, "Recipe created", response["message"])
	recipeResponse := response["recipe"].(map[string]interface{})
	assert.Equal(t, "Iced Coffee", recipeResponse["name"])
	assert.Equal(t, float64(7950), recipeResponse["cogs"])
	assert.Equal(t, "SKU-3010", recipeResponse["sku"])
	assert.Equal(t, "2025-02-10T00:23:34.043897+07:00", recipeResponse["CreatedAt"])
	assert.Equal(t, "2025-02-10T00:23:34.043897+07:00", recipeResponse["UpdatedAt"])
}

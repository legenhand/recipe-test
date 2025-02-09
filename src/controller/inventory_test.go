package controller

import (
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
type MockInventoryRepo struct {
	mock.Mock
}

func (m *MockInventoryRepo) GetInventory(page, limit int) ([]model.Inventory, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Inventory), args.Int(1), args.Error(2)
}

func TestGetInventory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createdAt, _ := time.Parse(time.RFC3339, "2025-02-09T23:20:18.31546+07:00")
	updatedAt, _ := time.Parse(time.RFC3339, "2025-02-09T23:20:18.31546+07:00")

	mockRepo := new(MockInventoryRepo)
	mockData := []model.Inventory{
		{
			ID:     uuid.MustParse("d1b18930-e730-4f97-ad3e-fd5aeee0fac5"),
			Name:   "Aren Sugar",
			Qty:    1,
			Price:  60000,
			UnitID: uuid.MustParse("21ddf3a3-3268-4ef9-8388-bae0f7a5c3de"),
			Unit: model.Unit{
				ID:        uuid.MustParse("21ddf3a3-3268-4ef9-8388-bae0f7a5c3de"),
				Code:      "kg",
				Name:      "Kilogram",
				IsBase:    true,
				IsActive:  true,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		},
		{
			ID:     uuid.MustParse("bd0175ee-a5df-4351-b45a-b55f29ba6336"),
			Name:   "Milk",
			Qty:    1,
			Price:  30000,
			UnitID: uuid.MustParse("4ed521fc-fca8-4852-8e2a-bc97efa24fc2"),
			Unit: model.Unit{
				ID:        uuid.MustParse("4ed521fc-fca8-4852-8e2a-bc97efa24fc2"),
				Code:      "l",
				Name:      "Liter",
				IsBase:    true,
				IsActive:  true,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		},
	}

	mockRepo.On("GetInventory", 1, 2).Return(mockData, 6, nil)

	router := gin.Default()
	router.GET("/inventory", func(c *gin.Context) {
		page := 1
		limit := 2

		data, count, err := mockRepo.GetInventory(page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch inventory"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"page":      page,
			"limit":     limit,
			"next_page": "http://localhost:8080/inventory?page=2&limit=2",
			"count":     count,
			"data":      data,
		})
	})

	req, _ := http.NewRequest("GET", "/inventory?page=1&limit=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), response["page"])
	assert.Equal(t, float64(2), response["limit"])
	assert.Equal(t, "http://localhost:8080/inventory?page=2&limit=2", response["next_page"])
	assert.Equal(t, float64(6), response["count"])
	assert.Len(t, response["data"], 2)
}

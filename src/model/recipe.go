package model

import (
	"github.com/google/uuid"
	"time"
)

type Recipe struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	COGS        float64   `json:"cogs"`
	SKU         string    `json:"sku" gorm:"type:varchar(50);unique;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Ingredients []RecipeIngredient `json:"ingredients" gorm:"foreignKey:RecipeID"`
}

type RecipeIngredient struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RecipeID    uuid.UUID `json:"recipe_id"`
	InventoryID uuid.UUID `json:"inventory_id"`
	Quantity    float64   `json:"quantity"`
	UnitID      uuid.UUID `json:"unit_id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Unit        Unit      `json:"unit" gorm:"foreignKey:UnitID"`
	Inventory   Inventory `json:"inventory" gorm:"foreignKey:InventoryID"`
}

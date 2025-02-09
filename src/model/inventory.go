package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Inventory struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name   string    `json:"name" gorm:"type:varchar(100);not null"`
	Qty    float64   `json:"qty" gorm:"not null"`
	Price  float64   `json:"price" gorm:"not null"`
	UnitID uuid.UUID `json:"unit_id" gorm:"type:uuid;not null"`
	Unit   Unit      `json:"unit" gorm:"foreignKey:UnitID"`
}

func (i *Inventory) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.New()
	return
}

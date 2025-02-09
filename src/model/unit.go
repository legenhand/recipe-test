package model

import (
	"github.com/google/uuid"
	"time"
)

type Unit struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Code      string    `json:"code" gorm:"column:code;type:varchar(100);uniqueIndex"`
	Name      string    `json:"name" gorm:"column:name;type:varchar(150)"`
	IsBase    bool      `json:"is_base" gorm:"column:is_base;default:false"`
	IsActive  bool      `json:"is_active" gorm:"column:is_active;default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UnitConversion struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	BaseID    uuid.UUID `json:"base_id" gorm:"column:base_id;index"`
	TargetID  uuid.UUID `json:"target_id" gorm:"column:target_id;index"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	BaseUnit   Unit `gorm:"foreignKey:BaseID"`
	TargetUnit Unit `gorm:"foreignKey:TargetID"`
}

package database

import (
	"time"

	"gorm.io/gorm"
)

type (
	Entity interface {
		GetID() string
	}

	BaseEntity struct {
		ID        string         `json:"id" param:"id" query:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
		DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	}
)

var _ Entity = (*BaseEntity)(nil)

func (e BaseEntity) GetID() string {
	return e.ID
}

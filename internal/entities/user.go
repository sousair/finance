package entities

import "github.com/sousair/go-finance/internal/infra/database"

type User struct {
	database.BaseEntity
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

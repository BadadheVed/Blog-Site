package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      *string   `json:"username" gorm:"type:text;default:null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

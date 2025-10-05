package models

import (
	"time"

	"github.com/google/uuid"
)

type Blog struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Title     string    `json:"title" gorm:"type:text;not null;unique"`
	Body      *string   `json:"body" gorm:"type:text;default:null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	AuthorId  uuid.UUID `json:"author_id" gorm:"type:uuid;not null"`
	Author    User      `gorm:"foreignKey:AuthorId;constraint:OnDelete:CASCADE"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"type:text;not null;unique"`
	Description *string   `json:"description" gorm:"type:text;default:null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type ChannelMember struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ChannelID uuid.UUID `json:"channel_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Channel Channel `gorm:"foreignKey:ChannelID;constraint:OnDelete:CASCADE"`
}

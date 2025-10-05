package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeNew    NotificationType = "NEW"
	NotificationTypeEdited NotificationType = "EDITED"
)

type Notification struct {
	ID        uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID        `json:"user_id" gorm:"type:uuid;not null"`
	Type      NotificationType `json:"type" gorm:"type:text;not null"`
	Message   string           `json:"message" gorm:"type:text;not null"`
	Read      bool             `json:"read" gorm:"default:false"`
	CreatedAt time.Time        `json:"created_at" gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

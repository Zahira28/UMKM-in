package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ParentID  *uint64   `gorm:"index" json:"parent_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relasi
	User          *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Product       *Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product,omitempty"`
	ParentComment *Comment  `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"parent_comment,omitempty"`
	Replies       []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

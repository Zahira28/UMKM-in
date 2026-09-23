package model

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_product_user_like" json:"product_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_product_user_like" json:"user_id"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relasi
	User    *User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product,omitempty"`
}

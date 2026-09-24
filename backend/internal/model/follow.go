package model

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_follower_following" json:"follower_id"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_follower_following" json:"following_id"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relasi
	Follower  *User `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"follower,omitempty"`
	Following *User `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE" json:"following,omitempty"`
}

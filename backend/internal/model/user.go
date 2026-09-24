package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Username         string     `gorm:"type:varchar(50);uniqueIndex" json:"username"`
	Email            string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash     string     `gorm:"type:varchar(255)" json:"-"`
	GoogleID         string     `gorm:"type:varchar(255);index" json:"-"`
	FullName         string     `gorm:"type:varchar(100);not null" json:"full_name"`
	PhoneNumber      string     `gorm:"type:varchar(20)" json:"phone_number"`
	Bio              string     `gorm:"type:text" json:"bio"`
	AvatarURL        string     `gorm:"type:text" json:"avatar_url"`
	City             string     `gorm:"type:varchar(100)" json:"city"`
	AddressDetail    string     `gorm:"type:text" json:"address_detail"`
	Latitude         *float64   `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude        *float64   `gorm:"type:decimal(11,8)" json:"longitude"`
	AICredits        int        `gorm:"type:int;default:10;not null" json:"ai_credits"`
	IsVerified       bool       `gorm:"default:false;not null" json:"is_verified"`
	VerificationCode string     `gorm:"type:varchar(10)" json:"-"`
	CodeExpiresAt    *time.Time `json:"-"`
	CreatedAt        time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relasi
	Products  []Product `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"products,omitempty"`
	Comments  []Comment `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
	Likes     []Like    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"likes,omitempty"`
	Followers []Follow  `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE" json:"followers,omitempty"`
	Following []Follow  `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"following,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

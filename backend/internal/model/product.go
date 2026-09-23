package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CategoryID     uint      `gorm:"not null;index" json:"category_id"`
	Name           string    `gorm:"type:varchar(200);not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`
	Price          float64   `gorm:"type:decimal(12,2);not null" json:"price"`
	Unit           string    `gorm:"type:varchar(50);not null" json:"unit"`
	RawImageURL    string    `gorm:"type:text" json:"raw_image_url"`
	StudioImageURL string    `gorm:"type:text" json:"studio_image_url"`
	IsAvailable    bool      `gorm:"default:true;not null" json:"is_available"`
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relasi
	User     *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Comments []Comment `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
	Likes    []Like    `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"likes,omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

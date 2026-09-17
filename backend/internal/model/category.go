package model

type Category struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Slug     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

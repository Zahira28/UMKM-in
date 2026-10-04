package database

import (
	"log"

	"backend/internal/model"

	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) error {
	categories := []model.Category{
		{Name: "Makanan", Slug: "makanan"},
		{Name: "Minuman", Slug: "minuman"},
		{Name: "Pakaian", Slug: "pakaian"},
	}

	for _, cat := range categories {
		var existing model.Category
		err := db.Where("slug = ?", cat.Slug).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&cat).Error; err != nil {
					return err
				}
				log.Printf("Seeded category: %s (%s)\n", cat.Name, cat.Slug)
			} else {
				return err
			}
		}
	}

	return nil
}

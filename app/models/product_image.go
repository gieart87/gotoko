package models

import (
	"time"
)

type ProductImage struct {
	ID         string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	Product    Product
	ProductID  string `gorm:"size:36;index"`
	Path       string `gorm:"type:text"`
	ExtraLarge string `gorm:"type:text"`
	Large      string `gorm:"type:text"`
	Medium     string `gorm:"type:text"`
	Small      string `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

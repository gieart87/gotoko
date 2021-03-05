package models

import "time"

type OrderCustomer struct {
	ID         string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	User       User
	UserID     string `gorm:"size:36;index"`
	Order      Order
	OrderID    string `gorm:"size:36;index"`
	FirstName  string `gorm:"size:100;not null"`
	LastName   string `gorm:"size:100;not null"`
	CityID     string `gorm:"size:100;"`
	ProvinceID string `gorm:"size:100;"`
	Address1   string `gorm:"size:100;"`
	Address2   string `gorm:"size:100;"`
	Phone      string `gorm:"size:50;"`
	Email      string `gorm:"size:100;"`
	PostCode   string `gorm:"size:100;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

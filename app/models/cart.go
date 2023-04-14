package models

import "github.com/shopspring/decimal"

type Cart struct {
	ID              string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	CartItems       []CartItem
	BaseTotalPrice  decimal.Decimal `gorm:"type:decimal(16,2)"`
	TaxAmount       decimal.Decimal `gorm:"type:decimal(16,2)"`
	TaxPercent      decimal.Decimal `gorm:"type:decimal(10,2)"`
	DiscountAmount  decimal.Decimal `gorm:"type:decimal(16,2)"`
	DiscountPercent decimal.Decimal `gorm:"type:decimal(10,2)"`
	GrandTotal      decimal.Decimal `gorm:"type:decimal(16,2)"`
}

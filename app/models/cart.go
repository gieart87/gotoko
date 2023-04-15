package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

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

func (c *Cart) GetCart(db *gorm.DB, cartID string) (*Cart, error) {
	var err error
	var cart Cart

	err = db.Debug().Model(Cart{}).Where("id = ?", cartID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (c *Cart) CreateCart(db *gorm.DB, cartID string) (*Cart, error) {
	cart := &Cart{
		ID:              cartID,
		BaseTotalPrice:  decimal.NewFromInt(0),
		TaxAmount:       decimal.NewFromInt(0),
		TaxPercent:      decimal.NewFromInt(11),
		DiscountAmount:  decimal.NewFromInt(0),
		DiscountPercent: decimal.NewFromInt(0),
		GrandTotal:      decimal.NewFromInt(0),
	}

	err := db.Debug().Create(&cart).Error
	if err != nil {
		return nil, err
	}

	return cart, nil
}

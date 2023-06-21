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
	TotalWeight     int             `gorm:"-"`
}

func (c *Cart) GetCart(db *gorm.DB, cartID string) (*Cart, error) {
	var err error
	var cart Cart

	err = db.Debug().Preload("CartItems").Model(Cart{}).Where("id = ?", cartID).First(&cart).Error
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

func (c *Cart) CalculateCart(db *gorm.DB, cartID string) (*Cart, error) {
	cartBaseTotalPrice := 0.0
	cartTaxAmount := 0.0
	cartDiscountAmount := 0.0
	cartGrandTotal := 0.0

	for _, item := range c.CartItems {
		itemBaseTotal, _ := item.BaseTotal.Float64()
		itemTaxAmount, _ := item.TaxAmount.Float64()
		itemSubTotalTaxAmount := itemTaxAmount * float64(item.Qty)
		itemDiscountAmount, _ := item.DiscountAmount.Float64()
		itemSubTotalDiscountAmount := itemDiscountAmount * float64(item.Qty)
		itemSubTotal, _ := item.SubTotal.Float64()

		cartBaseTotalPrice += itemBaseTotal
		cartTaxAmount += itemSubTotalTaxAmount
		cartDiscountAmount += itemSubTotalDiscountAmount
		cartGrandTotal += itemSubTotal
	}

	var updateCart, cart Cart

	updateCart.BaseTotalPrice = decimal.NewFromFloat(cartBaseTotalPrice)
	updateCart.TaxAmount = decimal.NewFromFloat(cartTaxAmount)
	updateCart.DiscountAmount = decimal.NewFromFloat(cartDiscountAmount)
	updateCart.GrandTotal = decimal.NewFromFloat(cartGrandTotal)

	err := db.Debug().First(&cart, "id = ?", c.ID).Updates(updateCart).Error
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (c *Cart) AddItem(db *gorm.DB, item CartItem) (*CartItem, error) {
	var existItem, updateItem CartItem
	var product Product

	err := db.Debug().Model(Product{}).Where("id = ?", item.ProductID).First(&product).Error
	if err != nil {
		return nil, err
	}

	basePrice, _ := product.Price.Float64()
	taxAmount := GetTaxAmount(basePrice)
	discountAmount := 0.0

	err = db.Debug().Model(CartItem{}).
		Where("cart_id = ?", c.ID).
		Where("product_id = ?", product.ID).
		First(&existItem).Error

	if err != nil {
		subTotal := float64(item.Qty) * (basePrice + taxAmount - discountAmount)

		item.CartID = c.ID
		item.BasePrice = product.Price
		item.BaseTotal = decimal.NewFromFloat(basePrice * float64(item.Qty))
		item.TaxPercent = decimal.NewFromFloat(GetTaxPercent())
		item.TaxAmount = decimal.NewFromFloat(taxAmount)
		item.DiscountPercent = decimal.NewFromFloat(0)
		item.DiscountAmount = decimal.NewFromFloat(discountAmount)
		item.SubTotal = decimal.NewFromFloat(subTotal)

		err = db.Debug().Create(&item).Error
		if err != nil {
			return nil, err
		}

		return &item, nil
	}

	updateItem.Qty = existItem.Qty + item.Qty
	updateItem.BaseTotal = decimal.NewFromFloat(basePrice * float64(updateItem.Qty))

	subTotal := float64(updateItem.Qty) * (basePrice + taxAmount - discountAmount)
	updateItem.SubTotal = decimal.NewFromFloat(subTotal)

	err = db.Debug().First(&existItem, "id = ?", existItem.ID).Updates(updateItem).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *Cart) GetItems(db *gorm.DB, cartID string) ([]CartItem, error) {
	var items []CartItem

	err := db.Debug().Preload("Product").Model(&CartItem{}).
		Where("cart_id = ?", cartID).
		Order("created_at desc").
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (c *Cart) UpdateItemQty(db *gorm.DB, itemID string, qty int) (*CartItem, error) {
	var existItem, updateItem CartItem

	err := db.Debug().Model(CartItem{}).
		Where("id = ?", itemID).
		First(&existItem).Error
	if err != nil {
		return nil, err
	}

	var product Product

	err = db.Debug().Model(Product{}).Where("id = ?", existItem.ProductID).First(&product).Error
	if err != nil {
		return nil, err
	}

	basePrice, _ := product.Price.Float64()
	taxAmount := GetTaxAmount(basePrice)
	discountAmount := 0.0

	updateItem.Qty = qty
	updateItem.BaseTotal = decimal.NewFromFloat(basePrice * float64(updateItem.Qty))

	subTotal := float64(updateItem.Qty) * (basePrice + taxAmount - discountAmount)
	updateItem.SubTotal = decimal.NewFromFloat(subTotal)

	err = db.Debug().First(&existItem, "id = ?", existItem.ID).Updates(updateItem).Error
	if err != nil {
		return nil, err
	}

	return &existItem, nil
}

func (c *Cart) RemoveItemByID(db *gorm.DB, itemID string) error {
	var err error
	var item CartItem

	err = db.Debug().Model(&CartItem{}).Where("id = ?", itemID).First(&item).Error
	if err != nil {
		return err
	}

	err = db.Debug().Delete(&item).Error
	if err != nil {
		return err
	}

	return nil
}

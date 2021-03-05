package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Payment struct {
	ID          string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	Order       Order
	OrderID     string          `gorm:"size:36;index"`
	Number      string          `gorm:"size:100;index"`
	Amount      decimal.Decimal `gorm:"type:decimal(16,2)"`
	Method      string          `gorm:"size:100"`
	Status      string          `gorm:"size:100"`
	Token       string          `gorm:"size:100;index"`
	Payload     string          `gorm:"type:text"`
	PaymentType string          `gorm:"size:100"`
	VaNumber    string          `gorm:"size:100"`
	BillCode    string          `gorm:"size:100"`
	BillKey     string          `gorm:"size:100"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

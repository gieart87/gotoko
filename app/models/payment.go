package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Payment struct {
	ID                string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	Order             Order
	OrderID           string           `gorm:"size:36;index"`
	Number            string           `gorm:"size:100;index"`
	Amount            decimal.Decimal  `gorm:"type:decimal(16,2)"`
	TransactionID     string           `gorm:"size:100;index"`
	TransactionStatus string           `gorm:"size:100;index"`
	Payload           *json.RawMessage `gorm:"type:json;not null;default:'{}'"`
	PaymentType       string           `gorm:"size:100"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

// MidtransNotification is a struct to handle midtrans notification
type MidtransNotification struct {
	TransactionTime        string          `json:"transaction_time"`
	TransactionStatus      string          `json:"transaction_status"`
	TransactionID          string          `json:"transaction_id"`
	StatusMessage          string          `json:"status_message"`
	StatusCode             string          `json:"status_code"`
	SignatureKey           string          `json:"signature_key"`
	PaymentType            string          `json:"payment_type"`
	OrderID                string          `json:"order_id"`
	MerchantID             string          `json:"merchant_id"`
	MaskedCard             string          `json:"masked_card"`
	GrossAmount            string          `json:"gross_amount"`
	SettlementTime         string          `json:"settlement_time"`
	Issuer                 string          `json:"issuer"`
	FraudStatus            string          `json:"fraud_status"`
	Eci                    string          `json:"eci"`
	Currency               string          `json:"currency"`
	ChannelResponseMessage string          `json:"channel_response_message"`
	ChannelResponseCode    string          `json:"channel_response_code"`
	CardType               string          `json:"card_type"`
	Bank                   string          `json:"bank"`
	ApprovalCode           string          `json:"approval_code"`
	BillKey                string          `json:"bill_key"`
	BillerCode             string          `json:"biller_code"`
	Store                  string          `json:"store"`
	VaNumbers              []VaNumber      `json:"va_numbers"`
	PaymentAmounts         []PaymentAmount `json:"payment_amounts"`
	PermataVaNumber        string          `json:"permata_va_number"`
}

// VaNumber is a struct to handle va number
type VaNumber struct {
	Bank     string `json:"bank"`
	VaNumber string `json:"va_number"`
}

// PaymentAmount is a struct to handle payment amount
type PaymentAmount struct {
	PaidAt string `json:"paid_at"`
	Amount string `json:"amount"`
}

func (p *Payment) BeforeCreate(db *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	p.Number = generatePaymentNumber(db)

	return nil
}

func generatePaymentNumber(db *gorm.DB) string {
	now := time.Now()
	month := now.Month()
	year := strconv.Itoa(now.Year())

	dateCode := "/PAYMENT/" + intToRoman(int(month)) + "/" + year

	var latestPayment Payment

	err := db.Debug().Order("created_at DESC").Find(&latestPayment).Error

	latestNumber, _ := strconv.Atoi(strings.Split(latestPayment.Number, "/")[0])
	if err != nil {
		latestNumber = 1
	}

	number := latestNumber + 1

	invoiceNumber := strconv.Itoa(number) + dateCode

	return invoiceNumber
}

func (p *Payment) CreatePayment(db *gorm.DB, payment *Payment) (*Payment, error) {
	result := db.Debug().Create(payment)
	if result.Error != nil {
		return nil, result.Error
	}

	return payment, nil
}

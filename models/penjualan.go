package models

import (
	"time"

	"gorm.io/gorm"
)

type Penjualan struct {
	ID uint `gorm:"primaryKey" json:"id"`

	KodeInvoice string         `json:"kode_invoice"`
	NamaPembeli string         `json:"nama_pembeli"`
	Subtotal    float64        `json:"subtotal"`
	KodeDiskon  *string        `json:"kode_diskon"`
	Diskon      float64        `json:"diskon"`
	Total       float64        `json:"total"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	CreatedBy   string         `json:"created_by"`
}

func (Penjualan) TableName() string {
	return "penjualans"
}

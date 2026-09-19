package models

import (
	"time"

	"gorm.io/gorm"
)

type ItemPenjualan struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	IDPenjualan uint           `gorm:"not null" json:"id_penjualan"`
	IDBarang    uint           `gorm:"not null" json:"id_barang"`
	Jumlah      int            `gorm:"not null" json:"jumlah"`
	Subtotal    float64        `gorm:"not null" json:"subtotal"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (ItemPenjualan) TableName() string {
	return "item_penjualan"
}

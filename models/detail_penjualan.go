package models

import "time"

// DetailPenjualan adalah model untuk menyimpan detail penjualan
type DetailPenjualan struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PenjualanID uint      `gorm:"not null" json:"penjualan_id"`
	BarangID    uint      `gorm:"not null" json:"barang_id"`
	Jumlah      int       `gorm:"not null" json:"jumlah"`
	Harga       float64   `gorm:"not null" json:"harga"`
	Subtotal    float64   `gorm:"not null" json:"subtotal"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName mengembalikan nama tabel untuk model DetailPenjualan
func (DetailPenjualan) TableName() string {
	return "detail_penjualans"
}

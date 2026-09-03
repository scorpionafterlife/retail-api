package models

import "time"

// Penjualan adalah model untuk menyimpan data penjualan
type Penjualan struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Tanggal   time.Time `gorm:"not null" json:"tanggal"`
	Total     float64   `gorm:"not null;default:0" json:"total"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName mengembalikan nama tabel untuk model penjualan
func (Penjualan) TableName() string {
	return "penjualans"
}

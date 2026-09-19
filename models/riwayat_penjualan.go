package models

import "time"
type RiwayatStok struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	BarangID   uint      `gorm:"not null" json:"barang_id"`
	Jenis      string    `gorm:"type:enum('masuk','keluar');not null" json:"jenis"`
	Jumlah     int       `gorm:"not null" json:"jumlah"`
	Keterangan string    `json:"keterangan"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (RiwayatStok) TableName() string {
	return "riwayat_stoks"
}

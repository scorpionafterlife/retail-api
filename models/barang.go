package models

import (
	"time"

	"gorm.io/gorm"
)

type Barang struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	KodeBarang string         `json:"kode_barang"`
	Nama       string         `json:"nama"`
	NamaBarang string         `json:"nama_barang"`
	HargaPokok float64        `json:"harga_pokok"`
	HargaJual  float64        `json:"harga_jual"`
	Harga      float64        `json:"harga"`
	TipeBarang string         `json:"tipe_barang"`
	Stok       uint           `json:"stok"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	CreatedBy  string         `json:"created_by"`
}

func (Barang) TableName() string {
	return "barangs"
}

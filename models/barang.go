package models

import "time"

//tabel barang untuk menyimpan data barang
type Barang struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	NamaBarang string    `gorm:"column:nama_barang;not null" json:"nama_barang"`
	Harga      float64   `gorm:"not null" json:"harga"`
	Stok       int       `gorm:"not null;default:0" json:"stok"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
// Tablename nama table di database mysql
func (Barang) TableName() string {
	return "barangs"
}

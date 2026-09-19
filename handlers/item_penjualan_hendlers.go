package handlers

import (
	"net/http"
	"strings"

	"retail-api/config"
	"retail-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateItemPenjualan(c *gin.Context) {

	var input struct {
		IDPenjualan uint `json:"id_penjualan" binding:"required"`
		IDBarang    uint `json:"id_barang" binding:"required"`
		Jumlah      int  `json:"jumlah" binding:"required"`
	}

	// BACA JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data yang dikirim tidak valid",
		})
		return
	}

	// VALIDASI JUMLAH
	if input.Jumlah <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Jumlah barang harus lebih dari 0",
		})
		return
	}

	var item models.ItemPenjualan
	var barang models.Barang
	var penjualan models.Penjualan

	// TRANSACTION
	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// CARI PENJUALAN
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&penjualan, input.IDPenjualan).Error; err != nil {

			return gorm.ErrRecordNotFound
		}

		// CARI BARANG
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&barang, input.IDBarang).Error; err != nil {

			return gorm.ErrRecordNotFound
		}

		// CEK STOK
		if barang.Stok < uint(input.Jumlah) {
			return &StokTidakCukupError{
				Tersedia: barang.Stok,
				Diminta:  uint(input.Jumlah),
			}
		}

		// HITUNG SUBTOTAL ITEM
		subtotalItem := barang.HargaJual * float64(input.Jumlah)

		// BUAT ITEM PENJUALAN
		item = models.ItemPenjualan{
			IDPenjualan: input.IDPenjualan,
			IDBarang:    input.IDBarang,
			Jumlah:      input.Jumlah,
			Subtotal:    subtotalItem,
		}

		if err := tx.Create(&item).Error; err != nil {
			return err
		}

		// KURANGI STOK BARANG
		barang.Stok -= uint(input.Jumlah)

		if err := tx.Save(&barang).Error; err != nil {
			return err
		}

		// SIMPAN RIWAYAT STOK
		riwayat := models.RiwayatStok{
			BarangID:   barang.ID,
			Jenis:      "keluar",
			Jumlah:     input.Jumlah,
			Keterangan: "Penjualan",
		}

		if err := tx.Create(&riwayat).Error; err != nil {
			return err
		}

		// HITUNG ULANG SUBTOTAL PENJUALAN
		var subtotalPenjualan float64

		if err := tx.
			Table("item_penjualan").
			Where("id_penjualan = ?", input.IDPenjualan).
			Where("deleted_at IS NULL").
			Select("COALESCE(SUM(subtotal), 0)").
			Scan(&subtotalPenjualan).Error; err != nil {

			return err
		}

		// Simpan subtotal terbaru
		penjualan.Subtotal = subtotalPenjualan

		// HITUNG ULANG DISKON
		var nilaiDiskon float64 = 0

		// Kalau transaksi sudah memiliki kode diskon
		if penjualan.KodeDiskon != nil && *penjualan.KodeDiskon != "" {

			var diskonData struct {
				KodeDiskon string
				Amount     float64
				Type       string
			}

			err := tx.
				Table("kode_diskon").
				Where("kode_diskon = ?", *penjualan.KodeDiskon).
				Where("deleted_at IS NULL").
				First(&diskonData).Error

			if err != nil {
				return err
			}

			// Ubah type menjadi huruf kecil
			tipeDiskon := strings.ToLower(
				strings.TrimSpace(diskonData.Type),
			)

			// DISKON PERSEN
			switch tipeDiskon {

			case "percent", "percentage", "persen":

				nilaiDiskon =
					subtotalPenjualan * diskonData.Amount / 100

			// DISKON NOMINAL
			default:

				nilaiDiskon = diskonData.Amount
			}

			// Diskon tidak boleh negatif
			if nilaiDiskon < 0 {
				nilaiDiskon = 0
			}

			// Diskon tidak boleh lebih besar dari subtotal
			if nilaiDiskon > subtotalPenjualan {
				nilaiDiskon = subtotalPenjualan
			}
		}

		// Hitung total
		total := subtotalPenjualan - nilaiDiskon

		// UPDATE PENJUALAN
		penjualan.Diskon = nilaiDiskon
		penjualan.Total = total

		if err := tx.Save(&penjualan).Error; err != nil {
			return err
		}

		return nil
	})

	// Error penjualan / Barang
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Penjualan atau barang tidak ditemukan",
		})
		return
	}

	// error stok
	if stokError, ok := err.(*StokTidakCukupError); ok {

		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "Stok tidak mencukupi",
			"tersedia": stokError.Tersedia,
			"diminta":  stokError.Diminta,
		})

		return
	}

	// error database
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{

		"message": "Item penjualan berhasil dibuat",

		"data": gin.H{
			"id":            item.ID,
			"id_penjualan":  item.IDPenjualan,
			"id_barang":     item.IDBarang,
			"nama_barang":   barang.Nama,
			"jumlah":        item.Jumlah,
			"harga":         barang.HargaJual,
			"subtotal":      item.Subtotal,
			"stok_sekarang": barang.Stok,
			"subtotal_jual": penjualan.Subtotal,
			"diskon":        penjualan.Diskon,
			"total":         penjualan.Total,
		},
	})
}

// error stok tidak cukup
type StokTidakCukupError struct {
	Tersedia uint
	Diminta  uint
}

func (e *StokTidakCukupError) Error() string {
	return "stok tidak mencukupi"
}

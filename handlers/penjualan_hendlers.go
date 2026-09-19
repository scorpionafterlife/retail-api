package handlers

import (
	"net/http"
	"strings"

	"retail-api/config"
	"retail-api/models"

	"github.com/gin-gonic/gin"
)

func CreatePenjualan(c *gin.Context) {

	var input struct {
		KodeInvoice string `json:"kode_invoice" binding:"required"`
		NamaPembeli string `json:"nama_pembeli"`
		CreatedBy   string `json:"created_by"`
	}

	// Ambil data JSON dari Postman
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data yang dikirim tidak valid",
		})
		return
	}

	penjualan := models.Penjualan{
		KodeInvoice: input.KodeInvoice,
		NamaPembeli: input.NamaPembeli,
		Subtotal:    0,
		KodeDiskon:  nil,
		Diskon:      0,
		Total:       0,
		CreatedBy:   input.CreatedBy,
	}

	// Simpan ke database
	if err := config.DB.Create(&penjualan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Penjualan berhasil dibuat",
		"data": gin.H{
			"id":           penjualan.ID,
			"kode_invoice": penjualan.KodeInvoice,
			"nama_pembeli": penjualan.NamaPembeli,
			"subtotal":     penjualan.Subtotal,
			"kode_diskon":  penjualan.KodeDiskon,
			"diskon":       penjualan.Diskon,
			"total":        penjualan.Total,
			"created_by":   penjualan.CreatedBy,
		},
	})
}

// APPLY DISKON
func ApplyDiskon(c *gin.Context) {

	// Data yang diterima dari Postman
	var input struct {
		KodeDiskon string `json:"kode_diskon" binding:"required"`
	}

	// Baca JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "kode_diskon wajib diisi",
		})
		return
	}

	// Ambil ID penjualan dari URL
	id := c.Param("id")

	// Cari data penjualan
	var penjualan models.Penjualan

	if err := config.DB.First(&penjualan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Penjualan tidak ditemukan",
		})
		return
	}

	// Cari kode diskon
	var diskonData struct {
		ID         uint
		KodeDiskon string
		Amount     float64
		Type       string
	}

	err := config.DB.
		Table("kode_diskon").
		Where("kode_diskon = ?", input.KodeDiskon).
		First(&diskonData).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Kode diskon tidak ditemukan",
		})
		return
	}

	// HITUNG SUBTOTAL DARI ITEM PENJUALAN
	var subtotal float64

	err = config.DB.
		Table("item_penjualan").
		Where("id_penjualan = ?", penjualan.ID).
		Select("COALESCE(SUM(subtotal), 0)").
		Scan(&subtotal).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghitung subtotal",
		})
		return
	}

	// Kalau tidak ada item
	if subtotal <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Penjualan belum memiliki item",
		})
		return
	}

	// HITUNG DISKON
	var nilaiDiskon float64

	// Ubah type menjadi lowercase
	tipeDiskon := strings.ToLower(
		strings.TrimSpace(diskonData.Type),
	)

	switch tipeDiskon {

	case "percent", "percentage", "persen":

		// Diskon berdasarkan persentase
		nilaiDiskon = subtotal * diskonData.Amount / 100

	default:

		// Diskon berdasarkan nominal
		nilaiDiskon = diskonData.Amount
	}

	// Diskon tidak boleh negatif
	if nilaiDiskon < 0 {
		nilaiDiskon = 0
	}

	// Diskon tidak boleh lebih besar dari subtotal
	if nilaiDiskon > subtotal {
		nilaiDiskon = subtotal
	}

	// HITUNG TOTAL
	total := subtotal - nilaiDiskon

	// UPDATE DATA PENJUALAN
	penjualan.Subtotal = subtotal
	penjualan.KodeDiskon = &diskonData.KodeDiskon
	penjualan.Diskon = nilaiDiskon
	penjualan.Total = total

	// Simpan perubahan ke database
	if err := config.DB.Save(&penjualan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan diskon",
		})
		return
	}

	// RESPONSE
	c.JSON(http.StatusOK, gin.H{
		"message": "Diskon berhasil diterapkan",

		"data": gin.H{
			"id_penjualan": penjualan.ID,
			"kode_invoice": penjualan.KodeInvoice,
			"nama_pembeli": penjualan.NamaPembeli,
			"subtotal":     penjualan.Subtotal,
			"kode_diskon":  *penjualan.KodeDiskon,
			"diskon":       penjualan.Diskon,
			"total":        penjualan.Total,
		},
	})
}

// GET PENJUALAN BULANAN
func GetPenjualanBulanan(c *gin.Context) {

	// Struktur untuk menampung hasil query
	var hasil []struct {
		Tahun          int     `json:"tahun"`
		Bulan          int     `json:"bulan"`
		TotalPenjualan float64 `json:"total_penjualan"`
	}

	// Mengambil total penjualan berdasarkan tahun dan bulan
	err := config.DB.
		Table("penjualans").
		Select(`
			YEAR(created_at) AS tahun,
			MONTH(created_at) AS bulan,
			SUM(total) AS total_penjualan
		`).
		Where("deleted_at IS NULL").
		Group("YEAR(created_at), MONTH(created_at)").
		Order("tahun ASC, bulan ASC").
		Scan(&hasil).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data penjualan bulanan",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Data penjualan bulanan berhasil diambil",
		"data":    hasil,
	})
}

// GET DETAIL PENJUALAN
func GetDetailPenjualan(c *gin.Context) {

	// Mengambil ID penjualan dari URL
	id := c.Param("id")

	// CARI DATA PENJUALAN
	var penjualan models.Penjualan

	if err := config.DB.First(&penjualan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Penjualan tidak ditemukan",
		})
		return
	}

	// STRUKTUR DATA ITEM PENJUALAN
	var items []struct {
		IDBarang   uint    `json:"id_barang"`
		NamaBarang string  `json:"nama_barang"`
		Jumlah     int     `json:"jumlah"`
		Harga      float64 `json:"harga"`
		Subtotal   float64 `json:"subtotal"`
	}

	// AMBIL ITEM DAN DATA BARANG
	err := config.DB.
		Table("item_penjualan AS ip").
		Select(`
			ip.id_barang,
			b.nama AS nama_barang,
			ip.jumlah,
			b.harga_jual AS harga,
			ip.subtotal
		`).
		Joins("JOIN barangs AS b ON b.id = ip.id_barang").
		Where("ip.id_penjualan = ?", penjualan.ID).
		Where("ip.deleted_at IS NULL").
		Where("b.deleted_at IS NULL").
		Scan(&items).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil item penjualan",
		})
		return
	}

	// RESPONSE
	c.JSON(http.StatusOK, gin.H{
		"message": "Detail penjualan berhasil diambil",

		"data": gin.H{
			"id_penjualan": penjualan.ID,
			"kode_invoice": penjualan.KodeInvoice,
			"nama_pembeli": penjualan.NamaPembeli,
			"subtotal":     penjualan.Subtotal,
			"kode_diskon":  penjualan.KodeDiskon,
			"diskon":       penjualan.Diskon,
			"total":        penjualan.Total,
			"created_at":   penjualan.CreatedAt,

			"items": items,
		},
	})
}

func GetPenjualan(c *gin.Context) {

	var penjualans []models.Penjualan

	// Ambil semua penjualan yang belum dihapus
	err := config.DB.
		Where("deleted_at IS NULL").
		Order("id ASC").
		Find(&penjualans).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data penjualan",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "Data penjualan berhasil diambil",
		"data":    penjualans,
	})
}

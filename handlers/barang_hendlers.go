package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"retail-api/config"
	"retail-api/models"
)

// GET /barang
func GetBarang(c *gin.Context) {
	var barangs []models.Barang

	result := config.DB.
		Where("deleted_at IS NULL").
		Order("id ASC").
		Find(&barangs)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": barangs,
	})
}

// GET /barang/:id
func GetBarangByID(c *gin.Context) {
	id := c.Param("id")

	var barang models.Barang

	if err := config.DB.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Barang tidak ditemukan",
		})
		return
	}

	var riwayats []models.RiwayatStok

	if err := config.DB.
		Where("barang_id = ?", barang.ID).
		Order("id DESC").
		Find(&riwayats).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":         barang,
		"histori_stok": riwayats,
	})
}

// POST /barang
func CreateBarang(c *gin.Context) {
	var input struct {
		KodeBarang string  `json:"kode_barang"`
		Nama       string  `json:"nama"`
		NamaBarang string  `json:"nama_barang"`
		HargaPokok float64 `json:"harga_pokok"`
		HargaJual  float64 `json:"harga_jual"`
		TipeBarang string  `json:"tipe_barang"`
		Stok       uint    `json:"stok"`
		CreatedBy  string  `json:"created_by"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid",
		})
		return
	}

	barang := models.Barang{
		KodeBarang: input.KodeBarang,
		Nama:       input.Nama,
		NamaBarang: input.NamaBarang,
		HargaPokok: input.HargaPokok,
		HargaJual:  input.HargaJual,
		Harga:      input.HargaJual,
		TipeBarang: input.TipeBarang,
		Stok:       input.Stok,
		CreatedBy:  input.CreatedBy,
	}

	result := config.DB.Create(&barang)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Barang berhasil ditambahkan",
		"data":    barang,
	})
}

// PUT /barang/:id
func UpdateBarang(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid",
		})
		return
	}

	var barang models.Barang

	result := config.DB.First(&barang, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Barang tidak ditemukan",
		})
		return
	}

	var input struct {
		KodeBarang *string  `json:"kode_barang"`
		Nama       *string  `json:"nama"`
		NamaBarang *string  `json:"nama_barang"`
		HargaPokok *float64 `json:"harga_pokok"`
		HargaJual  *float64 `json:"harga_jual"`
		Harga      *float64 `json:"harga"`
		TipeBarang *string  `json:"tipe_barang"`
		Stok       *uint    `json:"stok"`
		CreatedBy  *string  `json:"created_by"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format JSON tidak valid",
		})
		return
	}

	if input.KodeBarang != nil {
		barang.KodeBarang = *input.KodeBarang
	}

	if input.Nama != nil {
		barang.Nama = *input.Nama
	}

	if input.NamaBarang != nil {
		barang.NamaBarang = *input.NamaBarang
	}

	if input.HargaPokok != nil {
		barang.HargaPokok = *input.HargaPokok
	}

	if input.HargaJual != nil {
		barang.HargaJual = *input.HargaJual

		// Harga mengikuti harga jual
		barang.Harga = *input.HargaJual
	}

	if input.Harga != nil {
		barang.Harga = *input.Harga
	}

	if input.TipeBarang != nil {
		barang.TipeBarang = *input.TipeBarang
	}

	if input.Stok != nil {
		barang.Stok = *input.Stok
	}

	if input.CreatedBy != nil {
		barang.CreatedBy = *input.CreatedBy
	}

	if err := config.DB.Save(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan perubahan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Barang berhasil diubah",
		"data":    barang,
	})
}

// DELETE /barang/:id
func DeleteBarang(c *gin.Context) {
	id := c.Param("id")

	var barang models.Barang

	if err := config.DB.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Barang tidak ditemukan",
		})
		return
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// Soft delete barang
		if err := tx.Delete(&barang).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Barang berhasil dihapus",
		"barang_id": barang.ID,
		"nama":      barang.Nama,
	})
}

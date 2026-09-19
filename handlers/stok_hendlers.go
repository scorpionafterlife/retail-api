package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"retail-api/config"
	"retail-api/models"
)

func StokMasuk(c *gin.Context) {
	var input struct {
		BarangID   uint   `json:"barang_id" binding:"required"`
		Jumlah     int    `json:"jumlah" binding:"required"`
		Keterangan string `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid",
		})
		return
	}

	if input.Jumlah <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Jumlah stok harus lebih dari 0",
		})
		return
	}

	var barang models.Barang

	if result := config.DB.First(&barang, input.BarangID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Barang tidak ditemukan",
		})
		return
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// Tambahkan stok
		barang.Stok += uint(input.Jumlah)

		if err := tx.Save(&barang).Error; err != nil {
			return err
		}

		// Simpan riwayat stok
		riwayat := models.RiwayatStok{
			BarangID:   input.BarangID,
			Jenis:      "masuk",
			Jumlah:     input.Jumlah,
			Keterangan: input.Keterangan,
		}

		if err := tx.Create(&riwayat).Error; err != nil {
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
		"message":       "Stok berhasil ditambahkan",
		"barang_id":     barang.ID,
		"nama":          barang.Nama,
		"jumlah_masuk":  input.Jumlah,
		"stok_sekarang": barang.Stok,
	})
}

// STOK KELUAR
func StokKeluar(c *gin.Context) {
	var input struct {
		BarangID   uint   `json:"barang_id" binding:"required"`
		Jumlah     int    `json:"jumlah" binding:"required"`
		Keterangan string `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid",
		})
		return
	}

	if input.Jumlah <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Jumlah stok harus lebih dari 0",
		})
		return
	}

	var barang models.Barang

	if result := config.DB.First(&barang, input.BarangID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Barang tidak ditemukan",
		})
		return
	}

	// Cek apakah stok cukup
	if barang.Stok < uint(input.Jumlah) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Stok tidak mencukupi",
			"stok_tersedia": barang.Stok,
		})
		return
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// Kurangi stok
		barang.Stok -= uint(input.Jumlah)

		if err := tx.Save(&barang).Error; err != nil {
			return err
		}

		// Simpan riwayat stok
		riwayat := models.RiwayatStok{
			BarangID:   input.BarangID,
			Jenis:      "keluar",
			Jumlah:     input.Jumlah,
			Keterangan: input.Keterangan,
		}

		if err := tx.Create(&riwayat).Error; err != nil {
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
		"message":       "Stok berhasil dikurangi",
		"barang_id":     barang.ID,
		"nama":          barang.Nama,
		"jumlah_keluar": input.Jumlah,
		"stok_sekarang": barang.Stok,
	})
}
func UpdateStok(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Stok       uint   `json:"stok" binding:"required"`
		Keterangan string `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid",
		})
		return
	}

	var barang models.Barang

	if err := config.DB.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Barang tidak ditemukan",
		})
		return
	}

	stokLama := barang.Stok
	stokBaru := input.Stok

	// Hitung perubahan stok
	var jenis string
	var jumlah uint

	if stokBaru > stokLama {
		jenis = "masuk"
		jumlah = stokBaru - stokLama
	} else if stokBaru < stokLama {
		jenis = "keluar"
		jumlah = stokLama - stokBaru
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// Update stok barang
		barang.Stok = stokBaru

		if err := tx.Save(&barang).Error; err != nil {
			return err
		}

		// Simpan histori stok 
		if stokBaru != stokLama {
			riwayat := models.RiwayatStok{
				BarangID:   barang.ID,
				Jenis:      jenis,
				Jumlah:     int(jumlah),
				Keterangan: input.Keterangan,
			}

			if err := tx.Create(&riwayat).Error; err != nil {
				return err
			}
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
		"message":   "Stok berhasil diperbarui",
		"barang_id": barang.ID,
		"nama":      barang.Nama,
		"stok_lama": stokLama,
		"stok_baru": barang.Stok,
	})
}

// GET RIWAYAT STOK
func GetRiwayatStok(c *gin.Context) {
	var riwayats []models.RiwayatStok

	result := config.DB.
		Order("id DESC").
		Find(&riwayats)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": riwayats,
	})
}

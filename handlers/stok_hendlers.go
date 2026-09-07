package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"retail-api/config"
	"retail-api/models"
)

// ini adalah handler untuk menambahkan stok barang dan mengurangi stok barang, serta menampilkan riwayat stok barang

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
	// mencari barang berdasarkan ID yang dikirim dari postman, jika barang tidak ditemukan maka akan mengembalikan error not found
	var barang models.Barang

	if result := config.DB.First(&barang, input.BarangID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Barang tidak ditemukan",
		})
		return
	}
	// menggunakan transaksi untuk menambahkan stok barang dan membuat riwayat stok masuk, jika terjadi error maka akan mengembalikan error internal server
	err := config.DB.Transaction(func(tx *gorm.DB) error {

		barang.Stok += input.Jumlah

		if err := tx.Save(&barang).Error; err != nil {
			return err
		}

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
	//
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// mengembalikan response dengan stasus OK
	c.JSON(http.StatusOK, gin.H{
		"message":       "Stok berhasil ditambahkan",
		"barang_id":     barang.ID,
		"nama_barang":   barang.NamaBarang,
		"jumlah_masuk":  input.Jumlah,
		"stok_sekarang": barang.Stok,
	})
}

// StokKeluar mengurangi stok barang berdasarkan input yang dikirim dari postman, jika barang tidak ditemukan atau stok tidak mencukupi maka akan mengembalikan error not found atau bad request
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
	// mencari barang berdasarkan ID yang dikirim dari postman, jika barang tidak ditemukan maka akan mengembalikan error not found
	var barang models.Barang

	if result := config.DB.First(&barang, input.BarangID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Barang tidak ditemukan",
		})
		return
	}
	// jika stok barang kurang dari jumlah yang dikirim dari postman maka akan mengembalikan error bad request
	if barang.Stok < input.Jumlah {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         "Stok tidak mencukupi",
			"stok_tersedia": barang.Stok,
		})
		return
	}
	// menggunakan transaksi untuk mengurangi stok barang dan membuat riwayat stok keluar, jika terjadi error maka akan mengambalikan error internal server
	err := config.DB.Transaction(func(tx *gorm.DB) error {

		barang.Stok -= input.Jumlah

		if err := tx.Save(&barang).Error; err != nil {
			return err
		}	

		riwayat := models.RiwayatStok{
			BarangID:   input.BarangID,
			Jenis:      "keluar",
			Jumlah:     input.Jumlah,
			Keterangan: input.Keterangan,
		}
		// menyimpan riwayat stok keluar ke database
		if err := tx.Create(&riwayat).Error; err != nil {
			return err
		}
		// nil jika transaksi berhasil
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// mengembalikan response dengan status OK dan data barang yang stoknya telah dikurangi
	c.JSON(http.StatusOK, gin.H{
		"message":       "Stok berhasil dikurangi",
		"barang_id":     barang.ID,
		"nama_barang":   barang.NamaBarang,
		"jumlah_keluar": input.Jumlah,
		"stok_sekarang": barang.Stok,
	})
}

// mengambil data riwayat stok dari database dalam urutan descending berdasarkan ID
func GetRiwayatStok(c *gin.Context) {
	var riwayats []models.RiwayatStok

	result := config.DB.Order("id DESC").Find(&riwayats)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	// mengembalikan response dengan status OK dan data riwayat stok
	c.JSON(http.StatusOK, gin.H{
		"data": riwayats,
	})
}
  
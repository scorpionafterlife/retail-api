package handlers

import (
	"net/http"
	"strconv"

	"retail-api/config"
	"retail-api/models"

	"github.com/gin-gonic/gin"
)

// getbarang mengambil semua data barang dari database dan mengembalikannya dalam format json
func GetBarang(c *gin.Context) {
	var barangs []models.Barang

	result := config.DB.Find(&barangs)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	// mengembalikan response dengan status ok d=====-an data barang dalam format json
	c.JSON(http.StatusOK, gin.H{
		"data": barangs,
	})
}

// getbarangbyid  mengambil data barang berdasarkan id yang dikirim dari postman, jika barang tidak ditemukan maka akan mengembalikan error no found
func GetBarangByID(c *gin.Context) {
	id := c.Param("id")

	var barang models.Barang

	result := config.DB.First(&barang, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Barang tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": barang,
	})
}

// creatbarang membuat data barang baru berdasarkan input yang dikirm dari postman jika data tidak valid maka akan mengembalikan error bad request jika berhasil maka akan mengembalikan data barang yang baru dibuat
func CreateBarang(c *gin.Context) {
	var input struct {
		NamaBarang string  `json:"nama_barang" binding:"required"`
		Harga      float64 `json:"harga" binding:"required"`
		Stok       int     `json:"stok"`
	}
	// jika data yang dikirim dari postman tidak valid maka akan mengembalikan error nad requst
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid",
		})
		return
	}

	barang := models.Barang{
		NamaBarang: input.NamaBarang,
		Harga:      input.Harga,
		Stok:       input.Stok,
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
		NamaBarang *string  `json:"nama_barang"`
		Harga      *float64 `json:"harga"`
		Stok       *int     `json:"stok"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format JSON tidak valid",
		})
		return
	}

	if input.NamaBarang != nil {
		barang.NamaBarang = *input.NamaBarang
	}

	if input.Harga != nil {
		barang.Harga = *input.Harga
	}

	if input.Stok != nil {
		barang.Stok = *input.Stok
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

func DeleteBarang(c *gin.Context) {
	id := c.Param("id")

	var barang models.Barang

	if result := config.DB.First(&barang, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Barang tidak ditemukan",
		})
		return
	}
	config.DB.Delete(&barang)

	c.JSON(http.StatusOK, gin.H{
		"message": "Barang berhasil dihapus",
	})
}

package main

import (
	"log"

	"retail-api/config"
	"retail-api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDatabase()

	r := gin.Default()

	// barang	
	r.GET("/barang", handlers.GetBarang)
	r.GET("/barang/:id", handlers.GetBarangByID)
	r.POST("/barang", handlers.CreateBarang)
	r.PUT("/barang/:id", handlers.UpdateBarang)
	r.PUT("/barang/stok/:id", handlers.UpdateStok)
	r.DELETE("/barang/:id", handlers.DeleteBarang)

	// stok
	r.POST("/stok/masuk", handlers.StokMasuk)
	r.POST("/stok/keluar", handlers.StokKeluar)
	r.GET("/stok/riwayat", handlers.GetRiwayatStok)

	// penjualan
	r.POST("/penjualan", handlers.CreatePenjualan)
	r.GET("/penjualan", handlers.GetPenjualan)
	r.GET("/penjualan/bulanan", handlers.GetPenjualanBulanan)
	r.GET("/penjualan/:id", handlers.GetDetailPenjualan)
	r.POST("/item-penjualan", handlers.CreateItemPenjualan)
	r.PUT("/penjualan/:id/diskon", handlers.ApplyDiskon)

	log.Println("Server berjalan di http://localhost:8080")

	r.Run(":8080")
}

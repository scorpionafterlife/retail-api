package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"retail-api/config"
	"retail-api/handlers"
)

// file utama untuk menjalankan server dan mengatur routing endpoint
func main() {

	config.ConnectDatabase()

	r := gin.Default()

// Endpoint untuk barang berfunsi untu	
	r.GET("/barang", handlers.GetBarang)
	r.GET("/barang/:id", handlers.GetBarangByID)
	r.POST("/barang", handlers.CreateBarang)
	r.PUT("/barang/:id", handlers.UpdateBarang)
	r.DELETE("/barang/:id", handlers.DeleteBarang)

	// Endpoint untuk stok berfungsi untuk mengelola stok barang, termasuk menambahkan stok masuk, mengurangi stok keluar, dan mengambil riwayat stok.

	r.POST("/stok/masuk", handlers.StokMasuk)
	r.POST("/stok/keluar", handlers.StokKeluar)
	r.GET("/stok/riwayat", handlers.GetRiwayatStok)

	// Endpoint untuk penjualan berfungsi untuk mengelola data penjualan, termasuk membuat data penjualan baru, mengambil semua data penjualan, mengambil data penjualan berdasarkan ID, dan mengambil detail penjualan berdasarkan ID penjualan.
	
	log.Println("Server berjalan di http://localhost:8080")
	// menjalankan server pada port 8080 dan siap menerima permintaan dari client
	r.Run(":8080")
}

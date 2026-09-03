package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// menghubungkan ke database mysql menggunakan gorm

func ConnectDatabase() {
	username := "root"
	password := "gaktau"
	host := "127.0.0.1"
	port := "3306"
	database := "retail_db"
// membuat dsn (data source name) untuk menghubungkan ke database mysql
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
	)

	var err error
// membuka koneksi ke database mysql menggunakan gorm
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	log.Println("Database berhasil terhubung!")
}

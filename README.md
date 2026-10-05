# Retail API
# 
## Deskripsi Project

Project ini merupakan aplikasi Web Service Toko Retail yang dibuat menggunakan bahasa pemrograman Golang dan database MySQL. Aplikasi ini bertujuan untuk membantu mengelola data barang, stok barang, dan transaksi penjualan pada sebuah toko retail.

Sistem menyediakan beberapa fitur utama, yaitu manajemen barang, manajemen stok, pencatatan histori stok, dan proses penjualan. Dengan adanya sistem ini, pengelolaan barang dan transaksi dapat dilakukan secara lebih terstruktur dan mudah.

---

# Teknologi yang Digunakan

Project ini menggunakan beberapa teknologi:

- Golang
- Gin
- GORM
- MySQL
- Postmanhttps://github.com/btwedutech/kelas-beta-golang/blob/main/final-project%2Ftopik-1.md
- dbdiagram.io

## Setup Database MySQL

### 1. Buat Database

Buka MySQL, kemudian jalankan:

```sql
CREATE DATABASE retail_db;
USE retail_db;

```
### Fungsi Teknologi

| Teknologi \ Fungsi 
| 
| Golang Bahasa pemrograman backend 
| Gin Framework untuk membuat REST API 
| GORM ORM untuk menghubungkan Golang dengan MySQL 
| MySQL Database penyimpanan data 
| Postman  Digunakan untuk melakukan testing API 
| dbdiagram.io  Digunakan untuk membuat desain database 

---

# Fitur Project

Project Retail API memiliki beberapa fitur utama:

1. Menampilkan semua barang
2. Menampilkan detail barang
3. Menambahkan barang
4. Mengubah data barang
5. Menghapus barang menggunakan soft delete
6. Menambahkan stok barang
7. Mengurangi stok barang
8. Menampilkan histori stok
9. Membuat transaksi penjualan
10. Menambahkan barang ke transaksi
11. Mengurangi stok secara otomatis ketika terjadi penjualan
12. Mencatat histori stok ketika terjadi penjualan
13. Menggunakan kode diskon
14. Menghitung diskon secara otomatis
15. Menghitung total transaksi
16. Menampilkan detail penjualan
17. Menampilkan seluruh penjualan
18. Menampilkan laporan penjualan bulanan

---

# DB Design

Desain database dibuat menggunakan dbdiagram.io.

![DB Design](## Database Design

![Database Design](resources/DB-Design-Retail.png))

Database utama yang digunakan adalah:

## Struktur Database

Project ini menggunakan database `retail_db` yang terdiri dari beberapa tabel:

1. **barangs**
   Menyimpan data barang, seperti kode barang, nama barang, harga, dan stok.

2. **riwayat_stoks**  
   Menyimpan riwayat perubahan stok barang, baik stok masuk maupun stok keluar.

3. **penjualans**  
   Menyimpan data transaksi penjualan, seperti kode invoice, nama pembeli, subtotal, diskon, dan total.

4. **item_penjualan**  
   Menyimpan daftar barang yang terdapat dalam setiap transaksi penjualan, termasuk jumlah dan subtotal.

5. **kode_diskon**  
   Menyimpan kode diskon yang dapat digunakan dalam transaksi penjualan.

## Fitur Utama

Project ini memiliki beberapa fitur utama, yaitu:

1. **Manajemen Barang**
   - Menampilkan semua data barang.
   - Menampilkan detail barang berdasarkan ID.
   - Menambahkan barang baru.
   - Mengubah data barang.
   - Menghapus barang.

2. **Manajemen Stok**
   - Menambahkan stok barang.
   - Mengurangi stok barang.
   - Melihat riwayat stok masuk dan keluar.
   - Memeriksa ketersediaan stok sebelum penjualan.

3. **Manajemen Penjualan**
   - Membuat transaksi penjualan.
   - Menambahkan beberapa barang dalam satu transaksi.
   - Mengurangi stok secara otomatis setelah barang terjual.
   - Melihat daftar dan detail transaksi penjualan.

4. **Manajemen Diskon**
   - Menggunakan kode diskon.
   - Mendukung diskon persentase dan nominal.

5. **Laporan Penjualan**
   - Menampilkan total penjualan berdasarkan bulan

   ## Daftar Endpoint API

API ini menyediakan beberapa endpoint untuk mengelola barang, stok, dan penjualan.

### 1. Endpoint Barang

| Method | Endpoint | Fungsi |
|---|---|---|
| GET | `/barang` | Menampilkan semua barang |
| GET | `/barang/:id` | Menampilkan barang berdasarkan ID |
| POST | `/barang` | Menambahkan barang baru |
| PUT | `/barang/:id` | Mengubah data barang |
| PUT | `/barang/stok/:id` | Mengubah stok barang |
| DELETE | `/barang/:id` | Menghapus barang |

### 2. Endpoint Stok

| Method | Endpoint | Fungsi |
|---|---|---|
| POST | `/stok/masuk` | Menambahkan stok barang |
| POST | `/stok/keluar` | Mengurangi stok barang |
| GET | `/stok/riwayat` | Menampilkan riwayat stok |

### 3. Endpoint Penjualan

| Method | Endpoint | Fungsi |
|---|---|---|
| POST | `/penjualan` | Membuat transaksi penjualan |
| GET | `/penjualan` | Menampilkan semua transaksi |
| GET | `/penjualan/:id` | Menampilkan detail transaksi |
| GET | `/penjualan/bulanan` | Menampilkan laporan penjualan bulanan |
| POST | `/item-penjualan` | Menambahkan barang ke transaksi |
| PUT | `/penjualan/:id/diskon` | Menggunakan kode diskon |

## Contoh Penggunaan API (Postman)

Berikut beberapa contoh penggunaan API menggunakan Postman.

**Base URL:**
`http://localhost:8080`

### 1. Menampilkan Semua Barang

**Method:** GET  
**Endpoint:** `/barang`

**Request:**

Tidak memerlukan Body.

**Response:**

```json
{
  "data": [
    {
      "id": 1,
      "kode_barang": "BRG001",
      "nama_barang": "Beras 5kg",
      "harga": 75000,
      "stok": 10
    }
  ]
}
```

### 2. Menambahkan Barang


**Method:** POST  
**Endpoint:** `/barang`

**Request Body (JSON):**

```json
{
  "kode_barang": "BRG012",
  "nama_barang": "Tepung Terigu 1kg",
  "harga": 13000,
  "stok": 10,
  "created_by": "admin"
}
```

**Response:**

```json
{
  "message": "Barang berhasil ditambahkan",
  "data": {
    "id": 12,
    "kode_barang": "BRG012",
    "nama_barang": "Tepung Terigu 1kg",
    "harga": 13000,
    "stok": 10
  }
}
```

### 3. Menambahkan Stok Barang

**Method:** POST  
**Endpoint:** `/stok/masuk`

**Request Body (JSON):**

```json
{
  "barang_id": 1,
  "jumlah": 10,
  "keterangan": "Penambahan stok"
}
```

**Response:**

```json
{
  "message": "Stok berhasil ditambahkan"
}
```

### 4. Mengurangi Stok Barang

**Method:** POST  
**Endpoint:** `/stok/keluar`

**Request Body (JSON):**

```json
{
  "barang_id": 1,
  "jumlah": 5,
  "keterangan": "Pengurangan stok"
}
```

**Response:**

```json
{
  "message": "Stok berhasil dikurangi"
}
```

### 5. Membuat Transaksi Penjualan

**Method:** POST  
**Endpoint:** `/penjualan`

**Request Body (JSON):**

```json
{
  "kode_invoice": "INV-TEST-001",
  "nama_pembeli": "Pelanggan"
}
```

**Response:**

```json
{
  "message": "Penjualan berhasil dibuat",
  "data": {
    "id": 10,
    "kode_invoice": "INV-TEST-001",
    "nama_pembeli": "Pelanggan",
    "subtotal": 0,
    "diskon": 0,
    "total": 0
  }
}
```

### 6. Menambahkan Barang ke Transaksi

**Method:** POST  
**Endpoint:** `/item-penjualan`

**Request Body (JSON):**

```json
{
  "id_penjualan": 10,
  "id_barang": 1,
  "jumlah": 2
}
```

**Response:**

```json
{
  "message": "Item penjualan berhasil ditambahkan"
}
```

### 7. Menggunakan Kode Diskon

**Method:** PUT  
**Endpoint:** `/penjualan/10/diskon`

**Request Body (JSON):**

```json
{
  "kode_diskon": "HEMAT10"
}
```

**Response:**

```json
{
  "message": "Diskon berhasil diterapkan"
}
```

### 8. Melihat Laporan Penjualan Bulanan

**Method:** GET  
**Endpoint:** `/penjualan/bulanan`

**Request:**

Tidak memerlukan Body.

**Response:**

```json
{
  "data": [
    {
      "tahun": 2026,
      "bulan": 8,
      "total_penjualan": 178000
    },
    {
      "tahun": 2026,
      "bulan": 9,
      "total_penjualan": 626200
    }
  ],
  "message": "Data penjualan bulanan berhasil diambil"
}
```
## Cara Menjalankan Project

### 1. Clone repository

```bash
git clone https://github.com/scorpionafterlife/retail-api.git
cd retail-api

```
### 2. Install dependency
go mod tidy

### 3. Jalankan server
go run .
### jika berhasil server akan berjalan di http://localhost:8080



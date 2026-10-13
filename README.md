# Evermos Mini Project: REST API (Golang)

Service sederhana untuk transaksi penjualan di platform social commerce Evermos.
Stack: **Go, Fiber v2, GORM, MySQL/MariaDB, JWT**. Arsitektur: **Clean Architecture**.

## Struktur Folder

```
cmd/api/main.go                  Entry point & wiring dependency
internal/
  domain/                        Entity, interface repository, error domain (tidak bergantung ke layer lain)
  usecase/                       Logika bisnis (auth, user, toko, alamat, kategori, produk, transaksi, wilayah)
  repository/mysql/              Implementasi repository dengan GORM
  delivery/http/                 Handler Fiber, middleware JWT/admin, router
  pkg/security/                  JWT & bcrypt
  pkg/upload/                    Validasi & penyimpanan file upload
  config/                        Env, koneksi DB, auto-migrate
uploads/                         Hasil upload (dilayani di /uploads)
```

Arah dependency: `delivery -> usecase -> domain <- repository`.

## Menjalankan

1. Prasyarat: Go 1.22+, MySQL/MariaDB.
2. Buat database: `CREATE DATABASE evermos_mini CHARACTER SET utf8mb4;`
3. Salin `.env.example` menjadi `.env`, isi kredensial DB dan `JWT_SECRET`.
4. Unduh dependency lalu jalankan:
   ```
   go mod tidy
   go run ./cmd/api
   ```
   Tabel dibuat otomatis (AutoMigrate) saat aplikasi start.

## Menjadi Admin

Kategori hanya dapat dikelola admin. Ubah langsung di database, lalu **login ulang** (status admin dibawa di JWT):

```sql
UPDATE users SET is_admin = 1 WHERE id = 1;
```

## Autentikasi

Kirim JWT di header `token: <jwt>` (atau `Authorization: Bearer <jwt>`).

## Endpoint

Respons seragam: `{ "status": bool, "message": string, "errors": [], "data": ... }`

| Method | Path | Auth | Keterangan |
|---|---|---|---|
| POST | /auth/register | - | Daftar; toko otomatis dibuat |
| POST | /auth/login | - | Login dengan `no_telp` + `kata_sandi`, mengembalikan token |
| GET, PUT | /user | JWT | Profil sendiri |
| GET, POST | /user/alamat | JWT | List (filter `judul_alamat`) / tambah alamat |
| GET, PUT, DELETE | /user/alamat/:id | JWT | Detail / ubah / hapus alamat milik sendiri |
| GET | /toko | - | List toko (`nama`, `page`, `limit`) |
| GET | /toko/my | JWT | Toko milik sendiri |
| GET | /toko/:id_toko | - | Detail toko |
| PUT | /toko/:id_toko | JWT | Ubah toko milik sendiri (form-data: `nama_toko`, `photo`) |
| GET | /category, /category/:id | - | Baca kategori |
| POST, PUT, DELETE | /category, /category/:id | Admin | Kelola kategori |
| GET | /product | - | List produk |
| GET | /product/:id | - | Detail produk |
| POST | /product | JWT | Tambah produk ke toko sendiri (form-data) |
| PUT, DELETE | /product/:id | JWT | Ubah / hapus produk milik toko sendiri |
| GET | /trx | JWT | Riwayat transaksi sendiri (`search` kode invoice) |
| GET | /trx/:id | JWT | Detail transaksi sendiri |
| POST | /trx | JWT | Buat transaksi |
| GET | /provcity/listprovincies | - | Daftar provinsi (emsifa.com) |
| GET | /provcity/listcities/:prov_id | - | Kota/kabupaten per provinsi |
| GET | /provcity/detailprovince/:prov_id | - | Detail provinsi |
| GET | /provcity/detailcity/:city_id | - | Detail kota |

### Pagination & filter

- Pagination: `?page=1&limit=10` (default 1 dan 10, maksimal limit 100). Respons list berisi `data`, `page`, `limit`, `total`.
- Filter produk: `nama_produk`, `category_id`, `toko_id`, `min_harga`, `max_harga`.

### Contoh body

Register:
```json
{
  "nama": "Devi", "kata_sandi": "password123", "no_telp": "081234567890",
  "tanggal_lahir": "2000-01-31", "jenis_kelamin": "P", "tentang": "", "pekerjaan": "Mahasiswa",
  "email": "devi@example.com", "id_provinsi": "33", "id_kota": "3372"
}
```

Tambah produk (`multipart/form-data`): `nama_produk`, `category_id`, `harga_reseller`, `harga_konsumen`, `stok`, `deskripsi`, `photos` (file, boleh banyak; JPG/PNG maks 5MB).

Transaksi:
```json
{
  "method_bayar": "transfer",
  "alamat_kirim": 1,
  "detail_trx": [ { "product_id": 1, "kuantitas": 2 } ]
}
```

## Aturan yang Diterapkan

- Email dan no telp unik (dicek di usecase + unique index di database).
- Toko dibuat otomatis saat register, dalam satu transaksi database dengan user.
- User hanya bisa mengakses/mengubah data miliknya: profil (diambil dari `user_id` di token), alamat, toko, produk, dan transaksi. Akses ke data milik orang lain dijawab `403`.
- Transaksi bersifat atomik: baris produk dikunci, stok divalidasi lalu dikurangi, snapshot produk disimpan ke tabel `log_produks`, lalu `trxes` dan `detail_trxes` dibuat. Jika salah satu gagal, semuanya di-rollback.
- `detail_trxes` merujuk ke `log_produks` (bukan `produks`), sehingga riwayat transaksi tetap benar walau produk diubah atau dihapus.
- Produk dihapus secara soft delete.

## Catatan

- Desain database di link Google Drive pada soal tidak dapat dibuka saat pengerjaan, jadi skema disusun mengikuti kebutuhan soal dan collection Postman. Jika desain resmi punya nama tabel/kolom berbeda, sesuaikan tag GORM di `internal/domain/entity.go`.

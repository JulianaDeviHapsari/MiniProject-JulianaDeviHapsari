package domain

import (
	"time"

	"gorm.io/gorm"
)

// User: email dan no_telp unik. IsAdmin diubah langsung lewat database.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Nama         string    `gorm:"size:255;not null" json:"nama"`
	KataSandi    string    `gorm:"size:255;not null" json:"-"`
	NoTelp       string    `gorm:"size:20;not null;uniqueIndex" json:"no_telp"`
	Email        string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	TanggalLahir time.Time `gorm:"type:date" json:"tanggal_lahir"`
	JenisKelamin string    `gorm:"size:20" json:"jenis_kelamin"`
	Tentang      string    `gorm:"type:text" json:"tentang"`
	Pekerjaan    string    `gorm:"size:255" json:"pekerjaan"`
	IDProvinsi   string    `gorm:"size:20" json:"id_provinsi"`
	IDKota       string    `gorm:"size:20" json:"id_kota"`
	IsAdmin      bool      `gorm:"default:false" json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Toko struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IDUser    uint      `gorm:"not null;uniqueIndex" json:"id_user"`
	NamaToko  string    `gorm:"size:255;not null" json:"nama_toko"`
	UrlFoto   string    `gorm:"size:500" json:"url_foto"`
	User      *User     `gorm:"foreignKey:IDUser" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Alamat struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	IDUser       uint      `gorm:"not null;index" json:"id_user"`
	JudulAlamat  string    `gorm:"size:255;not null" json:"judul_alamat"`
	NamaPenerima string    `gorm:"size:255;not null" json:"nama_penerima"`
	NoTelp       string    `gorm:"size:20;not null" json:"no_telp"`
	DetailAlamat string    `gorm:"type:text;not null" json:"detail_alamat"`
	User         *User     `gorm:"foreignKey:IDUser" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Kategori struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	NamaKategori string    `gorm:"size:255;not null;uniqueIndex" json:"nama_category"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Produk struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	NamaProduk    string         `gorm:"size:255;not null;index" json:"nama_produk"`
	Slug          string         `gorm:"size:300;not null;uniqueIndex" json:"slug"`
	HargaReseller int            `gorm:"not null" json:"harga_reseller"`
	HargaKonsumen int            `gorm:"not null" json:"harga_konsumen"`
	Stok          int            `gorm:"not null" json:"stok"`
	Deskripsi     string         `gorm:"type:text" json:"deskripsi"`
	IDToko        uint           `gorm:"not null;index" json:"id_toko"`
	IDKategori    uint           `gorm:"not null;index" json:"id_kategori"`
	Toko          *Toko          `gorm:"foreignKey:IDToko" json:"toko,omitempty"`
	Kategori      *Kategori      `gorm:"foreignKey:IDKategori" json:"kategori,omitempty"`
	Foto          []FotoProduk   `gorm:"foreignKey:IDProduk" json:"foto"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type FotoProduk struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IDProduk  uint      `gorm:"not null;index" json:"id_produk"`
	Url       string    `gorm:"size:500;not null" json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

// LogProduk menyimpan snapshot produk pada saat transaksi terjadi.
// Sengaja tanpa foreign key ke produk agar tetap valid walau produk diubah/dihapus.
type LogProduk struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	IDProduk      uint      `gorm:"not null;index" json:"id_produk"`
	NamaProduk    string    `gorm:"size:255;not null" json:"nama_produk"`
	Slug          string    `gorm:"size:300" json:"slug"`
	HargaReseller int       `json:"harga_reseller"`
	HargaKonsumen int       `json:"harga_konsumen"`
	Deskripsi     string    `gorm:"type:text" json:"deskripsi"`
	UrlFoto       string    `gorm:"size:500" json:"url_foto"`
	IDToko        uint      `gorm:"not null" json:"id_toko"`
	IDKategori    uint      `gorm:"not null" json:"id_kategori"`
	CreatedAt     time.Time `json:"created_at"`
}

type Trx struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	IDUser      uint        `gorm:"not null;index" json:"id_user"`
	IDAlamat    uint        `gorm:"not null" json:"id_alamat"`
	HargaTotal  int         `gorm:"not null" json:"harga_total"`
	KodeInvoice string      `gorm:"size:100;not null;uniqueIndex" json:"kode_invoice"`
	MethodBayar string      `gorm:"size:50;not null" json:"method_bayar"`
	Alamat      *Alamat     `gorm:"foreignKey:IDAlamat" json:"alamat_kirim,omitempty"`
	Detail      []DetailTrx `gorm:"foreignKey:IDTrx" json:"detail_trx"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type DetailTrx struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	IDTrx       uint       `gorm:"not null;index" json:"id_trx"`
	IDLogProduk uint       `gorm:"not null" json:"id_log_produk"`
	IDToko      uint       `gorm:"not null" json:"id_toko"`
	Kuantitas   int        `gorm:"not null" json:"kuantitas"`
	HargaTotal  int        `gorm:"not null" json:"harga_total"`
	LogProduk   *LogProduk `gorm:"foreignKey:IDLogProduk" json:"log_produk,omitempty"`
	Toko        *Toko      `gorm:"foreignKey:IDToko" json:"toko,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

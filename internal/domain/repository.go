package domain

// Paginated adalah hasil list dengan metadata pagination.
type Paginated[T any] struct {
	Data  []T   `json:"data"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type Page struct {
	Page  int
	Limit int
}

func (p Page) Offset() int { return (p.Page - 1) * p.Limit }

type ProdukFilter struct {
	Page
	NamaProduk string
	IDKategori uint
	IDToko     uint
	MinHarga   int
	MaxHarga   int
}

type TrxItemInput struct {
	IDProduk  uint
	Kuantitas int
}

type UserRepository interface {
	// CreateWithToko menyimpan user dan toko dalam satu transaksi database.
	CreateWithToko(user *User, toko *Toko) error
	FindByID(id uint) (*User, error)
	FindByNoTelp(noTelp string) (*User, error)
	ExistsByEmail(email string, excludeID uint) (bool, error)
	ExistsByNoTelp(noTelp string, excludeID uint) (bool, error)
	Update(user *User) error
}

type TokoRepository interface {
	FindByID(id uint) (*Toko, error)
	FindByUserID(userID uint) (*Toko, error)
	List(nama string, page Page) ([]Toko, int64, error)
	Update(toko *Toko) error
}

type AlamatRepository interface {
	Create(a *Alamat) error
	FindByID(id uint) (*Alamat, error)
	ListByUser(userID uint, judul string) ([]Alamat, error)
	Update(a *Alamat) error
	Delete(id uint) error
}

type KategoriRepository interface {
	Create(k *Kategori) error
	FindByID(id uint) (*Kategori, error)
	List() ([]Kategori, error)
	Update(k *Kategori) error
	Delete(id uint) error
}

type ProdukRepository interface {
	Create(p *Produk) error
	FindByID(id uint) (*Produk, error)
	List(f ProdukFilter) ([]Produk, int64, error)
	Update(p *Produk) error
	AddFoto(fotos []FotoProduk) error
	Delete(id uint) error
}

type TrxRepository interface {
	// Create membuat transaksi secara atomik: mengunci stok, mengisi log_produk,
	// mengurangi stok, lalu menyimpan trx + detail_trx. Rollback jika ada yang gagal.
	Create(userID, alamatID uint, methodBayar string, items []TrxItemInput) (*Trx, error)
	FindByID(id uint) (*Trx, error)
	ListByUser(userID uint, search string, page Page) ([]Trx, int64, error)
}

package usecase

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"evermos-mini/internal/domain"
)

// ---------- Produk ----------

type ProdukInput struct {
	NamaProduk    string `validate:"required"`
	IDKategori    uint   `validate:"required"`
	HargaReseller int    `validate:"required,gt=0"`
	HargaKonsumen int    `validate:"required,gt=0"`
	Stok          int    `validate:"gte=0"`
	Deskripsi     string
}

type ProdukUsecase struct {
	produks   domain.ProdukRepository
	tokos     domain.TokoRepository
	kategoris domain.KategoriRepository
}

func NewProdukUsecase(p domain.ProdukRepository, t domain.TokoRepository, k domain.KategoriRepository) *ProdukUsecase {
	return &ProdukUsecase{p, t, k}
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.Trim(nonAlnum.ReplaceAllString(strings.ToLower(s), "-"), "-")
	return fmt.Sprintf("%s-%d", s, time.Now().UnixNano()/1e6)
}

func (u *ProdukUsecase) List(f domain.ProdukFilter) (*domain.Paginated[domain.Produk], error) {
	list, total, err := u.produks.List(f)
	if err != nil {
		return nil, err
	}
	return &domain.Paginated[domain.Produk]{Data: list, Page: f.Page.Page, Limit: f.Limit, Total: total}, nil
}

func (u *ProdukUsecase) Get(id uint) (*domain.Produk, error) { return u.produks.FindByID(id) }

// Create: produk otomatis masuk ke toko milik user yang login.
func (u *ProdukUsecase) Create(userID uint, in ProdukInput, fotoURLs []string) (*domain.Produk, error) {
	toko, err := u.tokos.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if _, err := u.kategoris.FindByID(in.IDKategori); err != nil {
		return nil, domain.BadRequest("kategori tidak ditemukan")
	}
	p := &domain.Produk{
		NamaProduk: in.NamaProduk, Slug: slugify(in.NamaProduk),
		HargaReseller: in.HargaReseller, HargaKonsumen: in.HargaKonsumen,
		Stok: in.Stok, Deskripsi: in.Deskripsi, IDToko: toko.ID, IDKategori: in.IDKategori,
	}
	for _, url := range fotoURLs {
		p.Foto = append(p.Foto, domain.FotoProduk{Url: url})
	}
	if err := u.produks.Create(p); err != nil {
		return nil, err
	}
	return u.produks.FindByID(p.ID)
}

func (u *ProdukUsecase) ownProduk(userID, id uint) (*domain.Produk, error) {
	p, err := u.produks.FindByID(id)
	if err != nil {
		return nil, err
	}
	toko, err := u.tokos.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if p.IDToko != toko.ID {
		return nil, domain.Forbidden("tidak dapat mengelola produk milik user lain")
	}
	return p, nil
}

func (u *ProdukUsecase) Update(userID, id uint, in ProdukInput, fotoURLs []string) (*domain.Produk, error) {
	p, err := u.ownProduk(userID, id)
	if err != nil {
		return nil, err
	}
	if _, err := u.kategoris.FindByID(in.IDKategori); err != nil {
		return nil, domain.BadRequest("kategori tidak ditemukan")
	}
	// Slug sengaja tidak diubah agar URL produk stabil.
	p.NamaProduk, p.IDKategori = in.NamaProduk, in.IDKategori
	p.HargaReseller, p.HargaKonsumen = in.HargaReseller, in.HargaKonsumen
	p.Stok, p.Deskripsi = in.Stok, in.Deskripsi
	if err := u.produks.Update(p); err != nil {
		return nil, err
	}
	var fotos []domain.FotoProduk
	for _, url := range fotoURLs {
		fotos = append(fotos, domain.FotoProduk{IDProduk: p.ID, Url: url})
	}
	if err := u.produks.AddFoto(fotos); err != nil {
		return nil, err
	}
	return u.produks.FindByID(p.ID)
}

func (u *ProdukUsecase) Delete(userID, id uint) error {
	if _, err := u.ownProduk(userID, id); err != nil {
		return err
	}
	return u.produks.Delete(id)
}

// ---------- Transaksi ----------

type TrxItemDTO struct {
	IDProduk  uint `json:"product_id" validate:"required"`
	Kuantitas int  `json:"kuantitas" validate:"required,gt=0"`
}

type CreateTrxInput struct {
	MethodBayar string       `json:"method_bayar" validate:"required"`
	AlamatKirim uint         `json:"alamat_kirim" validate:"required"`
	DetailTrx   []TrxItemDTO `json:"detail_trx" validate:"required,min=1,dive"`
}

type TrxUsecase struct {
	trxs    domain.TrxRepository
	alamats domain.AlamatRepository
}

func NewTrxUsecase(t domain.TrxRepository, a domain.AlamatRepository) *TrxUsecase {
	return &TrxUsecase{t, a}
}

func (u *TrxUsecase) Create(userID uint, in CreateTrxInput) (*domain.Trx, error) {
	alamat, err := u.alamats.FindByID(in.AlamatKirim)
	if err != nil {
		return nil, domain.BadRequest("alamat kirim tidak ditemukan")
	}
	if alamat.IDUser != userID {
		return nil, domain.Forbidden("alamat kirim bukan milik anda")
	}

	// Gabungkan item dengan produk yang sama agar pengecekan stok akurat.
	merged := map[uint]int{}
	var order []uint
	for _, it := range in.DetailTrx {
		if _, ok := merged[it.IDProduk]; !ok {
			order = append(order, it.IDProduk)
		}
		merged[it.IDProduk] += it.Kuantitas
	}
	items := make([]domain.TrxItemInput, 0, len(order))
	for _, id := range order {
		items = append(items, domain.TrxItemInput{IDProduk: id, Kuantitas: merged[id]})
	}
	return u.trxs.Create(userID, alamat.ID, in.MethodBayar, items)
}

func (u *TrxUsecase) List(userID uint, search string, page domain.Page) (*domain.Paginated[domain.Trx], error) {
	list, total, err := u.trxs.ListByUser(userID, search, page)
	if err != nil {
		return nil, err
	}
	return &domain.Paginated[domain.Trx]{Data: list, Page: page.Page, Limit: page.Limit, Total: total}, nil
}

func (u *TrxUsecase) Get(userID, id uint) (*domain.Trx, error) {
	t, err := u.trxs.FindByID(id)
	if err != nil {
		return nil, err
	}
	if t.IDUser != userID {
		return nil, domain.Forbidden("tidak dapat mengakses transaksi milik user lain")
	}
	return t, nil
}

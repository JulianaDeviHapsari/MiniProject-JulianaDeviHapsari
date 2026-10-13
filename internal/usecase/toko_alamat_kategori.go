package usecase

import "evermos-mini/internal/domain"

// ---------- Toko ----------

type TokoUsecase struct{ tokos domain.TokoRepository }

func NewTokoUsecase(t domain.TokoRepository) *TokoUsecase { return &TokoUsecase{t} }

func (u *TokoUsecase) List(nama string, page domain.Page) (*domain.Paginated[domain.Toko], error) {
	list, total, err := u.tokos.List(nama, page)
	if err != nil {
		return nil, err
	}
	return &domain.Paginated[domain.Toko]{Data: list, Page: page.Page, Limit: page.Limit, Total: total}, nil
}

func (u *TokoUsecase) GetByID(id uint) (*domain.Toko, error) { return u.tokos.FindByID(id) }

func (u *TokoUsecase) GetMine(userID uint) (*domain.Toko, error) {
	return u.tokos.FindByUserID(userID)
}

// Update: hanya pemilik toko. namaToko kosong = tidak diubah; urlFoto kosong = tidak diubah.
func (u *TokoUsecase) Update(userID, tokoID uint, namaToko, urlFoto string) (*domain.Toko, error) {
	t, err := u.tokos.FindByID(tokoID)
	if err != nil {
		return nil, err
	}
	if t.IDUser != userID {
		return nil, domain.Forbidden("tidak dapat mengelola toko milik user lain")
	}
	if namaToko != "" {
		t.NamaToko = namaToko
	}
	if urlFoto != "" {
		t.UrlFoto = urlFoto
	}
	if err := u.tokos.Update(t); err != nil {
		return nil, err
	}
	return t, nil
}

// ---------- Alamat ----------

type AlamatInput struct {
	JudulAlamat  string `json:"judul_alamat" validate:"required"`
	NamaPenerima string `json:"nama_penerima" validate:"required"`
	NoTelp       string `json:"no_telp" validate:"required,numeric,min=8,max=15"`
	DetailAlamat string `json:"detail_alamat" validate:"required"`
}

type AlamatUsecase struct{ alamats domain.AlamatRepository }

func NewAlamatUsecase(a domain.AlamatRepository) *AlamatUsecase { return &AlamatUsecase{a} }

func (u *AlamatUsecase) own(userID, id uint) (*domain.Alamat, error) {
	a, err := u.alamats.FindByID(id)
	if err != nil {
		return nil, err
	}
	if a.IDUser != userID {
		return nil, domain.Forbidden("tidak dapat mengelola alamat milik user lain")
	}
	return a, nil
}

func (u *AlamatUsecase) Create(userID uint, in AlamatInput) (*domain.Alamat, error) {
	a := &domain.Alamat{
		IDUser: userID, JudulAlamat: in.JudulAlamat, NamaPenerima: in.NamaPenerima,
		NoTelp: in.NoTelp, DetailAlamat: in.DetailAlamat,
	}
	if err := u.alamats.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (u *AlamatUsecase) List(userID uint, judul string) ([]domain.Alamat, error) {
	return u.alamats.ListByUser(userID, judul)
}

func (u *AlamatUsecase) Get(userID, id uint) (*domain.Alamat, error) { return u.own(userID, id) }

// AlamatUpdateInput: semua field opsional (collection tidak mengirim judul_alamat
// saat update). Field yang kosong tidak mengubah data lama.
type AlamatUpdateInput struct {
	JudulAlamat  string `json:"judul_alamat"`
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp" validate:"omitempty,numeric,min=8,max=15"`
	DetailAlamat string `json:"detail_alamat"`
}

func (u *AlamatUsecase) Update(userID, id uint, in AlamatUpdateInput) (*domain.Alamat, error) {
	a, err := u.own(userID, id)
	if err != nil {
		return nil, err
	}
	if in.JudulAlamat != "" {
		a.JudulAlamat = in.JudulAlamat
	}
	if in.NamaPenerima != "" {
		a.NamaPenerima = in.NamaPenerima
	}
	if in.NoTelp != "" {
		a.NoTelp = in.NoTelp
	}
	if in.DetailAlamat != "" {
		a.DetailAlamat = in.DetailAlamat
	}
	if err := u.alamats.Update(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (u *AlamatUsecase) Delete(userID, id uint) error {
	if _, err := u.own(userID, id); err != nil {
		return err
	}
	return u.alamats.Delete(id)
}

// ---------- Kategori (akses admin dijaga di middleware) ----------

// KategoriInput mengikuti collection Postman (`nama_category`);
// `nama_kategori` tetap diterima sebagai alias.
type KategoriInput struct {
	NamaCategory string `json:"nama_category" validate:"required_without=NamaKategori"`
	NamaKategori string `json:"nama_kategori"`
}

func (k KategoriInput) Nama() string {
	if k.NamaCategory != "" {
		return k.NamaCategory
	}
	return k.NamaKategori
}

type KategoriUsecase struct{ repo domain.KategoriRepository }

func NewKategoriUsecase(r domain.KategoriRepository) *KategoriUsecase { return &KategoriUsecase{r} }

func (u *KategoriUsecase) Create(in KategoriInput) (*domain.Kategori, error) {
	k := &domain.Kategori{NamaKategori: in.Nama()}
	if err := u.repo.Create(k); err != nil {
		return nil, err
	}
	return k, nil
}

func (u *KategoriUsecase) List() ([]domain.Kategori, error) { return u.repo.List() }

func (u *KategoriUsecase) Get(id uint) (*domain.Kategori, error) { return u.repo.FindByID(id) }

func (u *KategoriUsecase) Update(id uint, in KategoriInput) (*domain.Kategori, error) {
	k, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	k.NamaKategori = in.Nama()
	if err := u.repo.Update(k); err != nil {
		return nil, err
	}
	return k, nil
}

func (u *KategoriUsecase) Delete(id uint) error {
	if _, err := u.repo.FindByID(id); err != nil {
		return err
	}
	return u.repo.Delete(id)
}

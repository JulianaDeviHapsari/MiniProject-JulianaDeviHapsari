package mysql

import (
	"errors"

	"gorm.io/gorm"

	"evermos-mini/internal/domain"
)

// ---------- Toko ----------

type tokoRepo struct{ db *gorm.DB }

func NewTokoRepository(db *gorm.DB) domain.TokoRepository { return &tokoRepo{db} }

func (r *tokoRepo) FindByID(id uint) (*domain.Toko, error) {
	var t domain.Toko
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, mapErr(err, "toko tidak ditemukan")
	}
	return &t, nil
}

func (r *tokoRepo) FindByUserID(userID uint) (*domain.Toko, error) {
	var t domain.Toko
	if err := r.db.Where("id_user = ?", userID).First(&t).Error; err != nil {
		return nil, mapErr(err, "toko tidak ditemukan")
	}
	return &t, nil
}

func (r *tokoRepo) List(nama string, page domain.Page) ([]domain.Toko, int64, error) {
	var (
		list  []domain.Toko
		total int64
	)
	q := r.db.Model(&domain.Toko{})
	if nama != "" {
		q = q.Where("nama_toko LIKE ?", "%"+nama+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("id ASC").Limit(page.Limit).Offset(page.Offset()).Find(&list).Error
	return list, total, err
}

func (r *tokoRepo) Update(t *domain.Toko) error { return r.db.Save(t).Error }

// ---------- Alamat ----------

type alamatRepo struct{ db *gorm.DB }

func NewAlamatRepository(db *gorm.DB) domain.AlamatRepository { return &alamatRepo{db} }

func (r *alamatRepo) Create(a *domain.Alamat) error { return r.db.Create(a).Error }

func (r *alamatRepo) FindByID(id uint) (*domain.Alamat, error) {
	var a domain.Alamat
	if err := r.db.First(&a, id).Error; err != nil {
		return nil, mapErr(err, "alamat tidak ditemukan")
	}
	return &a, nil
}

func (r *alamatRepo) ListByUser(userID uint, judul string) ([]domain.Alamat, error) {
	var list []domain.Alamat
	q := r.db.Where("id_user = ?", userID)
	if judul != "" {
		q = q.Where("judul_alamat LIKE ?", "%"+judul+"%")
	}
	err := q.Order("id ASC").Find(&list).Error
	return list, err
}

func (r *alamatRepo) Update(a *domain.Alamat) error { return r.db.Save(a).Error }

func (r *alamatRepo) Delete(id uint) error { return r.db.Delete(&domain.Alamat{}, id).Error }

// ---------- Kategori ----------

type kategoriRepo struct{ db *gorm.DB }

func NewKategoriRepository(db *gorm.DB) domain.KategoriRepository { return &kategoriRepo{db} }

// dupToConflict mengubah pelanggaran unique index menjadi error domain 409.
func dupToConflict(err error, msg string) error {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return domain.Conflict(msg)
	}
	return err
}

func (r *kategoriRepo) Create(k *domain.Kategori) error {
	return dupToConflict(r.db.Create(k).Error, "nama kategori sudah ada")
}

func (r *kategoriRepo) FindByID(id uint) (*domain.Kategori, error) {
	var k domain.Kategori
	if err := r.db.First(&k, id).Error; err != nil {
		return nil, mapErr(err, "kategori tidak ditemukan")
	}
	return &k, nil
}

func (r *kategoriRepo) List() ([]domain.Kategori, error) {
	var list []domain.Kategori
	err := r.db.Order("id ASC").Find(&list).Error
	return list, err
}

func (r *kategoriRepo) Update(k *domain.Kategori) error {
	return dupToConflict(r.db.Save(k).Error, "nama kategori sudah ada")
}

func (r *kategoriRepo) Delete(id uint) error { return r.db.Delete(&domain.Kategori{}, id).Error }

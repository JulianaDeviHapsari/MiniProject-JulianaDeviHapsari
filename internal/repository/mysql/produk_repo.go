package mysql

import (
	"gorm.io/gorm"

	"evermos-mini/internal/domain"
)

type produkRepo struct{ db *gorm.DB }

func NewProdukRepository(db *gorm.DB) domain.ProdukRepository { return &produkRepo{db} }

func (r *produkRepo) Create(p *domain.Produk) error { return r.db.Create(p).Error }

func (r *produkRepo) FindByID(id uint) (*domain.Produk, error) {
	var p domain.Produk
	err := r.db.Preload("Toko").Preload("Kategori").Preload("Foto").First(&p, id).Error
	if err != nil {
		return nil, mapErr(err, "produk tidak ditemukan")
	}
	return &p, nil
}

func (r *produkRepo) List(f domain.ProdukFilter) ([]domain.Produk, int64, error) {
	var (
		list  []domain.Produk
		total int64
	)
	q := r.db.Model(&domain.Produk{})
	if f.NamaProduk != "" {
		q = q.Where("nama_produk LIKE ?", "%"+f.NamaProduk+"%")
	}
	if f.IDKategori != 0 {
		q = q.Where("id_kategori = ?", f.IDKategori)
	}
	if f.IDToko != 0 {
		q = q.Where("id_toko = ?", f.IDToko)
	}
	if f.MinHarga > 0 {
		q = q.Where("harga_konsumen >= ?", f.MinHarga)
	}
	if f.MaxHarga > 0 {
		q = q.Where("harga_konsumen <= ?", f.MaxHarga)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Preload("Toko").Preload("Kategori").Preload("Foto").
		Order("id DESC").Limit(f.Limit).Offset(f.Offset()).Find(&list).Error
	return list, total, err
}

func (r *produkRepo) Update(p *domain.Produk) error {
	// Omit association agar Save tidak menimpa toko/kategori/foto.
	return r.db.Omit("Toko", "Kategori", "Foto").Save(p).Error
}

func (r *produkRepo) AddFoto(fotos []domain.FotoProduk) error {
	if len(fotos) == 0 {
		return nil
	}
	return r.db.Create(&fotos).Error
}

// Delete adalah soft delete (gorm.DeletedAt).
func (r *produkRepo) Delete(id uint) error { return r.db.Delete(&domain.Produk{}, id).Error }

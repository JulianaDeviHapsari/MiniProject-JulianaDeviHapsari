package mysql

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"evermos-mini/internal/domain"
)

type trxRepo struct{ db *gorm.DB }

func NewTrxRepository(db *gorm.DB) domain.TrxRepository { return &trxRepo{db} }

func (r *trxRepo) Create(userID, alamatID uint, methodBayar string, items []domain.TrxItemInput) (*domain.Trx, error) {
	var result *domain.Trx

	err := r.db.Transaction(func(tx *gorm.DB) error {
		trx := &domain.Trx{
			IDUser:      userID,
			IDAlamat:    alamatID,
			MethodBayar: methodBayar,
			KodeInvoice: "INV-" + strings.ToUpper(strconv.FormatInt(time.Now().UnixNano(), 36)),
		}

		for _, it := range items {
			// Kunci baris produk agar stok tidak balapan antar transaksi.
			var p domain.Produk
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&p, it.IDProduk).Error
			if err != nil {
				return mapErr(err, fmt.Sprintf("produk %d tidak ditemukan", it.IDProduk))
			}
			if p.Stok < it.Kuantitas {
				return domain.Conflict(fmt.Sprintf("stok produk '%s' tidak mencukupi", p.NamaProduk))
			}

			var foto domain.FotoProduk
			_ = tx.Where("id_produk = ?", p.ID).Order("id ASC").Limit(1).Find(&foto).Error

			// Snapshot produk ke log_produk.
			logProduk := &domain.LogProduk{
				IDProduk:      p.ID,
				NamaProduk:    p.NamaProduk,
				Slug:          p.Slug,
				HargaReseller: p.HargaReseller,
				HargaKonsumen: p.HargaKonsumen,
				Deskripsi:     p.Deskripsi,
				UrlFoto:       foto.Url,
				IDToko:        p.IDToko,
				IDKategori:    p.IDKategori,
			}
			if err := tx.Create(logProduk).Error; err != nil {
				return err
			}

			if err := tx.Model(&domain.Produk{}).Where("id = ?", p.ID).
				Update("stok", gorm.Expr("stok - ?", it.Kuantitas)).Error; err != nil {
				return err
			}

			subtotal := p.HargaKonsumen * it.Kuantitas
			trx.HargaTotal += subtotal
			trx.Detail = append(trx.Detail, domain.DetailTrx{
				IDLogProduk: logProduk.ID,
				IDToko:      p.IDToko,
				Kuantitas:   it.Kuantitas,
				HargaTotal:  subtotal,
			})
		}

		if err := tx.Create(trx).Error; err != nil {
			return err
		}
		result = trx
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.FindByID(result.ID)
}

func (r *trxRepo) preloaded() *gorm.DB {
	return r.db.Preload("Alamat").Preload("Detail").
		Preload("Detail.LogProduk").Preload("Detail.Toko")
}

func (r *trxRepo) FindByID(id uint) (*domain.Trx, error) {
	var t domain.Trx
	if err := r.preloaded().First(&t, id).Error; err != nil {
		return nil, mapErr(err, "transaksi tidak ditemukan")
	}
	return &t, nil
}

func (r *trxRepo) ListByUser(userID uint, search string, page domain.Page) ([]domain.Trx, int64, error) {
	var (
		list  []domain.Trx
		total int64
	)
	q := r.db.Model(&domain.Trx{}).Where("id_user = ?", userID)
	if search != "" {
		q = q.Where("kode_invoice LIKE ?", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Preload("Alamat").Preload("Detail").
		Preload("Detail.LogProduk").Preload("Detail.Toko").
		Order("id DESC").Limit(page.Limit).Offset(page.Offset()).Find(&list).Error
	return list, total, err
}

package mysql

import (
	"errors"

	"gorm.io/gorm"

	"evermos-mini/internal/domain"
)

// mapErr mengubah error GORM menjadi error domain.
func mapErr(err error, notFoundMsg string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.NotFound(notFoundMsg)
	}
	return err
}

type userRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) domain.UserRepository { return &userRepo{db} }

func (r *userRepo) CreateWithToko(user *domain.User, toko *domain.Toko) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		toko.IDUser = user.ID
		return tx.Create(toko).Error
	})
}

func (r *userRepo) FindByID(id uint) (*domain.User, error) {
	var u domain.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, mapErr(err, "user tidak ditemukan")
	}
	return &u, nil
}

func (r *userRepo) FindByNoTelp(noTelp string) (*domain.User, error) {
	var u domain.User
	if err := r.db.Where("no_telp = ?", noTelp).First(&u).Error; err != nil {
		return nil, mapErr(err, "user tidak ditemukan")
	}
	return &u, nil
}

func (r *userRepo) exists(field, value string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.Model(&domain.User{}).Where(field+" = ?", value)
	if excludeID != 0 {
		q = q.Where("id <> ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

func (r *userRepo) ExistsByEmail(email string, excludeID uint) (bool, error) {
	return r.exists("email", email, excludeID)
}

func (r *userRepo) ExistsByNoTelp(noTelp string, excludeID uint) (bool, error) {
	return r.exists("no_telp", noTelp, excludeID)
}

func (r *userRepo) Update(user *domain.User) error { return r.db.Save(user).Error }

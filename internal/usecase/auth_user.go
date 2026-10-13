package usecase

import (
	"fmt"
	"strings"
	"time"

	"evermos-mini/internal/domain"
	"evermos-mini/internal/pkg/security"
)

// ---------- Input DTO ----------

type RegisterInput struct {
	Nama         string `json:"nama" validate:"required"`
	KataSandi    string `json:"kata_sandi" validate:"required,min=8"`
	NoTelp       string `json:"no_telp" validate:"required,numeric,min=8,max=15"`
	TanggalLahir string `json:"tanggal_lahir" validate:"required"`
	JenisKelamin string `json:"jenis_kelamin"`
	Tentang      string `json:"tentang"`
	Pekerjaan    string `json:"pekerjaan"`
	Email        string `json:"email" validate:"required,email"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
}

type LoginInput struct {
	NoTelp    string `json:"no_telp" validate:"required"`
	KataSandi string `json:"kata_sandi" validate:"required"`
}

type UpdateUserInput struct {
	Nama         string `json:"nama" validate:"required"`
	KataSandi    string `json:"kata_sandi" validate:"omitempty,min=8"`
	NoTelp       string `json:"no_telp" validate:"required,numeric,min=8,max=15"`
	TanggalLahir string `json:"tanggal_lahir" validate:"required"`
	JenisKelamin string `json:"jenis_kelamin"`
	Tentang      string `json:"tentang"`
	Pekerjaan    string `json:"pekerjaan"`
	Email        string `json:"email" validate:"required,email"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
}

type LoginResult struct {
	Nama       string `json:"nama"`
	NoTelp     string `json:"no_telp"`
	Email      string `json:"email"`
	IDProvinsi string `json:"id_provinsi"`
	IDKota     string `json:"id_kota"`
	Token      string `json:"token"`
}

// ---------- Auth ----------

type AuthUsecase struct {
	users  domain.UserRepository
	tokens *security.TokenManager
}

func NewAuthUsecase(users domain.UserRepository, tokens *security.TokenManager) *AuthUsecase {
	return &AuthUsecase{users: users, tokens: tokens}
}

func parseTanggal(s string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02", "02/01/2006"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, domain.BadRequest("format tanggal_lahir harus YYYY-MM-DD atau DD/MM/YYYY")
}

// Register membuat user sekaligus toko miliknya dalam satu transaksi.
func (u *AuthUsecase) Register(in RegisterInput) error {
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	tgl, err := parseTanggal(in.TanggalLahir)
	if err != nil {
		return err
	}
	if ok, err := u.users.ExistsByEmail(in.Email, 0); err != nil {
		return err
	} else if ok {
		return domain.Conflict("email sudah terdaftar")
	}
	if ok, err := u.users.ExistsByNoTelp(in.NoTelp, 0); err != nil {
		return err
	} else if ok {
		return domain.Conflict("no telp sudah terdaftar")
	}

	hash, err := security.HashPassword(in.KataSandi)
	if err != nil {
		return err
	}
	user := &domain.User{
		Nama: in.Nama, KataSandi: hash, NoTelp: in.NoTelp, Email: in.Email,
		TanggalLahir: tgl, JenisKelamin: in.JenisKelamin, Tentang: in.Tentang,
		Pekerjaan: in.Pekerjaan, IDProvinsi: in.IDProvinsi, IDKota: in.IDKota,
	}
	toko := &domain.Toko{NamaToko: fmt.Sprintf("Toko %s", in.Nama)}
	return u.users.CreateWithToko(user, toko)
}

func (u *AuthUsecase) Login(in LoginInput) (*LoginResult, error) {
	user, err := u.users.FindByNoTelp(in.NoTelp)
	if err != nil || !security.CheckPassword(user.KataSandi, in.KataSandi) {
		// Pesan sama untuk dua kasus agar tidak membocorkan akun yang terdaftar.
		return nil, domain.Unauthorized("no telp atau kata sandi salah")
	}
	token, err := u.tokens.Generate(user.ID, user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		Nama: user.Nama, NoTelp: user.NoTelp, Email: user.Email,
		IDProvinsi: user.IDProvinsi, IDKota: user.IDKota, Token: token,
	}, nil
}

// ---------- User ----------

type UserUsecase struct{ users domain.UserRepository }

func NewUserUsecase(users domain.UserRepository) *UserUsecase { return &UserUsecase{users} }

// GetMe hanya mengambil data milik userID dari token; tidak ada jalur ke user lain.
func (u *UserUsecase) GetMe(userID uint) (*domain.User, error) { return u.users.FindByID(userID) }

func (u *UserUsecase) UpdateMe(userID uint, in UpdateUserInput) (*domain.User, error) {
	user, err := u.users.FindByID(userID)
	if err != nil {
		return nil, err
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	tgl, err := parseTanggal(in.TanggalLahir)
	if err != nil {
		return nil, err
	}
	if ok, err := u.users.ExistsByEmail(in.Email, userID); err != nil {
		return nil, err
	} else if ok {
		return nil, domain.Conflict("email sudah digunakan user lain")
	}
	if ok, err := u.users.ExistsByNoTelp(in.NoTelp, userID); err != nil {
		return nil, err
	} else if ok {
		return nil, domain.Conflict("no telp sudah digunakan user lain")
	}

	if in.KataSandi != "" {
		hash, err := security.HashPassword(in.KataSandi)
		if err != nil {
			return nil, err
		}
		user.KataSandi = hash
	}
	user.Nama, user.NoTelp, user.Email = in.Nama, in.NoTelp, in.Email
	user.TanggalLahir, user.JenisKelamin = tgl, in.JenisKelamin
	user.Tentang, user.Pekerjaan = in.Tentang, in.Pekerjaan
	user.IDProvinsi, user.IDKota = in.IDProvinsi, in.IDKota

	if err := u.users.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

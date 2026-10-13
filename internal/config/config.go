package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"evermos-mini/internal/domain"
)

type Config struct {
	AppPort     string
	DBUser      string
	DBPass      string
	DBHost      string
	DBPort      string
	DBName      string
	JWTSecret   string
	JWTTTLHours int
	UploadDir   string
}

func Load() (*Config, error) {
	_ = godotenv.Load() // .env opsional

	cfg := &Config{
		AppPort:   get("APP_PORT", "8080"),
		DBUser:    get("DB_USER", "root"),
		DBPass:    get("DB_PASS", ""),
		DBHost:    get("DB_HOST", "127.0.0.1"),
		DBPort:    get("DB_PORT", "3306"),
		DBName:    get("DB_NAME", "evermos_mini"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		UploadDir: get("UPLOAD_DIR", "uploads"),
	}
	ttl, err := strconv.Atoi(get("JWT_TTL_HOURS", "24"))
	if err != nil || ttl <= 0 {
		return nil, fmt.Errorf("JWT_TTL_HOURS tidak valid")
	}
	cfg.JWTTTLHours = ttl
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET wajib diisi di .env")
	}
	return cfg, nil
}

func get(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func ConnectDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		// Ubah error driver (mis. MySQL 1062) menjadi gorm.ErrDuplicatedKey.
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{}, &domain.Toko{}, &domain.Alamat{}, &domain.Kategori{},
		&domain.Produk{}, &domain.FotoProduk{}, &domain.LogProduk{},
		&domain.Trx{}, &domain.DetailTrx{},
	)
}

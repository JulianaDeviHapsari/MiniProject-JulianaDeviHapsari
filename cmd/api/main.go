package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"evermos-mini/internal/config"
	httpdelivery "evermos-mini/internal/delivery/http"
	"evermos-mini/internal/pkg/security"
	"evermos-mini/internal/repository/mysql"
	"evermos-mini/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := config.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		log.Fatalf("upload dir: %v", err)
	}

	tokens := security.NewTokenManager(cfg.JWTSecret, cfg.JWTTTLHours)

	// Repository -> Usecase -> Handler (dependency mengarah ke dalam).
	userRepo := mysql.NewUserRepository(db)
	tokoRepo := mysql.NewTokoRepository(db)
	alamatRepo := mysql.NewAlamatRepository(db)
	kategoriRepo := mysql.NewKategoriRepository(db)
	produkRepo := mysql.NewProdukRepository(db)
	trxRepo := mysql.NewTrxRepository(db)

	handlers := httpdelivery.Handlers{
		Auth:     httpdelivery.NewAuthHandler(usecase.NewAuthUsecase(userRepo, tokens)),
		User:     httpdelivery.NewUserHandler(usecase.NewUserUsecase(userRepo)),
		Alamat:   httpdelivery.NewAlamatHandler(usecase.NewAlamatUsecase(alamatRepo)),
		Toko:     httpdelivery.NewTokoHandler(usecase.NewTokoUsecase(tokoRepo), cfg.UploadDir),
		Kategori: httpdelivery.NewKategoriHandler(usecase.NewKategoriUsecase(kategoriRepo)),
		Produk:   httpdelivery.NewProdukHandler(usecase.NewProdukUsecase(produkRepo, tokoRepo, kategoriRepo), cfg.UploadDir),
		Trx:      httpdelivery.NewTrxHandler(usecase.NewTrxUsecase(trxRepo, alamatRepo)),
		Wilayah:  httpdelivery.NewWilayahHandler(usecase.NewWilayahUsecase()),
	}

	app := fiber.New(fiber.Config{BodyLimit: 20 * 1024 * 1024})
	app.Use(recover.New(), logger.New())
	app.Static("/uploads", cfg.UploadDir)
	httpdelivery.RegisterRoutes(app, handlers, tokens)

	log.Fatal(app.Listen(":" + cfg.AppPort))
}

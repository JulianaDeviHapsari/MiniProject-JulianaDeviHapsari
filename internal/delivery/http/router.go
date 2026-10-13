package http

import (
	"github.com/gofiber/fiber/v2"

	"evermos-mini/internal/pkg/security"
)

type Handlers struct {
	Auth     *AuthHandler
	User     *UserHandler
	Alamat   *AlamatHandler
	Toko     *TokoHandler
	Kategori *KategoriHandler
	Produk   *ProdukHandler
	Trx      *TrxHandler
	Wilayah  *WilayahHandler
}

// RegisterRoutes memasang semua endpoint sesuai collection Postman.
func RegisterRoutes(app *fiber.App, h Handlers, tokens *security.TokenManager) {
	auth := AuthMiddleware(tokens)
	admin := AdminMiddleware()

	// Auth (publik)
	app.Post("/auth/register", h.Auth.Register)
	app.Post("/auth/login", h.Auth.Login)

	// User & alamat
	user := app.Group("/user", auth)
	user.Get("/", h.User.GetMe)
	user.Put("/", h.User.UpdateMe)
	user.Get("/alamat", h.Alamat.List)
	user.Get("/alamat/:id", h.Alamat.Get)
	user.Post("/alamat", h.Alamat.Create)
	user.Put("/alamat/:id", h.Alamat.Update)
	user.Delete("/alamat/:id", h.Alamat.Delete)

	// Toko. "/toko/my" harus didaftarkan sebelum "/toko/:id_toko".
	// Katalog (list/detail) publik; mengubah toko hanya oleh pemilik.
	app.Get("/toko", h.Toko.List)
	app.Get("/toko/my", auth, h.Toko.GetMine)
	app.Get("/toko/:id_toko", h.Toko.Get)
	app.Put("/toko/:id_toko", auth, h.Toko.Update)

	// Kategori: baca publik, kelola (POST/PUT/DELETE) khusus admin.
	app.Get("/category", h.Kategori.List)
	app.Get("/category/:id", h.Kategori.Get)
	app.Post("/category", auth, admin, h.Kategori.Create)
	app.Put("/category/:id", auth, admin, h.Kategori.Update)
	app.Delete("/category/:id", auth, admin, h.Kategori.Delete)

	// Produk: katalog publik, kelola oleh pemilik toko.
	app.Get("/product", h.Produk.List)
	app.Get("/product/:id", h.Produk.Get)
	app.Post("/product", auth, h.Produk.Create)
	app.Put("/product/:id", auth, h.Produk.Update)
	app.Delete("/product/:id", auth, h.Produk.Delete)

	// Transaksi
	app.Get("/trx", auth, h.Trx.List)
	app.Get("/trx/:id", auth, h.Trx.Get)
	app.Post("/trx", auth, h.Trx.Create)

	// Provinsi & kota (publik, sumber: emsifa.com)
	app.Get("/provcity/listprovincies", h.Wilayah.ListProvinces)
	app.Get("/provcity/listcities/:prov_id", h.Wilayah.ListCities)
	app.Get("/provcity/detailprovince/:prov_id", h.Wilayah.DetailProvince)
	app.Get("/provcity/detailcity/:city_id", h.Wilayah.DetailCity)
}

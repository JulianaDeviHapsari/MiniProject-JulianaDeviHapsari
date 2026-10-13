package http

import (
	"mime/multipart"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"evermos-mini/internal/domain"
	"evermos-mini/internal/pkg/upload"
	"evermos-mini/internal/usecase"
)

// ================= Auth =================

type AuthHandler struct{ uc *usecase.AuthUsecase }

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler { return &AuthHandler{uc} }

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var in usecase.RegisterInput
	if err := bindAndValidate(c, &in); err != nil {
		return handleError(c, "Failed to POST data", err)
	}
	if err := h.uc.Register(in); err != nil {
		return handleError(c, "Failed to POST data", err)
	}
	return created(c, "Succeed to POST data", "Register Succeed")
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var in usecase.LoginInput
	if err := bindAndValidate(c, &in); err != nil {
		return handleError(c, "Failed to POST data", err)
	}
	res, err := h.uc.Login(in)
	if err != nil {
		return handleError(c, "Failed to POST data", err)
	}
	return ok(c, "Succeed to POST data", res)
}

// ================= User =================

type UserHandler struct{ uc *usecase.UserUsecase }

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler { return &UserHandler{uc} }

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	u, err := h.uc.GetMe(currentUserID(c))
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", u)
}

func (h *UserHandler) UpdateMe(c *fiber.Ctx) error {
	var in usecase.UpdateUserInput
	if err := bindAndValidate(c, &in); err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	u, err := h.uc.UpdateMe(currentUserID(c), in)
	if err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	return ok(c, "Succeed to UPDATE data", u)
}

// ================= Alamat =================

type AlamatHandler struct{ uc *usecase.AlamatUsecase }

func NewAlamatHandler(uc *usecase.AlamatUsecase) *AlamatHandler { return &AlamatHandler{uc} }

func (h *AlamatHandler) List(c *fiber.Ctx) error {
	list, err := h.uc.List(currentUserID(c), c.Query("judul_alamat"))
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", list)
}

func (h *AlamatHandler) Get(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	a, err := h.uc.Get(currentUserID(c), id)
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", a)
}

func (h *AlamatHandler) Create(c *fiber.Ctx) error {
	var in usecase.AlamatInput
	if err := bindAndValidate(c, &in); err != nil {
		return handleError(c, "Failed to POST data", err)
	}
	a, err := h.uc.Create(currentUserID(c), in)
	if err != nil {
		return handleError(c, "Failed to POST data", err)
	}
	return created(c, "Succeed to POST data", a)
}

func (h *AlamatHandler) Update(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	var in usecase.AlamatUpdateInput
	if err := bindAndValidate(c, &in); err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	a, err := h.uc.Update(currentUserID(c), id, in)
	if err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	return ok(c, "Succeed to UPDATE data", a)
}

func (h *AlamatHandler) Delete(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return handleError(c, "Failed to DELETE data", err)
	}
	if err := h.uc.Delete(currentUserID(c), id); err != nil {
		return handleError(c, "Failed to DELETE data", err)
	}
	return ok(c, "Succeed to DELETE data", "")
}

// ================= Toko =================

type TokoHandler struct {
	uc        *usecase.TokoUsecase
	uploadDir string
}

func NewTokoHandler(uc *usecase.TokoUsecase, uploadDir string) *TokoHandler {
	return &TokoHandler{uc, uploadDir}
}

func (h *TokoHandler) List(c *fiber.Ctx) error {
	res, err := h.uc.List(c.Query("nama"), pageParams(c))
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", res)
}

func (h *TokoHandler) Get(c *fiber.Ctx) error {
	id, err := idParam(c, "id_toko")
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	t, err := h.uc.GetByID(id)
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", t)
}

func (h *TokoHandler) GetMine(c *fiber.Ctx) error {
	t, err := h.uc.GetMine(currentUserID(c))
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", t)
}

// Update menerima multipart/form-data: nama_toko (text) dan photo (file).
func (h *TokoHandler) Update(c *fiber.Ctx) error {
	id, err := idParam(c, "id_toko")
	if err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	urlFoto := ""
	if fh, ferr := c.FormFile("photo"); ferr == nil {
		urlFoto, err = upload.SaveImage(c, fh, h.uploadDir)
		if err != nil {
			return handleError(c, "Failed to UPDATE data", err)
		}
	}
	t, err := h.uc.Update(currentUserID(c), id, c.FormValue("nama_toko"), urlFoto)
	if err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	return ok(c, "Succeed to UPDATE data", t)
}

// ================= Kategori =================

type KategoriHandler struct{ uc *usecase.KategoriUsecase }

func NewKategoriHandler(uc *usecase.KategoriUsecase) *KategoriHandler { return &KategoriHandler{uc} }

func (h *KategoriHandler) List(c *fiber.Ctx) error {
	list, err := h.uc.List()
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", list)
}

func (h *KategoriHandler) Get(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	k, err := h.uc.Get(id)
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", k)
}

func (h *KategoriHandler) Create(c *fiber.Ctx) error {
	var in usecase.KategoriInput
	if err := bindAndValidate(c, &in); err != nil {
		return handleError(c, "Failed to POST data", err)
	}
	k, err := h.uc.Create(in)
	if err != nil {
		return handleError(c, "Failed to POST data", err)
	}
	return created(c, "Succeed to POST data", k)
}

func (h *KategoriHandler) Update(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	var in usecase.KategoriInput
	if err := bindAndValidate(c, &in); err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	k, err := h.uc.Update(id, in)
	if err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	return ok(c, "Succeed to UPDATE data", k)
}

func (h *KategoriHandler) Delete(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return handleError(c, "Failed to DELETE data", err)
	}
	if err := h.uc.Delete(id); err != nil {
		return handleError(c, "Failed to DELETE data", err)
	}
	return ok(c, "Succeed to DELETE data", "")
}

// ================= Produk =================

type ProdukHandler struct {
	uc        *usecase.ProdukUsecase
	uploadDir string
}

func NewProdukHandler(uc *usecase.ProdukUsecase, uploadDir string) *ProdukHandler {
	return &ProdukHandler{uc, uploadDir}
}

func queryUint(c *fiber.Ctx, key string) uint {
	n, _ := strconv.ParseUint(c.Query(key), 10, 64)
	return uint(n)
}

func queryInt(c *fiber.Ctx, key string) int {
	n, _ := strconv.Atoi(c.Query(key))
	return n
}

func (h *ProdukHandler) List(c *fiber.Ctx) error {
	res, err := h.uc.List(domain.ProdukFilter{
		Page:       pageParams(c),
		NamaProduk: c.Query("nama_produk"),
		IDKategori: queryUint(c, "category_id"),
		IDToko:     queryUint(c, "toko_id"),
		MinHarga:   queryInt(c, "min_harga"),
		MaxHarga:   queryInt(c, "max_harga"),
	})
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", res)
}

func (h *ProdukHandler) Get(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	p, err := h.uc.Get(id)
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", p)
}

// parseProdukForm membaca multipart/form-data dan menyimpan semua file `photos`.
func (h *ProdukHandler) parseProdukForm(c *fiber.Ctx) (usecase.ProdukInput, []string, error) {
	atoi := func(s string) int { n, _ := strconv.Atoi(s); return n }
	kat, _ := strconv.ParseUint(c.FormValue("category_id"), 10, 64)
	in := usecase.ProdukInput{
		NamaProduk:    c.FormValue("nama_produk"),
		IDKategori:    uint(kat),
		HargaReseller: atoi(c.FormValue("harga_reseller")),
		HargaKonsumen: atoi(c.FormValue("harga_konsumen")),
		Stok:          atoi(c.FormValue("stok")),
		Deskripsi:     c.FormValue("deskripsi"),
	}
	if err := validateStruct(in); err != nil {
		return in, nil, err
	}

	var files []*multipart.FileHeader
	if form, err := c.MultipartForm(); err == nil && form != nil {
		files = form.File["photos"]
	}
	var urls []string
	for _, fh := range files {
		url, err := upload.SaveImage(c, fh, h.uploadDir)
		if err != nil {
			return in, nil, err
		}
		urls = append(urls, url)
	}
	return in, urls, nil
}

func (h *ProdukHandler) Create(c *fiber.Ctx) error {
	in, urls, err := h.parseProdukForm(c)
	if err != nil {
		return handleError(c, "Failed to POST data", err)
	}
	p, err := h.uc.Create(currentUserID(c), in, urls)
	if err != nil {
		return handleError(c, "Failed to POST data", err)
	}
	return created(c, "Succeed to POST data", p)
}

func (h *ProdukHandler) Update(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	in, urls, err := h.parseProdukForm(c)
	if err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	p, err := h.uc.Update(currentUserID(c), id, in, urls)
	if err != nil {
		return handleError(c, "Failed to UPDATE data", err)
	}
	return ok(c, "Succeed to UPDATE data", p)
}

func (h *ProdukHandler) Delete(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return handleError(c, "Failed to DELETE data", err)
	}
	if err := h.uc.Delete(currentUserID(c), id); err != nil {
		return handleError(c, "Failed to DELETE data", err)
	}
	return ok(c, "Succeed to DELETE data", "")
}

// ================= Transaksi =================

type TrxHandler struct{ uc *usecase.TrxUsecase }

func NewTrxHandler(uc *usecase.TrxUsecase) *TrxHandler { return &TrxHandler{uc} }

func (h *TrxHandler) Create(c *fiber.Ctx) error {
	var in usecase.CreateTrxInput
	if err := bindAndValidate(c, &in); err != nil {
		return handleError(c, "Failed to POST data", err)
	}
	t, err := h.uc.Create(currentUserID(c), in)
	if err != nil {
		return handleError(c, "Failed to POST data", err)
	}
	return created(c, "Succeed to POST data", t)
}

func (h *TrxHandler) List(c *fiber.Ctx) error {
	res, err := h.uc.List(currentUserID(c), c.Query("search"), pageParams(c))
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", res)
}

func (h *TrxHandler) Get(c *fiber.Ctx) error {
	id, err := idParam(c, "id")
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	t, err := h.uc.Get(currentUserID(c), id)
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", t)
}

// ================= Provinsi & Kota =================

type WilayahHandler struct{ uc *usecase.WilayahUsecase }

func NewWilayahHandler(uc *usecase.WilayahUsecase) *WilayahHandler { return &WilayahHandler{uc} }

func (h *WilayahHandler) wrap(c *fiber.Ctx, data interface{}, err error) error {
	if err != nil {
		return handleError(c, "Failed to GET data", err)
	}
	return ok(c, "Succeed to GET data", data)
}

func (h *WilayahHandler) ListProvinces(c *fiber.Ctx) error {
	d, err := h.uc.ListProvinces()
	return h.wrap(c, d, err)
}

func (h *WilayahHandler) ListCities(c *fiber.Ctx) error {
	d, err := h.uc.ListCities(c.Params("prov_id"))
	return h.wrap(c, d, err)
}

func (h *WilayahHandler) DetailProvince(c *fiber.Ctx) error {
	d, err := h.uc.DetailProvince(c.Params("prov_id"))
	return h.wrap(c, d, err)
}

func (h *WilayahHandler) DetailCity(c *fiber.Ctx) error {
	d, err := h.uc.DetailCity(c.Params("city_id"))
	return h.wrap(c, d, err)
}

package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"evermos-mini/internal/domain"
)

const wilayahBase = "https://www.emsifa.com/api-wilayah-indonesia/api"

// WilayahUsecase membungkus API publik emsifa.com (data provinsi & kota/kabupaten).
type WilayahUsecase struct{ client *http.Client }

func NewWilayahUsecase() *WilayahUsecase {
	return &WilayahUsecase{client: &http.Client{Timeout: 10 * time.Second}}
}

func (u *WilayahUsecase) fetch(path string) (interface{}, error) {
	resp, err := u.client.Get(wilayahBase + path)
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi API wilayah: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.NotFound("data wilayah tidak ditemukan")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API wilayah mengembalikan status %d", resp.StatusCode)
	}
	var out interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (u *WilayahUsecase) ListProvinces() (interface{}, error) { return u.fetch("/provinces.json") }

func (u *WilayahUsecase) ListCities(provID string) (interface{}, error) {
	return u.fetch("/regencies/" + url.PathEscape(provID) + ".json")
}

func (u *WilayahUsecase) DetailProvince(id string) (interface{}, error) {
	return u.fetch("/province/" + url.PathEscape(id) + ".json")
}

func (u *WilayahUsecase) DetailCity(id string) (interface{}, error) {
	return u.fetch("/regency/" + url.PathEscape(id) + ".json")
}

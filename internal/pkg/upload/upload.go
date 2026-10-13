package upload

import (
	"crypto/rand"
	"encoding/hex"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"

	"evermos-mini/internal/domain"
)

const MaxSize = 5 * 1024 * 1024 // 5MB

var allowedExt = map[string]bool{".jpg": true, ".jpeg": true, ".png": true}

// SaveImage memvalidasi (ukuran & ekstensi) lalu menyimpan gambar ke dir,
// dan mengembalikan URL relatif yang dilayani di /uploads.
func SaveImage(c *fiber.Ctx, fh *multipart.FileHeader, dir string) (string, error) {
	if fh.Size > MaxSize {
		return "", domain.BadRequest("ukuran file maksimal 5MB")
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedExt[ext] {
		return "", domain.BadRequest("format file harus JPG, JPEG, atau PNG")
	}
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	name := hex.EncodeToString(buf) + ext
	if err := c.SaveFile(fh, filepath.Join(dir, name)); err != nil {
		return "", err
	}
	return "/uploads/" + name, nil
}

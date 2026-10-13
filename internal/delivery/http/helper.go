package http

import (
	"errors"
	"log"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"evermos-mini/internal/domain"
)

var validate = validator.New()

// Response adalah format respons seragam, mengikuti collection Postman.
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  []string    `json:"errors"`
	Data    interface{} `json:"data"`
}

func ok(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{Status: true, Message: message, Errors: nil, Data: data})
}

func created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{Status: true, Message: message, Errors: nil, Data: data})
}

func fail(c *fiber.Ctx, code int, message string, errs ...string) error {
	return c.Status(code).JSON(Response{Status: false, Message: message, Errors: errs, Data: nil})
}

// handleError memetakan error domain ke status HTTP. Error tak dikenal
// dicatat di server dan dibalas generik supaya detail internal tidak bocor.
func handleError(c *fiber.Ctx, message string, err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return fail(c, fiber.StatusNotFound, message, err.Error())
	case errors.Is(err, domain.ErrForbidden):
		return fail(c, fiber.StatusForbidden, message, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		return fail(c, fiber.StatusUnauthorized, message, err.Error())
	case errors.Is(err, domain.ErrConflict):
		return fail(c, fiber.StatusConflict, message, err.Error())
	case errors.Is(err, domain.ErrBadRequest):
		return fail(c, fiber.StatusBadRequest, message, err.Error())
	default:
		log.Printf("internal error on %s %s: %v", c.Method(), c.Path(), err)
		return fail(c, fiber.StatusInternalServerError, message, "terjadi kesalahan pada server")
	}
}

// bindAndValidate mem-parse body JSON lalu memvalidasi tag `validate`.
func bindAndValidate(c *fiber.Ctx, out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		return domain.BadRequest("body request tidak valid")
	}
	return validateStruct(out)
}

func validateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			msg := ""
			for i, fe := range verrs {
				if i > 0 {
					msg += "; "
				}
				msg += fe.Field() + " tidak valid (" + fe.Tag() + ")"
			}
			return domain.BadRequest(msg)
		}
		return domain.BadRequest(err.Error())
	}
	return nil
}

func pageParams(c *fiber.Ctx) domain.Page {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return domain.Page{Page: page, Limit: limit}
}

func idParam(c *fiber.Ctx, name string) (uint, error) {
	n, err := strconv.ParseUint(c.Params(name), 10, 64)
	if err != nil || n == 0 {
		return 0, domain.BadRequest("parameter " + name + " tidak valid")
	}
	return uint(n), nil
}

func currentUserID(c *fiber.Ctx) uint {
	id, _ := c.Locals("user_id").(uint)
	return id
}

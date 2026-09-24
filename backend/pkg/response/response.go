package response

import (
	"errors"
	"log"

	"backend/pkg/apperror"

	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type APIResponseWithMeta struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *fiber.Ctx, statusCode int, message string, data interface{}, meta interface{}) error {
	return c.Status(statusCode).JSON(APIResponseWithMeta{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *fiber.Ctx, statusCode int, message string, errorsDetail interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: false,
		Message: message,
		Errors:  errorsDetail,
	})
}

func HandleError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		if appErr.StatusCode >= 500 && appErr.RawErr != nil {
			log.Printf("[SERVER ERROR] %s: %v\n", appErr.Message, appErr.RawErr)
		}
		return Error(c, appErr.StatusCode, appErr.Message, appErr.Details)
	}

	log.Printf("[UNHANDLED ERROR] %v\n", err)
	return Error(c, fiber.StatusInternalServerError, "Terjadi kesalahan internal pada server", nil)
}

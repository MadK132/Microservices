package dto

import (
	"errors"
	"net/http"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/model"
)

type HTTPError struct {
	Code    int
	Message string
}

var (
	ErrInvalidCategory = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "Invalid category",
	}

	ErrInvalidDescription = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "Invalid description",
	}

	ErrProductNotFound = &HTTPError{
		Code:    http.StatusNotFound,
		Message: "Product not found",
	}

	ErrInvalidPrice = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "Price must be positive",
	}

	ErrInvalidStock = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "Stock must be non-negative",
	}

	ErrInvalidName = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "Name cannot be empty",
	}

	ErrProductAlreadyExists = &HTTPError{
		Code:    http.StatusConflict,
		Message: "Product already exists",
	}
)

func FromError(err error) *HTTPError {
	switch {
	case errors.Is(err, model.ErrInvalidCategory):
		return ErrInvalidCategory
	case errors.Is(err, model.ErrInvalidDescription):
		return ErrInvalidDescription
	case errors.Is(err, model.ErrProductNotFound):
		return ErrProductNotFound
	case errors.Is(err, model.ErrInvalidPrice):
		return ErrInvalidPrice
	case errors.Is(err, model.ErrInvalidStock):
		return ErrInvalidStock
	case errors.Is(err, model.ErrInvalidName):
		return ErrInvalidName
	case errors.Is(err, model.ErrProductAlreadyExists):
		return ErrProductAlreadyExists
	default:
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
}

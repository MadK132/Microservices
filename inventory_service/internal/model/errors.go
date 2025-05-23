package model

import (
	"errors"
)

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrInvalidCategory      = errors.New("invalid category")
	ErrInsufficientStock    = errors.New("insufficient stock")
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrInvalidDescription   = errors.New("invalid description")
	ErrInvalidPrice         = errors.New("price must be positive")
	ErrInvalidStock         = errors.New("stock must be non-negative")
	ErrInvalidName          = errors.New("name cannot be empty")
	ErrDiscnoutNotFound = errors.New("Discount not found")
	ErrReviewNotFound = errors.New("Review not found")
)

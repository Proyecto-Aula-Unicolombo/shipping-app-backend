package utils

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

const (
	DefaultLimit = 10
	MaxLimit     = 100
)

type PaginationParams struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"-"`
}

type PaginationResponse struct {
	Items      interface{} `json:"items"`
	TotalPages int         `json:"total_pages"`
	TotalItems int         `json:"total_items"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
}

func GetPaginationParams(ctx fiber.Ctx) PaginationParams {
	page := getIntQuery(ctx, "page", 1)
	limit := getIntQuery(ctx, "limit", DefaultLimit)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = DefaultLimit
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}
	offset := (page - 1) * limit

	return PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

func NewPaginationResponse(items interface{}, totalItems, page, limit int) PaginationResponse {
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	return PaginationResponse{
		Items:      items,
		TotalPages: totalPages,
		TotalItems: totalItems,
		Page:       page,
		Limit:      limit,
	}
}

func getIntQuery(ctx fiber.Ctx, key string, defaultValue int) int {
	value := ctx.Query(key)
	if value == "" {
		return defaultValue
	}

	intvalue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intvalue
}

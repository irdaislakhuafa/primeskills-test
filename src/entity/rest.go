package entity

import "math"

type (
	HTTPMessage struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	MetaError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	Meta struct {
		Path       string     `json:"path"`
		StatusCode int        `json:"statusCode"`
		StatusStr  string     `json:"statusStr"`
		Message    string     `json:"message"`
		Timestamp  string     `json:"timestamp"`
		Error      *MetaError `json:"error"`
		RequestID  string     `json:"requestId"`
	}

	Pagination struct {
		CurrentPage     int      `json:"currentPage"`
		CurrentElements int      `json:"currentElements"`
		TotalPages      int      `json:"totalPages"`
		TotalElements   int      `json:"totalElements"`
		SortBy          []string `json:"sortBy"`
		CursorStart     *string  `json:"cursorStart"`
		CursorEnd       *string  `json:"cursorEnd"`
	}

	HTTPRes struct {
		Message    HTTPMessage `json:"message"`
		Meta       Meta        `json:"meta"`
		Data       any         `json:"data,omitempty"`
		Pagination *Pagination `json:"pagination,omitempty"`
	}
)

func GenPagination(page, limit, totalItems int) Pagination {
	if limit <= 0 {
		limit = 15
	}
	if page < 0 {
		page = 0
	}

	if limit > totalItems {
		limit = totalItems
	}

	totalPages := math.Ceil(float64(totalItems) / float64(limit))
	return Pagination{
		CurrentPage:     page,
		CurrentElements: limit,
		TotalPages:      int(totalPages) - 1,
		TotalElements:   totalItems,
		SortBy:          []string{"id DESC"},
		CursorStart:     new(string),
		CursorEnd:       new(string),
	}
}

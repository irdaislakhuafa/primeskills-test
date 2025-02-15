package entity

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

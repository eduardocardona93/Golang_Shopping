package dto

// Meta carries pagination metadata returned alongside list endpoints.
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewMeta computes pagination metadata from the requested page/limit and the
// total number of matching records.
func NewMeta(page, limit int, total int64) Meta {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}

	return Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

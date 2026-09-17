package repository

// Pagination carries page-based pagination parameters shared by all list
// endpoints.
type Pagination struct {
	Page  int
	Limit int
}

const (
	defaultPage  = 1
	defaultLimit = 20
	maxLimit     = 100
)

// Normalize applies sane defaults/bounds and returns the SQL OFFSET/LIMIT
// values to use.
func (p Pagination) Normalize() (offset, limit int) {
	page := p.Page
	if page < 1 {
		page = defaultPage
	}

	limit = p.Limit
	if limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	offset = (page - 1) * limit
	return offset, limit
}

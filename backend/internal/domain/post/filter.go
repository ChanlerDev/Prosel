package post

type Pagination struct {
	Page       int
	PerPage    int
	Total      int64
	TotalPages int
}

type PostListFilter struct {
	Page       int
	PerPage    int
	Status     *PostStatus
	Search     string
	CategoryID string
	TagID      string
	Featured   *bool
}

func NormalizePagination(page int, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}

func NewPagination(page int, perPage int, total int64) Pagination {
	page, perPage = NormalizePagination(page, perPage)
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(perPage) - 1) / int64(perPage))
	}
	return Pagination{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages}
}

package models

// PaginationQuery: Standar query param untuk list data
type PaginationQuery struct {
	Page   int    `query:"page" validate:"min=1"`
	Limit  int    `query:"limit" validate:"min=1,max=100"`
	Search string `query:"search"`
	Sort   string `query:"sort"` // e.g. "created_at desc"
}

// PaginationMeta: Metadata untuk response list
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
	TotalData   int64 `json:"total_data"`
	PerPage     int   `json:"per_page"`
}

// IDResponse: Response standar setelah Create (mengembalikan ID baru)
type IDResponse struct {
	ID interface{} `json:"id"` // Bisa uint atau UUID
}
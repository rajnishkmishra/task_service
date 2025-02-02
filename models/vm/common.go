package vm

type IDRequest struct {
	ID uint64 `uri:"id" binding:"required,gte=1"`
}

type PaginationRequest struct {
	PageNumber int64 `form:"p"`
	Limit      int64 `form:"l"`
}

func NewPaginationRequest(p int64, l int64) PaginationRequest {
	return PaginationRequest{
		PageNumber: p,
		Limit:      l,
	}
}

type MetaResponse struct {
	TotalPages  int64 `json:"total_pages"`
	TotalRecord int64 `json:"total_records"`
	PageNumber  int64 `json:"page_number"`
}
